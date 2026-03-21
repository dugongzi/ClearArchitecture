package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DataFlowOptions struct {
	FeatureName string
	Prefix      string
	Submodule   string
	FlowType    string
}

func CreateDataFlow(options DataFlowOptions) error {
	flowSuffix, err := resolveFlowSuffix(options.FlowType)
	if err != nil {
		return err
	}

	moduleName := SnakeToPascal(options.Prefix)
	variableName := LowerCamel(moduleName)

	files := []generatedFile{
		{
			path: buildFlowFilePath(options.FeatureName, "data/datasources", options.Submodule, options.Prefix+"_"+flowSuffix+"_datasource.dart"),
			content: fmt.Sprintf(
				"class %s%sDatasource {\n  final HttpService _httpService;\n\n  %s%sDatasource({required HttpService httpService})\n      : _httpService = httpService;\n}\n",
				moduleName,
				flowSuffix,
				moduleName,
				flowSuffix,
			),
		},
		{
			path: buildFlowFilePath(options.FeatureName, "domain/repositories", options.Submodule, options.Prefix+"_"+flowSuffix+"_repository.dart"),
			content: fmt.Sprintf(
				"abstract class %s%sRepository {\n}\n",
				moduleName,
				flowSuffix,
			),
		},
		{
			path: buildFlowFilePath(options.FeatureName, "data/repositories", options.Submodule, options.Prefix+"_"+flowSuffix+"_repository_impl.dart"),
			content: fmt.Sprintf(
				"class %s%sRepositoryImpl implements %s%sRepository {\n  final %s%sDatasource dataSource;\n\n  %s%sRepositoryImpl({required this.dataSource});\n}\n",
				moduleName,
				flowSuffix,
				moduleName,
				flowSuffix,
				moduleName,
				flowSuffix,
				moduleName,
				flowSuffix,
			),
		},
		{
			path: buildFlowFilePath(options.FeatureName, "presentation/providers", options.Submodule, options.Prefix+"_"+flowSuffix+"_provider.dart"),
			content: fmt.Sprintf(
				"@riverpod\n%s%sRepository %s%sRepository(Ref ref) {\n  final httpService = ref.watch(httpServiceProvider);\n  final dataSource = %s%sDatasource(httpService: httpService);\n  return %s%sRepositoryImpl(dataSource: dataSource);\n}\n",
				moduleName,
				flowSuffix,
				variableName,
				flowSuffix,
				moduleName,
				flowSuffix,
				moduleName,
				flowSuffix,
			),
		},
	}

	for _, file := range files {
		if err := os.MkdirAll(filepath.Dir(file.path), 0755); err != nil {
			return err
		}

		if err := os.WriteFile(file.path, []byte(file.content), 0644); err != nil {
			return err
		}

		fmt.Printf("已创建: %s\n", file.path)
	}

	return nil
}

type generatedFile struct {
	path    string
	content string
}

func buildFlowFilePath(featureName, baseDir, submodule, fileName string) string {
	parts := []string{"features", featureName, baseDir}
	if strings.TrimSpace(submodule) != "" {
		parts = append(parts, submodule)
	}
	parts = append(parts, fileName)
	return filepath.Join(parts...)
}

func resolveFlowSuffix(flowType string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(flowType)) {
	case "query":
		return "Query", nil
	case "action":
		return "Action", nil
	default:
		return "", fmt.Errorf("未知数据流类型: %s，支持 query 或 action", flowType)
	}
}
