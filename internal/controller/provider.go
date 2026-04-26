package controller

import "github.com/google/wire"

var ProviderSetController = wire.NewSet(
	NewAuthController,
	NewDDLController,
	NewExamController,
	NewUserController,
)
