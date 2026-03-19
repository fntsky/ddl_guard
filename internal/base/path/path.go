package path

import "sync"

const (
	DefaultConfigFileName = "config.yaml"
)

var (
	ConfigFileDir     = "/conf/"
	formatAllPathOnce sync.Once
)

func FormatAllPath(dataDirPath string) {
	formatAllPathOnce.Do(func() {
		ConfigFileDir = dataDirPath + ConfigFileDir
	})
}
