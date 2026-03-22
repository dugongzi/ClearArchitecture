package childs

import (
	"ClearArchitecture/core/env"
	"ClearArchitecture/core/model"
	"fmt"
	"os"
	"path/filepath"
)

func InitCommand() model.Command {
	return model.Command{
		Name:        "init",
		Description: "初始化清晰架构",
		Usage: []string{
			"init",
			"init <targetPath>",
		},
		Examples: []string{
			"init",
			"init ./demo_project",
		},
		Run: runInitCommand,
	}
}

func runInitCommand(args []string) error {
	targetPath := env.GetRootPath()
	if len(args) > 0 && args[0] != "" {
		targetPath = args[0]
	}

	fmt.Printf("正在初始化清晰架构....\n目标根目录:%s\n", targetPath)

	for _, dir := range initDirectories {
		fullPath := filepath.Join(targetPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
		fmt.Printf("已创建: %s\n", fullPath)
	}

	return nil
}

var initDirectories = []string{
	"common",
	"core",
	"features",
}
