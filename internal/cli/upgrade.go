package cli

import (
	"fmt"
	"path/filepath"

	"github.com/fntsky/ddl_guard/internal/base/conf"
	"github.com/fntsky/ddl_guard/internal/base/data"
	"github.com/fntsky/ddl_guard/internal/base/path"
	"github.com/fntsky/ddl_guard/internal/migrations"
)

func UpgradeDB(configFilePath string) error {
	if len(configFilePath) == 0 {
		configFilePath = filepath.Join(path.ConfigFileDir, path.DefaultConfigFileName)
	}
	if _, err := conf.LoadGlobal(configFilePath); err != nil {
		return fmt.Errorf("load global config failed: %w", err)
	}

	db, err := data.NewDB(true)
	if err != nil {
		return fmt.Errorf("create DB engine failed: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	if err := migrations.Migrate(db); err != nil {
		return fmt.Errorf("migrate database failed: %w", err)
	}
	return nil
}
