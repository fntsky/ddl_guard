package path

import (
	"path/filepath"
	"sync"
)

const (
	DefaultConfigFileName = "config.yaml"
)

var (
	ConfigFileDir     = "conf"
	formatAllPathOnce sync.Once
)

func FormatAllPath(dataDirPath string) {
	formatAllPathOnce.Do(func() {
		ConfigFileDir = filepath.Join(dataDirPath, ConfigFileDir)
	})
}
