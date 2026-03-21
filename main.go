package main

import (
	"ClearArchitecture/core/registry/commands"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("请输入命令")
		return
	}

	input := os.Args[1]
	for _, cmd := range commands.Commands {
		if cmd.Name == input {
			if err := cmd.Run(os.Args[2:]); err != nil {
				fmt.Println("执行失败:", err)
			}
			return
		}
	}

	fmt.Println("未知命令:", input)
}
