package utils

import (
	"os"
)

func PathExists(path string) (bool, os.FileInfo) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}

	return true, info
}

func DirExists(dirname string) bool {
	exists, info := PathExists(dirname)
	if !exists {
		return false
	}

	return info.IsDir()
}
func FileExists(filename string) bool {
	exists, info := PathExists(filename)
	if !exists {
		return false
	}

	return !info.IsDir()
}
