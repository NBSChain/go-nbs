package core

import (
	"io"
	"os"
)

type FileImporter interface {
	io.Closer

	NextChunk() (chunk []byte, err error)

	FileName() string

	FullPath() string

	IsDirectory() bool

	NextFile() (FileImporter, error)
}

func ImportFile(importer FileImporter) error {

	logger.Info("\nfileName:" + importer.FileName())
	logger.Info("\nfullPath:" + importer.FullPath())

	file, err := os.Create("server-" + importer.FileName())
	if err != nil {
		return err
	}

	data, err := importer.NextChunk()
	for err == nil {
		file.Write(data)
	}

	if err != io.EOF {
		return err
	}

	file.Close()
	importer.Close()

	return nil
}
