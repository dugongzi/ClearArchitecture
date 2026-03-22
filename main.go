package main

import (
	"ClearArchitecture/core/registry/commands"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Print(commands.HelpText())
		return
	}

	input := os.Args[1]
	if commands.IsHelpArg(input) {
		fmt.Print(commands.HelpText())
		return
	}

	cmd, ok := commands.Find(input)
	if !ok {
		fmt.Println("未知命令:", input)
		fmt.Print(commands.HelpText())
		return
	}

	if len(os.Args) >= 3 && commands.IsHelpArg(os.Args[2]) {
		fmt.Print(commands.CommandHelpText(cmd))
		return
	}

	if err := cmd.Run(os.Args[2:]); err != nil {
		fmt.Println("执行失败:", err)
		fmt.Print(commands.CommandHelpText(cmd))
	}
}
