package dir

import "os"

func CheckExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func CheckFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func CheckDirExist(dirPath string) bool {
	_, err := os.Stat(dirPath)
	return !os.IsNotExist(err)
}

func CreateDir(dirPath string) error {
	return os.MkdirAll(dirPath, os.ModePerm)
}

func CreateDirifNotExist(dirPath string) error {
	if CheckDirExist(dirPath) {
		return nil
	}
	return CreateDir(dirPath)
}
