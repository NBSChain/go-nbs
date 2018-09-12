package rpcServiceImpl

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/core/pb"
	"github.com/NBSChain/go-nbs/storage/merkledag"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/pkg/errors"
	"path"
	"strings"
	"sync"
	"time"
)

type Directory struct {
	name      string
	lock      sync.Mutex
	modTime   time.Time
	childDirs map[string]*Directory
	files     map[string]*ipld.DagNode
}

func NewDir() *Directory {

	dagNode := ipld.NodeWithData(unixfs_pb.FolderPBData())

	return &Directory{
		name:      dagNode.String(),
		childDirs: make(map[string]*Directory),
		files:     make(map[string]*ipld.DagNode),
		modTime:   time.Now(),
	}
}

func (d *Directory) PutNode(filePath string, node ipld.DagNode) error {

	parentDir, fileName := path.Split(filePath)
	if fileName == "" {
		return errors.New("cannot create file with empty name")
	}

	dir, err := d.LookupDir(parentDir)
	if err != nil {
		return err
	}

	return dir.AddChild(fileName, node)
}

func (d *Directory) LookupDir(dirPath string) (*Directory, error) {

	dirPath = strings.Trim(dirPath, "/")

	parts := strings.Split(dirPath, "/")

	if len(parts) == 1 && parts[0] == "" {
		return d, nil
	} else {
		panic(" TODO:: we need consider the directory ")
		return nil, fmt.Errorf(" TODO:: we need consider the directory ")
	}
}

func (d *Directory) AddChild(fileName string, node ipld.DagNode) error {

	//TODO:: sync node cache

	dagService := merkledag.GetInstance()
	dagService.Add(node)

	//TODO:: AddUnixFSChild

	d.modTime = time.Now()

	return nil
}
