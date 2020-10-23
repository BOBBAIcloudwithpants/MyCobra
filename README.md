# Bobra

## 简介
bobra 是一个模仿了 [github.com/spf13/cobra]("https://github.com/spf13/cobra.git") 的包，bobra 中实现了精简版的 cobra 的功能, 使得命令行程序开发者能够快速的建立耦合度低，高度模块化的命令行程序。

## 功能
- 快速定义命令行程序
- 支持定义多个子命令
- 默认生成命令的介绍，也支持命令介绍的自定义

## 文档链接
- API文档: https://github.com/BOBBAIcloudwithpants/bobra/wiki/Bobra-API-Document
- 设计文档: https://github.com/BOBBAIcloudwithpants/bobra/blob/main/specification.md

## 环境以及获取方式
- golang 版本: golang 1.14及以上
- 操作系统: mac/linux
- 获取方式:
```
go get github.com/bobbaicloudwithpants/bobra
```

## 使用
该使用示例的代码全部都可以在 [resume 仓库]("https://github.com/BOBBAIcloudwithpants/resume.git") 中获取。    
假设您此时要开发一款命令行的 app，通过命令行来按照需求输出您的学历以及个人信息，即:
```
resume          // 显示resume的用法
resume name     // 显示姓名
resume edu -p   // 显示小学
resume edu -m   // 显示中学
resume edu -c   // 显示大学
```
可以看到，这个需求中有2个子命令：`name`, `edu`; 其中 `edu` 下还有3个参数 `-pmc`。    
使用`bobra`, 您的项目目录结构大致如下:
```go
.
├── cmd                     // 存放每条指令的目录
│   ├── edu.go
│   ├── name.go
│   └── resume.go
└── main.go                 // 命令行程序入口
```
`main.go` 中只需要直接调用cmd中的 `Execute` 即可，即:
```go
package main

import (
	"github.com/bobbaicloudwithpants/resume/cmd"
)
func main(){
	cmd.Execute()
}
```
Execute 实际上是对于根命令的 `Run` 的一层封装。该项目中，我们的根命令是 `resume` , resume 的定义在 cmd/resume.go 中，如下: 
```go
package cmd

import (
	cobra "github.com/bobbaicloudwithpants/bobra"
)

var resume = &cobra.Command{
	Use: "resume",			// Use 指定了这个命令的名字
	Short: "resume is a simple self-introduction cli program",	// Short 是对于该命令的简短介绍
	Long: "resume makes you organize your personal resume properly, and display in a user-friendly and cleary way.",	// Long 是命令的比较完整的介绍
	Run: func(c *cobra.Command, args []string) {
		c.Usage()
	},
}

func Execute() {
	resume.Execute()
}
```

子命令的定义分别位于 `edu.go` 和 `name.go` 中，子命令的定义中，除了要规定命令的基本属性之外，还要在 `func init()` 中将该子命令加入到根命令下，比如以 name.go 为例:
```go
package cmd

import (
	"fmt"
	cobra "github.com/bobbaicloudwithpants/bobra"
)
var name = &cobra.Command{
	Use : "name",
	Short: "name screens user's name",
	Long: "name displays the user's name put into the memory in advance.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bjd")
	},
}

func init() {       // 这里调用根命令的 AddCommand 方法，便将子命令加入其中
	resume.AddCommand(name)
}
```

我们的需求中还设计一些命令行参数的添加，所以在需要添加参数的指令中，我们可以很方便的在 init 函数中添加想要的参数，以 `edu.go` 为例:   
```go
package cmd

import (
	"fmt"
	cobra "github.com/bobbaicloudwithpants/bobra"
)

var edu = &cobra.Command{
	Use: "edu",
	Short: "edu stores and displays you educational background",
	Long: "edu reads the data put into the program in advance, and dynamically chooses which item to show based on the given parameters.",
	Run: func(cmd *cobra.Command, args []string) {
		if ok, _ := cmd.Flags().GetBool("college");ok {
			fmt.Println("SYSU")
		} else if ok, _ := cmd.Flags().GetBool("middle");ok{
			fmt.Println("THSchool")
		} else if ok, _ := cmd.Flags().GetBool("primary");ok{
			fmt.Println("TS2FX")
		}
	},
}

func init(){
  // 可以看到，我们按照需求为 edu 添加了三个命令行参数、
  // bobra 在参数上引用了 github.com/spf13/pflags 库，这里仅展示部分用法，更多用法您可以在官方文档上查阅
	edu.Flags().BoolP("college", "c", false, "whether show college")
	edu.Flags().BoolP("middle", "m", false, "whether show middle school")
	edu.Flags().BoolP("primary", "p", false, "whether show primary")
	resume.AddCommand(edu)
}
```

项目完成以后，执行 `go install`, 并在终端执行操作，得到的结果如下:     
1. resume
![](https://tva1.sinaimg.cn/large/0081Kckwgy1gjz14jqmhlj31ak0dydhj.jpg)

2. resume name
![](https://tva1.sinaimg.cn/large/0081Kckwgy1gjz15pxjtaj30t002ct8u.jpg)

3. resume edu -c
![](https://tva1.sinaimg.cn/large/0081Kckwgy1gjz16phf6jj30su0220sv.jpg)

