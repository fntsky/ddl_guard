package cli

import (
	"fmt"
	"path/filepath"

	"github.com/fntsky/ddl_guard/configs"
	"github.com/fntsky/ddl_guard/internal/base/conf"
	"github.com/fntsky/ddl_guard/internal/base/data"
	"github.com/fntsky/ddl_guard/internal/base/path"
	"github.com/fntsky/ddl_guard/internal/migrations"
	"github.com/fntsky/ddl_guard/pkg/dir"
	"github.com/fntsky/ddl_guard/pkg/writer"
)

func Install(dataDirPath string) {
	path.FormatAllPath(dataDirPath)
	if err := InstallConfig(""); err != nil {
		fmt.Printf("[Install] failed to install config: %v\n", err)
		return
	}
	InstallDB()
}

func InstallConfig(configFilePath string) error {
	if len(configFilePath) == 0 {
		configFilePath = filepath.Join(path.ConfigFileDir, path.DefaultConfigFileName)
	}
	if CheckConfig(configFilePath) {
		fmt.Printf("[InstallConfig] config file already exists at %s\n", configFilePath)
	}
	if err := dir.CreateDirifNotExist(path.ConfigFileDir); err != nil {
		fmt.Printf("[InstallConfig] failed to create config directory: %v\n", err)
		return err
	}
	if err := writer.WriterFile(configFilePath, string(configs.ConfigYaml)); err != nil {
		fmt.Printf("[InstallConfig] failed to write config file: %v\n", err)
		return err
	}
	fmt.Printf("[InstallConfig] config file installed successfully at %s\n", configFilePath)
	if _, err := conf.LoadGlobal(configFilePath); err != nil {
		fmt.Printf("[InstallConfig] failed to load global config: %v\n", err)
		return err
	}
	return nil
}

func InstallDB() {
	db, err := data.NewDB(true)
	if err != nil {
		fmt.Printf("[InstallDB] failed to create DB engine: %v\n", err)
		return
	}
	m := migrations.NewMentor(db)
	if err := m.InitDB(); err != nil {
		fmt.Printf("[InstallDB] failed to initialize database: %v\n", err)
		return
	}
	fmt.Println("[InstallDB] database initialized successfully")

}

func CheckConfig(configFilePath string) bool {
	return dir.CheckFileExist(configFilePath)
}
