package app

type App struct {
	Routes map[string]func(args []string)
}

func CreateNewApp() *App {
	return &App{
		Routes: make(map[string]func(args []string)),
	}
}

func (app *App) AddHandler(command string, handler func(args []string)) {
	app.Routes[command] = handler
}

func (app *App) Handle(command string, args []string) {
	if handler, ok := app.Routes[command]; ok {
		handler(args)
	} else {
		println("Comando no reconocido")
	}
}
