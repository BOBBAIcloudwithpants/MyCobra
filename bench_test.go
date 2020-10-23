package bobra

import (
	"fmt"
	"os"
	"testing"
)

func BenchmarkCommand_Execute(b *testing.B) {
	// 子命令1
	var s1 = &Command{
		Use: "test",
		Short: "test",
		Long: "test",
		Example: "test",
	}

	var s2 = &Command{
		Use: "subtest",
		Short: "subtest",
		Long: "subtest",
		Run: func(cmd *Command, args []string) {
			fmt.Println("this is a test")
		},
	}

	// 根命令
	var r = &Command{
		Use: "root",
		Short: "root",
		Long: "root test",
		Example: "root test",
	}
	r.AddCommand(s1)
	s1.AddCommand(s2)
	os.Args = []string{"root", "test", "subtest"}
	for i := 0; i < b.N; i++ {
		r.Execute()
	}
}
