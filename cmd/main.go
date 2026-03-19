package ddlcmd

import "github.com/fntsky/ddl_guard/internal/base/conf"

func Main() {
	Execute()
}

func runApp(configPath string) {
	_, err := conf.LoadGlobal(configPath)
	if err != nil {
		panic(err)
	}
	app, cleanup, err := initApplication(true)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	if err := app.HttpServer.Run(conf.Global().Server.HTTP.Addr); err != nil {
		panic(err)
	}
}
