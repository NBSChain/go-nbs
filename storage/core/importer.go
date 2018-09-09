package core

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/core/pb"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/gogo/protobuf/proto"
	"io"
	"sync"
	"time"
)

const BlockSizeLimit = 1048576 // 1 MB

const roughLinkBlockSize = 1 << 13 // 8KB
const roughLinkSize = 34 + 8 + 5   // sha256 multihash + size + no name + protobuf framing

const DefaultLinksPerBlock = roughLinkBlockSize / roughLinkSize

const (
	TRaw       = unixfs_pb.Data_Raw
	TFile      = unixfs_pb.Data_File
	TDirectory = unixfs_pb.Data_Directory
	TMetadata  = unixfs_pb.Data_Metadata
	TSymlink   = unixfs_pb.Data_Symlink
	THAMTShard = unixfs_pb.Data_HAMTShard
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

	adder := &Adder{
		importer: importer,
	}

	_, err := adder.buildNodeLayout()
	if err != nil {
		return err
	}

	importer.Close()

	return nil
}

type Adder struct {
	rootNode ipld.DagNode
	tempRoot *cid.Cid
	rootDir  *Directory
	importer FileImporter
}

func (adder *Adder) buildNodeLayout() (ipld.DagNode, error) {

	root, fileSize, err := adder.newImportNode()

	for depth := 1; err == nil; depth++ {

		newRoot, _, _ := adder.newImportNode()
		newRoot.AddChild(root, fileSize)
	}

	return adder.AddNodeAndClose(root)
}

func (adder *Adder) newImportNode() (*ImportNode, int64, error) {

	data, err := adder.importer.NextChunk()
	if err != nil {
		return nil, 0, err
	}

	dataLen := int64(len(data))
	if dataLen > BlockSizeLimit {
		return nil, 0, fmt.Errorf("object size limit exceeded")
	}

	node := new(ImportNode)

	node.format = &unixfs_pb.Data{
		Type:     TFile,
		Data:     data,
		Filesize: dataLen,
		DataLen:  dataLen,
	}

	fileData, err := proto.Marshal(node.format)
	if err != nil {
		return nil, 0, err
	}

	node.dag = new(ipld.ProtoDagNode)
	node.dag.SetData(fileData)

	return node, dataLen, nil
}

func (adder *Adder) AddNodeAndClose(node *ImportNode) (ipld.DagNode, error) {
	return node.dag, nil
}

type ImportNode struct {
	dag    *ipld.ProtoDagNode
	format *unixfs_pb.Data
}

func (node *ImportNode) AddChild(child *ImportNode, dataSize int64) {

}

type Directory struct {
	name      string
	lock      sync.Mutex
	modTime   time.Time
	childDirs map[string]*Directory
	files     map[string]*ipld.DagNode
}
