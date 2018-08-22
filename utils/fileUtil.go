package utils

import "os"

func FileExists(fileName string) (os.FileInfo, bool) {

	fileInfo, err := os.Lstat(fileName)

	if fileInfo != nil || (err != nil && !os.IsNotExist(err)) {
		return fileInfo, true
	}

	return nil, false
}
