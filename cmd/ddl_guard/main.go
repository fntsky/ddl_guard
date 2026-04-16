package main

import (
	ddlcmd "github.com/fntsky/ddl_guard/cmd"
	_ "github.com/fntsky/ddl_guard/docs"
)

// @title           DDL Guard API
// @version         1.0
// @description     This is a DDL Guard server API.
// @termsOfService  http://swagger.io/terms/
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
// @BasePath  /api/v1
func main() {
	ddlcmd.Main()
}
