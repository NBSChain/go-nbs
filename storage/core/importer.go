package core

import (
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"io"
	"os"
	"sync"
	"time"
)

type FileImporter interface {
	io.Closer

	NextChunk() ([]byte, error)

	FileName() string

	FullPath() string

	IsDirectory() bool

	NextFile() (FileImporter, error)
}

type File struct {
	name      string
	node      ipld.DagNode
	RawLeaves bool
}

type Directory struct {
	name      string
	lock      sync.Mutex
	modTime   time.Time
	childDirs map[string]*Directory
	files     map[string]*File
}

type Adder struct {
	localStore BlockStore
	rootNode   ipld.DagNode
	tempRoot   *cid.Cid
	rootDir    *Directory
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
