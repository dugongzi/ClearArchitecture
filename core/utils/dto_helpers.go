package utils

import (
	"fmt"
	"sort"
	"strings"
)

type field struct {
	Type string
	Name string
}

func ParseRequiredFields(content string) []field {
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

func BuildDtoFields(fields []field) string {
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

func BuildDtoImports(packageName, featureName string, fields []field) string {
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

func BuildToEntityFields(fields []field) string {
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
