package bobra

import (
	"errors"
	"fmt"
	"os"
)

var(
	// 当找到 "help" 等命令行参数时抛出
	FoundHelp = errors.New("Found Help")
)
// 当命令没有找到时抛出
type ObjectNotFound struct {
	Type string
	Name string
}

func (e ObjectNotFound) Error() string {
	return fmt.Sprintf("An instance of %s, name '%s' doesn't exist.", e.Type, e.Name)
}

// 打印异常的函数
func LogError(e error) {
	fmt.Fprintln(os.Stderr, e.Error())
}