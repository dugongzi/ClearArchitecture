package childs

import (
	"ClearArchitecture/core/model"
	"fmt"
	"strings"
)

func EchoCommand() model.Command {
	return model.Command{
		Name:        "echo",
		Description: "打印输入内容",
		Usage: []string{
			"echo <text...>",
		},
		Examples: []string{
			"echo hello world",
		},
		Run: func(args []string) error {
			fmt.Println(strings.Join(args, " "))
			return nil
		},
	}
}
