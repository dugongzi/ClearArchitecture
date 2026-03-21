package commands

import (
	"ClearArchitecture/core/model"
	"ClearArchitecture/core/registry/commands/childs"
)

var Commands = []model.Command{
	childs.EchoCommand(),
	childs.InitCommand(),
	childs.CreateCommand(),
}
