package file

import (
	"os"
)

func InitFileStorage(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModePerm)
}
