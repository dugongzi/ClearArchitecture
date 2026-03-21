package childs

import (
	"ClearArchitecture/core/model"
	"ClearArchitecture/core/utils"
	"fmt"
	"os"
	"path/filepath"
)

func CreateCommand() model.Command {
	return model.Command{
		Name:        "create",
		Description: "创建清晰架构模块",
		Run:         runCreateCommand,
	}
}

func runCreateCommand(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("缺少参数，示例: create -feature user [data|domain|presentation]")
	}

	switch args[0] {
	case "-feature":
		return createFeature(args)
	case "-model":
		return createModel(args)
	case "-flow":
		return createFlow(args)
	default:
		return fmt.Errorf("未知命令类型: %s", args[0])
	}
}

func createFeature(args []string) error {
	featureName := args[1]
	featureRoot := filepath.Join("features", featureName)

	if err := os.MkdirAll(featureRoot, 0755); err != nil {
		return err
	}

	paths, err := resolveFeaturePaths(args)
	if err != nil {
		return err
	}

	for _, dir := range paths {
		targetPath := filepath.Join(featureRoot, dir)
		if err := os.MkdirAll(targetPath, 0755); err != nil {
			return err
		}
		fmt.Printf("已创建: %s\n", targetPath)
	}

	return nil
}

func createModel(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("缺少参数，示例: create -model featureName modelName [-sync]")
	}

	featureName := args[1]
	modelName := args[2]
	modelClassName := utils.UpperFirst(modelName)
	featureRoot := filepath.Join("features", featureName)
	modelPath := filepath.Join(featureRoot, "domain/models", modelName+".dart")
	modelDtoPath := filepath.Join(featureRoot, "data/models", modelName+"_dto.dart")

	if len(args) > 3 && args[3] == "-sync" {
		return utils.SyncModelToDto(featureName, featureRoot, modelName, modelClassName, modelPath, modelDtoPath)
	}

	return utils.CreateDomainModel(modelPath, modelName, modelClassName)
}

func createFlow(args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("缺少参数，示例: create -flow query featureName prefix [submodule]")
	}

	flowType := args[1]
	featureName := args[2]
	prefix := args[3]
	submodule := ""
	if len(args) > 4 {
		submodule = args[4]
	}

	return utils.CreateDataFlow(utils.DataFlowOptions{
		FeatureName: featureName,
		Prefix:      prefix,
		Submodule:   submodule,
		FlowType:    flowType,
	})
}

func resolveFeaturePaths(args []string) ([]string, error) {
	if len(args) <= 2 || args[2] == "" {
		return allFeaturePaths, nil
	}

	switch args[2] {
	case "data":
		return dataFeaturePaths, nil
	case "domain":
		return domainFeaturePaths, nil
	case "presentation":
		return presentationFeaturePaths, nil
	default:
		return nil, fmt.Errorf("未知模块类型: %s", args[2])
	}
}

var dataFeaturePaths = []string{
	"data/datasources",
	"data/repositories",
	"data/models",
}

var domainFeaturePaths = []string{
	"domain/models",
	"domain/repositories",
}

var presentationFeaturePaths = []string{
	"presentation/page",
	"presentation/widgets",
	"presentation/states",
	"presentation/providers",
}

var allFeaturePaths = []string{
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
