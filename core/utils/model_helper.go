package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateDomainModel(modelPath, modelName, modelClassName string) error {
	if err := os.MkdirAll(filepath.Dir(modelPath), 0755); err != nil {
		return err
	}

	modelContent := fmt.Sprintf(
		"import 'package:freezed_annotation/freezed_annotation.dart';\n\npart '%s.freezed.dart';\n\n@freezed\nabstract class %s with _$%s {\n  const %s._();\n\n  const factory %s({\n  }) = _%s;\n}\n",
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
	return nil
}

func SyncModelToDto(featureName, featureRoot, modelName, modelClassName, modelPath, modelDtoPath string) error {
	entityContent, err := os.ReadFile(modelPath)
	if err != nil {
		return err
	}

	fields := ParseRequiredFields(string(entityContent))
	if len(fields) == 0 {
		return fmt.Errorf("未在实体中找到 required 字段: %s", modelPath)
	}

	if err := os.MkdirAll(filepath.Dir(modelDtoPath), 0755); err != nil {
		return err
	}

	dtoClassName := modelClassName + "Dto"
	packageName := ResolveFlutterPackageName(TargetPathOrDefault(featureRoot))
	dtoImports := BuildDtoImports(packageName, featureName, fields)
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
		BuildDtoFields(fields),
		dtoClassName,
		dtoClassName,
		dtoClassName,
		modelClassName,
		modelClassName,
		BuildToEntityFields(fields),
	)

	if err := os.WriteFile(modelDtoPath, []byte(dtoContent), 0644); err != nil {
		return err
	}

	fmt.Printf("已创建: %s\n", modelDtoPath)
	return nil
}
