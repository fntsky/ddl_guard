package service

import (
	ai "github.com/fntsky/ddl_guard/internal/service/ai"
	"github.com/fntsky/ddl_guard/internal/service/ddl"
	"github.com/google/wire"
)

var ProviderSetService = wire.NewSet(
	ai.NewAIProvider,
	ddl.NewDDLService,
)
