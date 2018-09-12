package rpcService

import (
	"github.com/NBSChain/go-nbs/storage/core/pb"
	"github.com/NBSChain/go-nbs/storage/merkledag"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/gogo/protobuf/proto"
	"io"
)

const BlockSizeLimit = 1048576 // 1 MB

const roughLinkBlockSize = 1 << 13 // 8KB
const roughLinkSize = 34 + 8 + 5   // sha256 multihash + size + no name + protobuf framing
const adderOutChanSize = 8
const DefaultLinksPerBlock = roughLinkBlockSize / roughLinkSize

const (
	TRaw       = unixfs_pb.Data_Raw
	TFile      = unixfs_pb.Data_File
	TDirectory = unixfs_pb.Data_Directory
	TMetadata  = unixfs_pb.Data_Metadata
	TSymlink   = unixfs_pb.Data_Symlink
	THAMTShard = unixfs_pb.Data_HAMTShard
)

var logger = utils.GetLogInstance()

type FileImporter interface {
	io.Closer

	NextChunk() ([]byte, error)

	FileName() string

	FullPath() string

	IsDirectory() bool

	NextFile() (FileImporter, error)
}

//TODO:: add args and optional settings.
func ImportFile(importer FileImporter) error {

	adder := &Adder{
		importer: importer,
		batch:    merkledag.NewBatch(),
		rootDir:  NewDir(),
		Out:      make(chan interface{}, adderOutChanSize),
	}

	rootNode, err := adder.buildNodeLayout()
	if err != nil {
		return err
	}

	logger.Info("rootNode:->", rootNode.String())

	adder.addNode(rootNode, importer.FileName())

	importer.Close()

	return nil
}

/*****************************************************************
*
*		ImportNode
*
*****************************************************************/

type ImportNode struct {
	dag    *ipld.ProtoDagNode
	format *unixfs_pb.Data
}

func (node *ImportNode) AddChild(adder *Adder, child *ImportNode, dataSize int64) error {

	err := node.dag.AddNodeLink("", child.dag)

	if err != nil {
		return err
	}

	logger.Info("===3=== newRoot->", node.dag.String())

	node.format.AddBlockSize(dataSize)

	logger.Info("===4=== newRoot->", node.dag.String())

	return adder.batch.Add(child.dag)
}

func (node *ImportNode) NumChildren() int {
	return len(node.format.Blocksizes)
}

func (node *ImportNode) FileSize() int64 {
	return int64(node.format.GetFilesize())
}

func (node *ImportNode) Commit() error {

	fileData, err := proto.Marshal(node.format)
	if err != nil {
		return err
	}

	node.dag.SetData(fileData)

	return nil
}
