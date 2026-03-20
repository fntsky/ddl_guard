package service

import (
	"github.com/fntsky/ddl_guard/internal/service/ddl"
	"github.com/google/wire"
)

var ProviderSetService = wire.NewSet(
	ddl.NewDDLService,
)
