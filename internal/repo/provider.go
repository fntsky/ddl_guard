package repo

import (
	"github.com/fntsky/ddl_guard/internal/base/data"
	"github.com/fntsky/ddl_guard/internal/repo/ddl"
	"github.com/google/wire"
)

var ProviderSetRepo = wire.NewSet(
	data.NewDB,
	data.NewData,
	ddl.NewDDLRepo,
)
