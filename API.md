

# bobra
`import "github.com/bobbaicloudwithpants/bobra"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Examples](#pkg-examples)

## <a name="pkg-overview">Overview</a>
bobra 是一个模仿了 github.com/spf13/cobra 的包。
bobra 中实现了精简版的 cobra 的功能, 使得命令行程序开发者能够快速的建立耦合度低，高度模块化的命令行程序。




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [func LogError(e error)](#LogError)
* [type Command](#Command)
  * [func (c *Command) AddCommand(cmds ...*Command)](#Command.AddCommand)
  * [func (c *Command) CommandPath() string](#Command.CommandPath)
  * [func (c *Command) Commands() []*Command](#Command.Commands)
  * [func (c *Command) Execute() error](#Command.Execute)
  * [func (c *Command) ExecuteC() (err error)](#Command.ExecuteC)
  * [func (c *Command) Find(args []string) (*Command, []string, error)](#Command.Find)
  * [func (c *Command) Flags() *flag.FlagSet](#Command.Flags)
  * [func (c *Command) GlobalFlags() *flag.FlagSet](#Command.GlobalFlags)
  * [func (c *Command) HasAvailableFlags() bool](#Command.HasAvailableFlags)
  * [func (c *Command) HasAvailableGlobalFlags() bool](#Command.HasAvailableGlobalFlags)
  * [func (c *Command) HasAvailableLocalFlags() bool](#Command.HasAvailableLocalFlags)
  * [func (c *Command) HasAvailableSubCmds() bool](#Command.HasAvailableSubCmds)
  * [func (c *Command) HasParent() bool](#Command.HasParent)
  * [func (c *Command) HasSubCommands() bool](#Command.HasSubCommands)
  * [func (c *Command) IsAvailable() bool](#Command.IsAvailable)
  * [func (c *Command) LocalFlags() *flag.FlagSet](#Command.LocalFlags)
  * [func (c *Command) LongIntroduction() string](#Command.LongIntroduction)
  * [func (c *Command) Name() string](#Command.Name)
  * [func (c *Command) Parent() *Command](#Command.Parent)
  * [func (c *Command) ParseFlags(args []string) error](#Command.ParseFlags)
  * [func (c *Command) Root() *Command](#Command.Root)
  * [func (c *Command) Runnable() bool](#Command.Runnable)
  * [func (c *Command) SetGlobalFlags(flags *flag.FlagSet)](#Command.SetGlobalFlags)
  * [func (c *Command) ShortIntroduction() string](#Command.ShortIntroduction)
  * [func (c *Command) Usage() error](#Command.Usage)
  * [func (c *Command) UsageFunc() (f func(*Command) error)](#Command.UsageFunc)
  * [func (c *Command) UsageTemplate() string](#Command.UsageTemplate)
  * [func (c *Command) UseLine() string](#Command.UseLine)
* [type ObjectNotFound](#ObjectNotFound)
  * [func (e ObjectNotFound) Error() string](#ObjectNotFound.Error)

#### <a name="pkg-examples">Examples</a>
* [Command.AddCommand](#example_Command_AddCommand)
* [Command.CommandPath](#example_Command_CommandPath)
* [Command.Execute](#example_Command_Execute)

#### <a name="pkg-files">Package files</a>
[command.go](/src/github.com/bobbaicloudwithpants/bobra/command.go) [error.go](/src/github.com/bobbaicloudwithpants/bobra/error.go) [utils.go](/src/github.com/bobbaicloudwithpants/bobra/utils.go) 



## <a name="pkg-variables">Variables</a>
``` go
var (
    // 当找到 "help" 等命令行参数时抛出
    FoundHelp = errors.New("Found Help")
)
```


## <a name="LogError">func</a> [LogError](/src/target/error.go?s=388:410#L24)
``` go
func LogError(e error)
```
打印异常的函数




## <a name="Command">type</a> [Command](/src/target/command.go?s=306:1126#L15)
``` go
type Command struct {
    // 命令的使用名称
    Use string
    // 命令的较短介绍
    Short string
    // 命令的较长介绍
    Long string
    // 命令使用介绍
    Example string

    // 运行这个命令执行的函数
    Run func(cmd *Command, args []string)
    // contains filtered or unexported fields
}

```









### <a name="Command.AddCommand">func</a> (\*Command) [AddCommand](/src/target/command.go?s=3726:3772#L168)
``` go
func (c *Command) AddCommand(cmds ...*Command)
```
添加子命令




### <a name="Command.CommandPath">func</a> (\*Command) [CommandPath](/src/target/command.go?s=5982:6020#L253)
``` go
func (c *Command) CommandPath() string
```
返回这条命令从根命令开始向下，直到当前命令c的命令名称组合，用 ' ' 分割




### <a name="Command.Commands">func</a> (\*Command) [Commands](/src/target/command.go?s=5812:5851#L248)
``` go
func (c *Command) Commands() []*Command
```



### <a name="Command.Execute">func</a> (\*Command) [Execute](/src/target/command.go?s=2099:2132#L98)
``` go
func (c *Command) Execute() error
```
执行命令，调用链为：Execute--->ExecuteC--->execute




### <a name="Command.ExecuteC">func</a> (\*Command) [ExecuteC](/src/target/command.go?s=1717:1757#L77)
``` go
func (c *Command) ExecuteC() (err error)
```
找到要执行的命令，或者抛出异常




### <a name="Command.Find">func</a> (\*Command) [Find](/src/target/command.go?s=5008:5073#L208)
``` go
func (c *Command) Find(args []string) (*Command, []string, error)
```
从参数中找到要执行的子命令, 如果没有子命令则返回这个命令本身，如果找不到则返回错误




### <a name="Command.Flags">func</a> (\*Command) [Flags](/src/target/command.go?s=3361:3400#L152)
``` go
func (c *Command) Flags() *flag.FlagSet
```
返回命令的参数列表, 如果 flags 为空则初始化这个flag




### <a name="Command.GlobalFlags">func</a> (\*Command) [GlobalFlags](/src/target/command.go?s=2340:2385#L112)
``` go
func (c *Command) GlobalFlags() *flag.FlagSet
```
获取全局的flags




### <a name="Command.HasAvailableFlags">func</a> (\*Command) [HasAvailableFlags](/src/target/command.go?s=7369:7411#L322)
``` go
func (c *Command) HasAvailableFlags() bool
```
判断命令是否存在有效的flags




### <a name="Command.HasAvailableGlobalFlags">func</a> (\*Command) [HasAvailableGlobalFlags](/src/target/command.go?s=7527:7575#L328)
``` go
func (c *Command) HasAvailableGlobalFlags() bool
```
判断命令是否存在全局有效的flags




### <a name="Command.HasAvailableLocalFlags">func</a> (\*Command) [HasAvailableLocalFlags](/src/target/command.go?s=7697:7744#L334)
``` go
func (c *Command) HasAvailableLocalFlags() bool
```
判断命令是否存在局部有效的flags




### <a name="Command.HasAvailableSubCmds">func</a> (\*Command) [HasAvailableSubCmds](/src/target/command.go?s=6954:6998#L299)
``` go
func (c *Command) HasAvailableSubCmds() bool
```
判断该命令是否有有效的子命令




### <a name="Command.HasParent">func</a> (\*Command) [HasParent](/src/target/command.go?s=7234:7268#L314)
``` go
func (c *Command) HasParent() bool
```
判断 c 是否有父命令




### <a name="Command.HasSubCommands">func</a> (\*Command) [HasSubCommands](/src/target/command.go?s=7130:7169#L309)
``` go
func (c *Command) HasSubCommands() bool
```
判断 c 是否有子命令




### <a name="Command.IsAvailable">func</a> (\*Command) [IsAvailable](/src/target/command.go?s=6789:6825#L291)
``` go
func (c *Command) IsAvailable() bool
```
判断该命令是否有效




### <a name="Command.LocalFlags">func</a> (\*Command) [LocalFlags](/src/target/command.go?s=2984:3028#L137)
``` go
func (c *Command) LocalFlags() *flag.FlagSet
```
返回仅子命令可以使用的局部flags




### <a name="Command.LongIntroduction">func</a> (\*Command) [LongIntroduction](/src/target/command.go?s=5469:5512#L230)
``` go
func (c *Command) LongIntroduction() string
```
返回这条命令的完整介绍，应放在 Usage 的开头




### <a name="Command.Name">func</a> (\*Command) [Name](/src/target/command.go?s=5274:5305#L220)
``` go
func (c *Command) Name() string
```
返回命令的名字




### <a name="Command.Parent">func</a> (\*Command) [Parent](/src/target/command.go?s=1977:2012#L93)
``` go
func (c *Command) Parent() *Command
```
返回当前命令的父命令




### <a name="Command.ParseFlags">func</a> (\*Command) [ParseFlags](/src/target/command.go?s=1165:1214#L49)
``` go
func (c *Command) ParseFlags(args []string) error
```
将args参数转换为flags参数




### <a name="Command.Root">func</a> (\*Command) [Root](/src/target/command.go?s=5714:5747#L240)
``` go
func (c *Command) Root() *Command
```
返回该命令的根命令




### <a name="Command.Runnable">func</a> (\*Command) [Runnable](/src/target/command.go?s=6698:6731#L286)
``` go
func (c *Command) Runnable() bool
```
根据是否存在 Run 函数指针来判断这个命令能否运行




### <a name="Command.SetGlobalFlags">func</a> (\*Command) [SetGlobalFlags](/src/target/command.go?s=2234:2287#L107)
``` go
func (c *Command) SetGlobalFlags(flags *flag.FlagSet)
```
设置全局可用的flags




### <a name="Command.ShortIntroduction">func</a> (\*Command) [ShortIntroduction](/src/target/command.go?s=5617:5661#L235)
``` go
func (c *Command) ShortIntroduction() string
```
返回这条命令的简短介绍，会返回在命令Usage的子命令列表中




### <a name="Command.Usage">func</a> (\*Command) [Usage](/src/target/command.go?s=7824:7855#L339)
``` go
func (c *Command) Usage() error
```
显示命令的使用方法




### <a name="Command.UsageFunc">func</a> (\*Command) [UsageFunc](/src/target/command.go?s=7941:7995#L344)
``` go
func (c *Command) UsageFunc() (f func(*Command) error)
```
返回能够用于输出【使用方法】的函数




### <a name="Command.UsageTemplate">func</a> (\*Command) [UsageTemplate](/src/target/command.go?s=8269:8309#L361)
``` go
func (c *Command) UsageTemplate() string
```



### <a name="Command.UseLine">func</a> (\*Command) [UseLine](/src/target/command.go?s=6160:6194#L261)
``` go
func (c *Command) UseLine() string
```
输出对于这条命令的完整描述




## <a name="ObjectNotFound">type</a> [ObjectNotFound](/src/target/error.go?s=178:234#L14)
``` go
type ObjectNotFound struct {
    Type string
    Name string
}

```
当命令没有找到时抛出










### <a name="ObjectNotFound.Error">func</a> (ObjectNotFound) [Error](/src/target/error.go?s=236:274#L19)
``` go
func (e ObjectNotFound) Error() string
```







- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
