package env

import "os"

func GetRootPath() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path
}
