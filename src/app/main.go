package app

import (
	"fmt"
	"os"
)

func Start() {

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Debe ingresar los argumentos")
		return
	}

	comando := args[0]

	app := CreateNewApp()
	app.AddHandler("paths", handlePaths)
	app.AddHandler("group", handleGroup)
	app.AddHandler("electivas", handleElectivas)
	app.AddHandler("deploy", handleDeploy)
	app.AddHandler("test", handleTest)
	app.AddHandler("extract", handleExtract)

	app.Handle(comando, args)
}
