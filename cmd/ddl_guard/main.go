package main

import (
	ddlcmd "github.com/fntsky/ddl_guard/cmd"
	_ "github.com/fntsky/ddl_guard/docs"
)

// @title           DDL Guard API
// @version         1.0
// @description     This is a DDL Guard server API.
// @termsOfService  http://swagger.io/terms/
// @BasePath  /api/v1
func main() {
	ddlcmd.Main()
}
