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

func SnakeToPascal(value string) string {
	if value == "" {
		return value
	}

	parts := strings.Split(value, "_")
	var builder strings.Builder
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		builder.WriteString(UpperFirst(strings.ToLower(part)))
	}

	return builder.String()
}

func LowerCamel(value string) string {
	if value == "" {
		return value
	}

	return strings.ToLower(value[:1]) + value[1:]
}
