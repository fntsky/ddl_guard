package ddlcmd

func Main() {
	Execute()
}

func runApp() {
	app, cleanup, err := initApplication(true)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	if err := app.HttpServer.Run(":8080"); err != nil {
		panic(err)
	}
}
