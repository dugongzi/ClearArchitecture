package registry

import (
	"ClearArchitecture/core/env"
	"ClearArchitecture/core/model"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type field struct {
	Type string
	Name string
}

var Commands = []model.Command{
	{
		Name:        "echo",
		Description: "打印问候语",
		Run: func(args []string) error {
			fmt.Println(strings.Join(args, " "))
			return nil
		},
	},

	{
		Name:        "init",
		Description: "初始化清晰架构",
		Run: func(args []string) error {
			dirs := []string{
				"common",
				"core",
				"features",
			}

			targetPath := env.GetRootPath()

			if len(args) > 0 && args[0] != "" {
				targetPath = args[0]
			}

			fmt.Printf("正在初始化清晰架构....\n目标根目录:%s\n", targetPath)

			for _, dir := range dirs {
				fullPath := filepath.Join(targetPath, dir)
				if err := os.MkdirAll(fullPath, 0755); err != nil {
					return err
				}
				fmt.Printf("已创建: %s\n", fullPath)
			}

			return nil
		},
	},

	{
		Name:        "create",
		Description: "创建清晰架构模块",
		Run: func(args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("缺少参数，示例: create -feature user [data|domain|presentation]")
			}

			switch args[0] {
			case "-feature":
				featureName := args[1]

				featureRoot := filepath.Join("features", featureName)
				if err := os.MkdirAll(featureRoot, 0755); err != nil {
					return err
				}

				var paths []string

				if len(args) > 2 && args[2] != "" {
					switch args[2] {
					case "data":
						paths = []string{
							"data/datasources",
							"data/repositories",
							"data/models",
						}
					case "domain":
						paths = []string{
							"domain/models",
							"domain/repositories",
						}
					case "presentation":
						paths = []string{
							"presentation/page",
							"presentation/widgets",
							"presentation/states",
							"presentation/providers",
						}
					default:
						return fmt.Errorf("未知模块类型: %s", args[2])
					}
				} else {
					paths = []string{
						"data/datasources",
						"data/repositories",
						"data/models",
						"domain/models",
						"domain/repositories",
						"presentation/page",
						"presentation/widgets",
						"presentation/states",
						"presentation/providers",
					}
				}

				for _, dir := range paths {
					targetPath := filepath.Join(featureRoot, dir)
					if err := os.MkdirAll(targetPath, 0755); err != nil {
						return err
					}
					fmt.Printf("已创建: %s\n", targetPath)
				}

			case "-model":
				if len(args) < 3 {
					return fmt.Errorf("缺少参数，示例: create -model featureName modelName [-sync]")
				}

				isSync := len(args) > 3 && args[3] == "-sync"
				featureName := args[1]
				modelName := args[2]
				modelClassName := upperFirst(modelName)

				featureRoot := filepath.Join("features", featureName)
				modelPath := filepath.Join(featureRoot, "domain/models", modelName+".dart")
				modelDtoPath := filepath.Join(featureRoot, "data/models", modelName+"_dto.dart")
				if !isSync {
					if err := os.MkdirAll(filepath.Dir(modelPath), 0755); err != nil {
						return err
					}

					modelContent := fmt.Sprintf("import 'package:freezed_annotation/freezed_annotation.dart';\n\npart '%s.freezed.dart';\n\n@freezed\nabstract class %s with _$%s {\n  const %s._();\n\n  const factory %s({\n  }) = _%s;\n}\n",
						modelName,
						modelClassName,
						modelClassName,
						modelClassName,
						modelClassName,
						modelClassName,
					)

					if err := os.WriteFile(modelPath, []byte(modelContent), 0644); err != nil {
						return err
					}

					fmt.Printf("已创建: %s\n", modelPath)
				} else {
					entityContent, err := os.ReadFile(modelPath)
					if err != nil {
						return err
					}

					fields := parseRequiredFields(string(entityContent))
					if len(fields) == 0 {
						return fmt.Errorf("未在实体中找到 required 字段: %s", modelPath)
					}

					if err := os.MkdirAll(filepath.Dir(modelDtoPath), 0755); err != nil {
						return err
					}

					dtoClassName := modelClassName + "Dto"
					packageName := resolveFlutterPackageName(targetPathOrDefault(featureRoot))
					dtoImports := buildDtoImports(packageName, featureName, fields)
					dtoContent := fmt.Sprintf(
						"import 'package:freezed_annotation/freezed_annotation.dart';\nimport 'package:%s/features/%s/domain/models/%s.dart';\n%s\npart '%s_dto.freezed.dart';\npart '%s_dto.g.dart';\n\n@freezed\nabstract class %s with _$%s {\n  const %s._();\n\n  const factory %s({\n%s  }) = _%s;\n\n  factory %s.fromJson(Map<String, dynamic> json) =>\n      _$%sFromJson(json);\n\n  %s toEntity() {\n    return %s(\n%s    );\n  }\n}\n",
						packageName,
						featureName,
						modelName,
						dtoImports,
						modelName,
						modelName,
						dtoClassName,
						dtoClassName,
						dtoClassName,
						dtoClassName,
						buildDtoFields(fields),
						dtoClassName,
						dtoClassName,
						dtoClassName,
						modelClassName,
						modelClassName,
						buildToEntityFields(fields),
					)

					if err := os.WriteFile(modelDtoPath, []byte(dtoContent), 0644); err != nil {
						return err
					}

					fmt.Printf("已创建: %s\n", modelDtoPath)
				}
			default:
				return fmt.Errorf("未知命令类型: %s", args[0])
			}

			return nil
		},
	},
}

func parseRequiredFields(content string) []field {
	lines := strings.Split(content, "\n")
	fields := make([]field, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "required ") {
			continue
		}

		line = strings.TrimSuffix(line, ",")
		parts := strings.Fields(line)
		if len(parts) != 3 {
			continue
		}

		fields = append(fields, field{
			Type: parts[1],
			Name: parts[2],
		})
	}

	return fields
}

func buildDtoFields(fields []field) string {
	var builder strings.Builder

	for _, item := range fields {
		dtoType := toDtoType(item.Type)
		defaultExpr := getDefaultExpr(dtoType)
		builder.WriteString("    @Default(")
		builder.WriteString(defaultExpr)
		builder.WriteString(") ")
		builder.WriteString(dtoType)
		builder.WriteString(" ")
		builder.WriteString(item.Name)
		builder.WriteString(",\n")
	}

	return builder.String()
}

func buildDtoImports(packageName, featureName string, fields []field) string {
	imports := make(map[string]struct{})

	for _, item := range fields {
		if isListType(item.Type) {
			innerType := getListInnerType(item.Type)
			if !isBuiltinType(innerType) {
				imports[fmt.Sprintf("import 'package:%s/features/%s/data/models/%s_dto.dart';", packageName, featureName, toSnakeCase(innerType))] = struct{}{}
			}
			continue
		}

		if !isBuiltinType(item.Type) {
			imports[fmt.Sprintf("import 'package:%s/features/%s/data/models/%s_dto.dart';", packageName, featureName, toSnakeCase(item.Type))] = struct{}{}
		}
	}

	if len(imports) == 0 {
		return ""
	}

	lines := make([]string, 0, len(imports))
	for line := range imports {
		lines = append(lines, line)
	}
	sort.Strings(lines)

	return strings.Join(lines, "\n") + "\n"
}

func buildToEntityFields(fields []field) string {
	var builder strings.Builder

	for _, item := range fields {
		builder.WriteString("      ")
		builder.WriteString(item.Name)
		builder.WriteString(": ")

		if isListType(item.Type) {
			if isBuiltinListType(item.Type) {
				builder.WriteString(item.Name)
			} else {
				builder.WriteString(item.Name)
				builder.WriteString(".map((item) => item.toEntity()).toList()")
			}
		} else if isBuiltinType(item.Type) {
			builder.WriteString(item.Name)
		} else {
			builder.WriteString(item.Name)
			builder.WriteString(".toEntity()")
		}

		builder.WriteString(",\n")
	}

	return builder.String()
}

func toDtoType(fieldType string) string {
	if isListType(fieldType) {
		innerType := getListInnerType(fieldType)
		if isBuiltinType(innerType) {
			return "List<" + innerType + ">"
		}
		return "List<" + innerType + "Dto>"
	}

	if isBuiltinType(fieldType) {
		return fieldType
	}

	return fieldType + "Dto"
}

func getDefaultExpr(dtoType string) string {
	switch dtoType {
	case "String":
		return `""`
	case "int", "double", "num":
		return "0"
	case "bool":
		return "false"
	}

	if isListType(dtoType) {
		return "[]"
	}

	return dtoType + "()"
}

func isBuiltinType(fieldType string) bool {
	switch fieldType {
	case "String", "int", "double", "num", "bool", "dynamic":
		return true
	default:
		return false
	}
}

func isListType(fieldType string) bool {
	return strings.HasPrefix(fieldType, "List<") && strings.HasSuffix(fieldType, ">")
}

func isBuiltinListType(fieldType string) bool {
	if !isListType(fieldType) {
		return false
	}
	return isBuiltinType(getListInnerType(fieldType))
}

func getListInnerType(fieldType string) string {
	return strings.TrimSuffix(strings.TrimPrefix(fieldType, "List<"), ">")
}

func upperFirst(value string) string {
	if value == "" {
		return value
	}
	return strings.ToUpper(value[:1]) + value[1:]
}

func toSnakeCase(value string) string {
	if value == "" {
		return value
	}

	var builder strings.Builder
	for index, char := range value {
		if index > 0 && char >= 'A' && char <= 'Z' {
			builder.WriteByte('_')
		}
		builder.WriteRune(char)
	}

	return strings.ToLower(builder.String())
}

func resolveFlutterPackageName(startPath string) string {
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

func targetPathOrDefault(path string) string {
	if path != "" {
		return path
	}
	return env.GetRootPath()
}
