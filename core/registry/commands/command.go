package commands

import (
	"ClearArchitecture/core/model"
	"ClearArchitecture/core/registry/commands/childs"
	"fmt"
	"strings"
)

var Commands = []model.Command{
	childs.EchoCommand(),
	childs.InitCommand(),
	childs.CreateCommand(),
}

func IsHelpArg(input string) bool {
	return input == "-help" || input == "--help" || input == "help"
}

func Find(name string) (model.Command, bool) {
	for _, cmd := range Commands {
		if cmd.Name == name {
			return cmd, true
		}
	}

	return model.Command{}, false
}

func HelpText() string {
	var builder strings.Builder
	builder.WriteString("可用命令:\n\n")

	for _, cmd := range Commands {
		builder.WriteString(CommandHelpText(cmd))
		builder.WriteString("\n")
	}

	builder.WriteString("帮助命令:\n")
	builder.WriteString("  -help | --help | help\n")
	builder.WriteString("  <command> -help\n")

	return builder.String()
}

func CommandHelpText(cmd model.Command) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s\n", cmd.Name))

	if cmd.Description != "" {
		builder.WriteString(fmt.Sprintf("  说明: %s\n", cmd.Description))
	}

	if len(cmd.Usage) > 0 {
		builder.WriteString("  用法:\n")
		for _, usage := range cmd.Usage {
			builder.WriteString(fmt.Sprintf("    %s\n", usage))
		}
	}

	if len(cmd.Examples) > 0 {
		builder.WriteString("  示例:\n")
		for _, example := range cmd.Examples {
			builder.WriteString(fmt.Sprintf("    %s\n", example))
		}
	}

	return builder.String()
}
