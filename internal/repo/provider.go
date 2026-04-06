package repo

import (
	"github.com/fntsky/ddl_guard/internal/base/data"
	"github.com/fntsky/ddl_guard/internal/repo/ddl"
	"github.com/fntsky/ddl_guard/internal/repo/session"
	"github.com/fntsky/ddl_guard/internal/repo/user"
	"github.com/google/wire"
)

var ProviderSetRepo = wire.NewSet(
	data.NewDB,
	data.NewRedisClient,
	data.NewData,
	ddl.NewDDLRepo,
	session.NewSessionRepo,
	user.NewUserRepo,
)
