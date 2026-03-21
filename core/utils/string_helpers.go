package utils

import "strings"

func UpperFirst(value string) string {
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
