package migrations

import (
	"context"
	"fmt"

	"github.com/fntsky/ddl_guard/internal/entity"
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

func Migrate(ctx context.Context, engine *xorm.Engine) error {
	currentVersion, err := GetCurrentVersion(ctx, engine)
	if err != nil {
		return fmt.Errorf("get current version failed: %w", err)
	}
	if currentVersion < minDBVersion+int64(len(migrations)) {
		fmt.Printf("Now update database to next version: %s\n", migrations[currentVersion+1].Version())
		fmt.Printf("Description: %s\n", migrations[currentVersion+1].Description())
		if err := migrations[currentVersion+1].Migrate(ctx, engine); err != nil {
			return fmt.Errorf("migrate to version %s failed: %w", migrations[currentVersion+1].Version(), err)
		}
		currentVersion++
		// Update version in database
		if _, err := engine.Context(ctx).ID(1).Update(&entity.Version{VersionNumber: currentVersion}); err != nil {
			return fmt.Errorf("update version to %d failed: %w", currentVersion, err)
		}
		fmt.Printf("Database version updated to %d successfully\n", currentVersion)
	}
	return nil
}

var migrations = []Migration{
	NewMigration("0.0.1", "this is first version", nil),
}

func ExpectVersion() int64 {
	return minDBVersion + int64(len(migrations))
}

func GetCurrentVersion(ctx context.Context, engine *xorm.Engine) (int64, error) {
	if err := engine.Context(ctx).Sync2(&entity.Version{}); err != nil {

		return -1, fmt.Errorf("sync version fail: %v", err)
	}
	var version, err = GetDBVersion(ctx, engine)
	if err != nil {
		return -1, fmt.Errorf("get current version fail: %v", err)
	}
	return version, nil
}

func GetDBVersion(ctx context.Context, engine *xorm.Engine) (int64, error) {
	if err := engine.Context(ctx).Sync2(&entity.Version{}); err != nil {

		return -1, fmt.Errorf("sync version fail: %v", err)
	}
	currentVersion := &entity.Version{ID: 1}
	has, err := engine.Context(ctx).Get(currentVersion)
	if err != nil {
		return -1, fmt.Errorf("get version fail: %v", err)
	}
	if !has {
		return -1, fmt.Errorf("version record not found")
	}
	return currentVersion.VersionNumber, nil

}
