package bobra

import (
	"fmt"
	"os"
)

// 两个自定义异常
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