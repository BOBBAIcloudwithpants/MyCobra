# Bobra 设计报告
## 个人信息
|      |          |
| ---- | -------- |
| 姓名 | 白家栋   |
| 学号 | 18342001 |
| 专业 | 软件工程 |

## 结构设计

### 1. Command
本次 Bobra 的设计我大量的参考了 Cobra，也就是我们要模仿的本尊，整个包的结构都是通过 `Command` 这个结构体支撑起来的，下面我们首先看一下 Command 的定义:
```go
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
	localflags *flag.FlagSet

	// 存放FlagSet错误输出的缓冲区
	flagErrorBuf *bytes.Buffer
	// 命令的介绍模版
	usageTemplate string
	// 子命令的列表
	commands []*Command

	// 父命令的指针
	parent *Command

	// 运行这个命令执行的函数
	Run func(cmd *Command, args []string)

	// 该 Command 的使用方法介绍
	usageFunc func(*Command) error
}
```
可以看到，Command 定义中包含了 `Run` 函数成员，这是执行该 command 时运行的主体函数，这使得用户在定义命令行程序的时候能够进行基本的功能逻辑设置；同时，为了满足命令行程序能够定义**命令行参数**的需求，我引入了第三方包 [pflag](https://github.com/spf13/pflag). 每个 Command 都有 `flag.Flagset` 类型的指针成员： `flags`, `globalflags`, `localflags`。为什么要设计这么多种 flagset 成员，设置flag的时候应当如何使用？这就涉及到命令之间的组合问题。因为我们要满足命令行程序中的子命令，这就需要每个 Command 都拥有自己的**子Command**, 同时也要有自己的**父Command**, 这样的结构通过 **commands**, `[]*Command` 类型的成员变量，和 **parent**, `*Command` 类型的成员变量实现，最后的效果类似于数据结构中的**多叉树**，如下图所示:    
![](https://tva1.sinaimg.cn/large/0081Kckwgy1gjz9zpj3yzj30dd08xdg7.jpg)


### 2. flags
在设计的过程中，我定义了这样3中 **flags**:    

#### (1). GlobalFlags
- 含义: 全局命令行参数，在每个 Command 节点都有效。**在一棵命令树中仅能全局维护同一个 globalflag 指针**
- 用法: 对于一些比较通用的参数，比如 `-v, version`; `-h, help`，可以通过 GlobalFlags 的语义来设置。

#### (2). LocalFlags
- 含义: 局部命令行参数，仅在某个 Command 节点有效。
- 用法: 对于一些功能比较专一的参数可以通过 GlobalFlags 的语义来设置。


#### (3). Flags
- 含义: 某个节点对应的全部Flags, 包含了这个节点的 **globalflags** 和 **localflags**。
- 用法: 在从flags中得到用户的命令行参数时，访问这个flags指针来获取参数的值


对于这三种 flags，我分别定义了它们的 Get, Set 方法；当用户需要定义命令行参数的时候，只能够调用 `GlobalFlags()` 以及 `LocalFlags()` 来进行定义。    
举个例子，比如此时要定义一个bool类型的命令行参数，使得这个参数仅在 `sub1` 有效，在其他的所有命令都无效, 那么在命令行程序中，需要这样定义:
```go
sub1.LocalFlags().BoolVarP()
```

那如果要定义一个bool类型的命令行参数，使得这个参数在所有的命令下都有效，那么在命令行程序中，需要这样定义:
```go
sub1.GlobalFlags().BoolVarP()
```

## 功能实现
要执行一个指令对应的功能，首先要通过命令行的输入找到真正要执行的命令. 比如, 命令行输入为: `root sub1`，那么 bobra 包就会根据解析到的两个命令行输入：root 和 sub1 来按照顺序判断这个命令是否存在，如果不存在则抛出异常，如果存在则继续向下寻找，直到到达 `sub1`, sub1 就是本次输入真正要执行的指令。此时，bobra会将 `「非命令」` 的命令行参数解析出来，并调用本次要执行的指令的 `Parse` 函数进行解析。如果解析发生异常就返回，如果没有发生异常就进入到本次指令的 `Run` 函数，也就是指令对应的真正逻辑。总结一下上述的过程，如下图所示: 
![](https://tva1.sinaimg.cn/large/0081Kckwgy1gjzeyix8wvj30ct0g5ab0.jpg)

## 指令功能显示

特别要提到的是，在阅读 **cobra** 的文档的时候，它在输出指令的用法的时候使用了 golang 的 `text/template` 官方库，这个库能够让用户自定义输出的模版，然后根据传入的参数动态替换模版的内容，从而很方便的动态输出具有特定格式的标准输出，我的模版定义如下:

```go
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
  {{.Name}}: {{.ShortIntroduction}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}
LocalFlags:
  {{.LocalFlags.FlagUsages}}
{{end}}{{if .HasAvailableGlobalFlags}}
GlobalFlags:
  {{.GlobalFlags.FlagUsages}}
{{end}} {{if .HasAvailableSubCmds}}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}
```
在要输出的时候动态的将 command 对象代入模版即可输出指令的简介，例如:
![](https://tva1.sinaimg.cn/large/0081Kckwgy1gjzf7lo701j31bi0dyq4n.jpg)