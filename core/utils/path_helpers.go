package utils

import (
	"ClearArchitecture/core/env"
	"os"
	"path/filepath"
	"strings"
)

func ResolveFlutterPackageName(startPath string) string {
	searchPath := startPath
	for {
		pubspecPath := filepath.Join(searchPath, "pubspec.yaml")
		content, err := os.ReadFile(pubspecPath)
		if err == nil {
			for _, line := range strings.Split(string(content), "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "name:") {
					name := strings.TrimSpace(strings.TrimPrefix(line, "name:"))
					if name != "" {
						return name
					}
				}
			}
		}

		parentPath := filepath.Dir(searchPath)
		if parentPath == searchPath {
			break
		}
		searchPath = parentPath
	}

	return "your_project_name"
}

func TargetPathOrDefault(path string) string {
	if path != "" {
		return path
	}
	return env.GetRootPath()
}
