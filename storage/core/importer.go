package core

import (
	"io"
	"os"
)

type FileImporter interface {
	io.Closer

	NextChunk() ([]byte, error)

	FileName() string

	FullPath() string

	IsDirectory() bool

	NextFile() (FileImporter, error)
}

func ImportFile(importer FileImporter) error {

	importer.Close()

	return nil
}

func ImportFileTest(importer FileImporter) error {

	logger.Info("fileName:" + importer.FileName())
	logger.Info("fullPath:" + importer.FullPath())

	file, err := os.Create("server-" + importer.FileName())
	if err != nil {
		return err
	}

	data, err := importer.NextChunk()
	for err == nil {
		file.Write(data)
		data, err = importer.NextChunk()
	}

	if err != io.EOF {
		return err
	}

	file.Close()
	importer.Close()

	return nil
}
