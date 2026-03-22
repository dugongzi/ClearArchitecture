package model

type Command struct {
	Name        string
	Description string
	Usage       []string
	Examples    []string
	Run         func(args []string) error
}
