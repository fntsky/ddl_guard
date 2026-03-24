package migrations

import (
	"context"

	"xorm.io/xorm"
)

const minDBVersion int64 = 0

type Migration interface {
	Version() string
	Description() string
	Migrate(ctx context.Context, x *xorm.Engine) error
}

type migration struct {
	version     string
	description string
	migrate     func(ctx context.Context, x *xorm.Engine) error
}

func (m *migration) Version() string {
	return m.version
}

func (m *migration) Description() string {
	return m.description
}

func (m *migration) Migrate(ctx context.Context, x *xorm.Engine) error {
	return m.migrate(ctx, x)
}

func NewMigration(version, desc string, fn func(ctx context.Context, x *xorm.Engine) error) Migration {
	return &migration{
		version:     version,
		description: desc,
		migrate:     fn,
	}
}

func Migrate(engine *xorm.Engine) error {
	return nil
}

var migrations = []Migration{
	NewMigration("0.0.1", "this is first version", nil),
}

func ExpectVersion() int64 {
	return minDBVersion + int64(len(migrations))
}
