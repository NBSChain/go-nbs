package core

import "io"

type FileImporter interface {
	io.Closer

	NextChunk() (chunk []byte, err error)

	FileName() string

	FullPath() string

	IsDirectory() bool

	NextFile() (FileImporter, error)
}

func ImportFile(importer FileImporter) error {
	return nil
}
