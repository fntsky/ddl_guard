package writer

import "os"

func WriterFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
