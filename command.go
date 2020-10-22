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

	// 存放FlagSet错误输出的缓冲区
	flagErrorBuf *bytes.Buffer
	// 命令的介绍模版
	usageTemplate string
	// 子命令的列表
	commands []*Command

	// 父命令的指针
	parent *Command

	// Max lengths of commands' string lengths for use in padding.
	commandsMaxUseLen         int
	commandsMaxCommandPathLen int
	commandsMaxNameLen        int


	// 运行这个命令执行的函数
	Run func(cmd *Command, args []string)

	usageFunc func(*Command) error

}

//// 为指令设置helpCommand
//func (c *Command) SetHelpCommand(cmd *Command) {
//	c.helpCommand = cmd
//}

//// 初始化默认的 helpCommand, 如果 c 没有任何子命令或者 c 已经有设置好的helpCommand, 则不设置
//func (c *Command) SetDefaultHelpCmd() {
//	if !c.HasSubCommands() || c.helpCommand != nil{
//		return
//	}
//
//	usePath := c.UsePath()
//	c.helpCommand = &Command{
//		Use: usePath + " [command]",
//		Short: "Help about the usage of several subcommands",
//		Long: "Help you to have ideas of how to use subcommands to satisfy your demands"
//
//	}
//}


// 将args参数转换为flags参数
func (c *Command) ParseFlags(args []string) error{

	if c.flagErrorBuf == nil {
		c.flagErrorBuf = new(bytes.Buffer)
	}

	beforeBufferLen := c.flagErrorBuf.Len()

	c.inheritGlobalFlags()
	err := c.Flags().Parse(args)
	if c.flagErrorBuf.Len() - beforeBufferLen > 0 && err == nil {
		fmt.Println(c.flagErrorBuf.String())
	}
	return err
}



// 根据flag参数执行该命令
func (c *Command) execute(a []string) error {

	err := c.ParseFlags(a)
	if err != nil {
		return err
	}
	c.Usage()
	c.Run(c, a)
	return nil
}

// 找到要执行的命令，或者抛出异常
func (c *Command) ExecuteC() (err error) {
	args := os.Args
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

func (c *Command) SetGlobalFlags(flags *flag.FlagSet) {
	c.globalflags = flags
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
		return
	}

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

		nameLen := len(x.Name())
		if nameLen > c.commandsMaxNameLen {
			c.commandsMaxNameLen = nameLen
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

// 返回命令的名字
func (c* Command) Name() string {
	name := c.Use
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// 返回这条命令的完整介绍，应放在 Usage 的开头
func (c *Command) LongIntroduction() string {
	return c.Long
}


func (c *Command) ShortIntroduction() string {
	return c.Short
}

// 返回该命令的根命令
func (c *Command) Root() *Command {
	p := c
	for p.parent != nil {
		p = c.parent
	}
	return p
}

func (c *Command) Commands() []*Command {
	return c.commands
}

// 返回这条命令从根命令开始向下，直到当前命令c的命令名称组合，用 ' ' 分割
func (c *Command) CommandPath() string {
	if c.HasParent() {
		return c.Parent().CommandPath() + " " + c.Name()
	}
	return c.Name()
}

// 输出对于这条命令的完整描述
func (c *Command) UseLine() string {
	var useline string
	if c.HasParent() {
		useline = c.parent.CommandPath() + " " + c.Use
	} else {
		useline = c.Use
	}

	if c.HasAvailableFlags() && !strings.Contains(useline, "[flags]") {
		useline += " [flags]"
	}
	return useline
}

// 根据命令的名称寻找子命令
func (c *Command) findSubCmd(cmdUse string) *Command {
	for _, cmd := range c.commands {
		if cmd.Name() == cmdUse {
			return cmd
		}
	}
	return nil
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}



func (c *Command) IsAvailable() bool {
	if c.Runnable() || c.HasAvailableSubCmds() {
		return true
	}
	return false
}

func (c *Command) HasAvailableSubCmds() bool {
	for _, sub := range c.commands {
		if sub.IsAvailable() {
			return true
		}
	}
	return false
}

// 判断 c 是否有子命令
func (c *Command) HasSubCommands() bool {
	return len(c.commands) > 0
}

// 判断 c 是否有父命令
func (c *Command) HasParent() bool{
	if c.parent != nil {
		return true
	}
	return false
}

func (c *Command) HasAvailableFlags() bool {
	return c.Flags().HasAvailableFlags()
}

func (c *Command) HasAvailableGlobalFlags() bool {
	return c.GlobalFlags().HasAvailableFlags()
}

func (c *Command) HasAvailableLocalFlags() bool {
	return c.LocalFlags().HasAvailableFlags()
}


//
//var minNamePadding = 11
//
//// NamePadding returns padding for the name.
//func (c *Command) NamePadding() int {
//	if c.parent == nil || minNamePadding > c.parent.commandsMaxNameLen {
//		return minNamePadding
//	}
//	return c.parent.commandsMaxNameLen
//}

func (c *Command) Usage() error {
	return c.UsageFunc()(c)
}

func (c *Command) UsageFunc() (f func(*Command) error) {
	if c.usageFunc != nil {
		return c.usageFunc
	}
	if c.HasParent() {
		return c.Parent().UsageFunc()
	}
	return func(c *Command) error {
		c.inheritGlobalFlags()
		err := tmpl(os.Stdout, c.UsageTemplate(), c)
		if err != nil {
			LogError(err)
		}
		return err
	}
}

func (c *Command) UsageTemplate() string {
	if c.usageTemplate != "" {
		return c.usageTemplate
	}

	if c.HasParent() {
		return c.parent.UsageTemplate()
	}
	return `
{{.LongIntroduction}}

Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCmds}}
  {{.CommandPath}} [command]

Available Commands:{{range .Commands}}{{if .IsAvailable}}
  {{.Name}}: {{.ShortIntroduction}}{{end}}
  {{end}}{{end}}{{if .HasAvailableLocalFlags}}
LocalFlags:
  {{.LocalFlags.FlagUsages}}
{{end}}{{if .HasAvailableGlobalFlags}}
GlobalFlags:
  {{.GlobalFlags.FlagUsages}}
{{end}} {{if .HasAvailableSubCmds}}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}

`
}


//`Usage:{{if .Runnable}}
//  {{.UseLine}}{{end}}{{if .HasAvailableSubCmds}}
//  {{.CommandPath}} [command]{{end}}
//
//Examples:
//{{.Example}}{{end}}{{if .HasAvailableSubCmds}}
//
//Available Commands:{{range .Commands}}{{if .IsAvailable}}
//  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}
//
//Local Flags:
//{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableGlobalFlags}}
//
//Global Flags:
//{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableSubCmds}}
//
//Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
//`
//
