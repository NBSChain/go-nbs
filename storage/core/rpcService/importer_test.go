package rpcService

import (
	"io"
	"os"
	"testing"
)

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

//TODO:: test more file size and more rpc chunk size.
func TestImportFile(t *testing.T) {

}
