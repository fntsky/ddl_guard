package cli

import (
	"fmt"

	"github.com/fntsky/ddl_guard/configs"
	"github.com/fntsky/ddl_guard/internal/base/path"
	"github.com/fntsky/ddl_guard/pkg/dir"
	"github.com/fntsky/ddl_guard/pkg/writer"
)

func Install(dataDirPath string) {
	path.FormatAllPath(dataDirPath)
	if err := InstallConfig(""); err != nil {
		fmt.Printf("[Install] failed to install config: %v\n", err)
		return
	}
}

func InstallConfig(configFilePath string) error {
	if len(configFilePath) == 0 {
		configFilePath = path.ConfigFileDir + path.DefaultConfigFileName
	}
	if CheckConfig(configFilePath) {
		fmt.Printf("[InstallConfig] config file already exists at %s\n", configFilePath)
		return nil
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
	return nil
}

func CheckConfig(configFilePath string) bool {
	return dir.CheckFileExist(configFilePath)
}
