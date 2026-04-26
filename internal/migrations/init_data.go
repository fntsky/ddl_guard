package migrations

import "github.com/fntsky/ddl_guard/internal/entity"

var (
	tables = []any{
		&entity.User{},
		&entity.UserAuth{},
		&entity.UserSession{},
		&entity.DDL{},
		&entity.Exam{},
		&entity.Version{},
	}
)
