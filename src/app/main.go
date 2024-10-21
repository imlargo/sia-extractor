package app

import (
	"fmt"
	"os"
)

const (
	comPaths     string = "paths"
	comGroup     string = "group"
	comElectivas string = "electivas"
	comDeploy    string = "deploy"
	comTest      string = "test"
	comExtract   string = "extract"
)

func Start() {

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Debe ingresar los argumentos")
		return
	}

	comando := args[0]

	app := CreateNewApp()
	app.AddHandler(comPaths, handlePaths)
	app.AddHandler(comGroup, handleGroup)
	app.AddHandler(comElectivas, handleElectivas)
	app.AddHandler(comDeploy, handleDeploy)
	app.AddHandler(comTest, handleTest)
	app.AddHandler(comExtract, handleExtract)

	app.Handle(comando, args)
}
