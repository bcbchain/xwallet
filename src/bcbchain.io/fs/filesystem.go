package fs

import (
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func MakeDir(dir string) (bool, error) {
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
