package bobra

import (
	"bytes"
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
	"strings"
)


type Command struct {
	// 命令的使用名称
	Use string
	// 命令的较短介绍
	Short string
	// 命令的较长介绍
	Long string
	// 命令使用介绍
	Example string
	// 这个命令对应的全部flags,为 globalflags + localflags
	flags *flag.FlagSet
	// 这个命令集合对应的全部全局可用的flag
	globalflags *flag.FlagSet
	// 这个命令集合对应的局部可用的flag，即仅当前命令可以使用的flag
	localflags  *flag.FlagSet

	// 这个命令调用时用到的参数
	args []string

	// 存放FlagSet错误输出的缓冲区
	flagErrorBuf *bytes.Buffer

	commands []*Command
	parent *Command
	// The *Run functions are executed in the following order:
	//   * PersistentPreRun()
	//   * PreRun()
	//   * Run()
	//   * PostRun()
	//   * PersistentPostRun()
	// All functions get the same args, the arguments after the command name.
	//
	// PersistentPreRun: children of this command will inherit and execute.

	Run func(cmd *Command, args []string)
	// RunE: Run but returns an error.
}





// 将args参数转换为flags参数
func (c *Command) ParseFlags(args []string) error{

	if c.flagErrorBuf == nil {
		c.flagErrorBuf = new(bytes.Buffer)
	}

	beforeBufferLen := c.flagErrorBuf.Len()
	fmt.Println("args ready to be parse")
	fmt.Println(args)

	c.inheritGlobalFlags()
	err := c.Flags().Parse(args)
	if c.flagErrorBuf.Len() - beforeBufferLen > 0 && err == nil {
		fmt.Println(c.flagErrorBuf.String())
	}
	fmt.Println(c.GlobalFlags().GetString("author"))
	return err
}

// 根据flag参数执行该命令
func (c *Command) execute(a []string) error {
	fmt.Println("args: ")
	fmt.Println(a)
	err := c.ParseFlags(a)
	if err != nil {
		return err
	}
	c.Run(c, a)
	return nil
}

// 找到要执行的命令，或者抛出异常
func (c *Command) ExecuteC() (err error) {
	args := os.Args
	fmt.Println(args)
	cmd, flags, err := c.Find(args)

	if err != nil {
		return err
	}
	return cmd.execute(flags)
}

// 返回当前命令的父命令
func (c *Command) Parent() *Command {
	return c.parent
}

// 执行命令，调用链为：Execute--->ExecuteC--->execute
func (c *Command) Execute()  {
	err := c.ExecuteC()
	if err != nil {
		panic(err)
	}
}

// 获取全局的flags
func (c *Command) GlobalFlags() *flag.FlagSet {
	c.inheritGlobalFlags()
	if c.globalflags == nil {
		c.globalflags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		if c.flagErrorBuf == nil {
			c.flagErrorBuf = new(bytes.Buffer)
		}
		c.globalflags.SetOutput(c.flagErrorBuf)
	}
	return c.globalflags
}


func (c *Command) inheritGlobalFlags() {
	// 如果为根命令，终止
	if c.Parent() == nil {
		fmt.Println("is parent")
		return
	}
	fmt.Println("is not parent")
	// 否则继承父亲的globalflags, 一个指令集下应当维护一个全局唯一的globalflags指针
	c.globalflags = c.Parent().GlobalFlags()
}

// 返回仅子命令可以使用的局部flags
func (c *Command) LocalFlags() *flag.FlagSet {
	c.inheritGlobalFlags()

	if c.localflags == nil {
		c.localflags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		if c.flagErrorBuf == nil {
			c.flagErrorBuf = new(bytes.Buffer)
		}
		c.localflags.SetOutput(c.flagErrorBuf)
	}

	//addToLocal := func(f *flag.Flag) {
	//	if c.localflags.Lookup(f.Name) == nil && c.globalflags.Lookup(f.Name) == nil {
	//		c.localflags.AddFlag(f)
	//	}
	//}
	//c.Flags().VisitAll(addToLocal)
	//c.GlobalFlags().VisitAll(addToLocal)
	//c.Flags().AddFlagSet(c.localflags)
	return c.localflags
}


// 返回命令的参数列表, 如果 flags 为空则初始化这个flag
func (c *Command) Flags() *flag.FlagSet {
	c.inheritGlobalFlags()
	if c.flags == nil {
		c.flags = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		if c.flagErrorBuf == nil {
			c.flagErrorBuf = new(bytes.Buffer)
		}
		c.flags.SetOutput(c.flagErrorBuf)
	}
	c.flags.AddFlagSet(c.localflags)
	c.flags.AddFlagSet(c.globalflags)

	return c.flags
}


// 添加子命令
func (c *Command) AddCommand(cmds ...*Command) {
	for i, x := range cmds {
		if cmds[i] == c {
			panic("Command can't be a child of itself")
		}
		cmds[i].parent = c
		c.commands = append(c.commands, x)
	}
}

// 递归寻找下一个要执行的子命令，如果找不到则抛出异常
func innerFind(cmd *Command, innerArgs []string)(*Command, []string, error) {


	// 参数列表中的第一个一定与cmd的 Name 相同
	if innerArgs[0] != cmd.Name() {
		return cmd, nil, ObjectNotFound{Type: "Command", Name: innerArgs[0]}
	}

	innerArgsWithoutFlags := stripFlags(innerArgs[1:], cmd)
	fmt.Println("after strip")
	fmt.Println(innerArgsWithoutFlags)
	fmt.Println(innerArgs)
	// 如果此时已经没有向下的子命令了
	if len(innerArgsWithoutFlags) == 0 {
		return cmd, innerArgs[1:], nil
	}
	// 否则此时已经有一个子命令了
	sub := innerArgsWithoutFlags[0]

	subCmd := cmd.findSubCmd(sub)
	if subCmd == nil {
		return cmd, nil, ObjectNotFound{Type: "Command", Name: sub}
	}
	return innerFind(subCmd, innerArgs[1:])
}

// 从参数中找到要执行的子命令, 如果没有子命令则返回这个命令本身，如果找不到则返回错误
func (c *Command) Find(args []string) (*Command,[]string , error) {
	cmd , flags, err := innerFind(c, args)
	if err != nil {
		return cmd, flags, err
	}
	return cmd, flags, nil
}

func (c* Command) Name() string {
	name := c.Use
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) findSubCmd(cmdUse string) *Command {
	for _, cmd := range c.commands {
		if cmd.Name() == cmdUse {
			return cmd
		}
	}
	return nil
}


