package core

import "io"

type FileImporter interface {
	io.ReadCloser

	FileName() string

	FullPath() string

	IsDirectory() bool

	NextFile() (FileImporter, error)
}

func ImportFile(importer FileImporter) error {
	return nil
}
