package childs

import (
	"ClearArchitecture/core/model"
	"fmt"
	"strings"
)

func EchoCommand() model.Command {
	return model.Command{
		Name:        "echo",
		Description: "打印问候语",
		Run: func(args []string) error {
			fmt.Println(strings.Join(args, " "))
			return nil
		},
	}
}
