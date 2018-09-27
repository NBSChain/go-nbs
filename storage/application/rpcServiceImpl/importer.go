package rpcServiceImpl

import (
	"github.com/NBSChain/go-nbs/storage/application/pb"
	"github.com/NBSChain/go-nbs/storage/merkledag"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/NBSChain/go-nbs/utils/cmdKits/pb"
	"github.com/gogo/protobuf/proto"
	"io"
)

const BlockSizeLimit = 1048576 // 1 MB

const roughLinkBlockSize = 1 << 13 // 8KB
const roughLinkSize = 34 + 8 + 5   // sha256 multihash + size + no name + protobuf framing
const adderOutChanSize = 8
const DefaultLinksPerBlock = roughLinkBlockSize / roughLinkSize

var logger = utils.GetLogInstance()

type FileImporter interface {
	io.Closer

	NextChunk() ([]byte, error)

	FileName() string

	FullPath() string

	IsDirectory() bool

	NextFile() (FileImporter, error)

	ResultCh() chan *pb.AddResponse
}

//TODO:: add args and optional settings.
func ImportFile(importer FileImporter) error {

	adder := &Adder{
		importer: importer,
		batch:    merkledag.NewBatch(),
		rootDir:  NewDir(),
	}

	rootNode, err := adder.buildNodeLayout()
	if err != nil {
		return err
	}

	logger.Info("currentNode:->", rootNode.String())

	adder.AddNode(rootNode, importer.FileName())

	adder.Finalize()

	adder.PinRoot()

	importer.Close()

	return nil
}

/*****************************************************************
*
*		DagDataBridge
*
*****************************************************************/

type DagDataBridge struct {
	dag    *ipld.ProtoDagNode
	format *unixfs_pb.Data
}

func (node *DagDataBridge) AddChild(adder *Adder, child *DagDataBridge, dataSize int64) error {

	err := node.dag.AddNodeLink("", child.dag)

	if err != nil {
		return err
	}

	logger.Debug("===3=== newRoot->", node.dag.String())

	node.format.AddBlockSize(dataSize)

	logger.Debug("===4=== newRoot->", node.dag.String())

	return adder.batch.Add(child.dag)
}

func (node *DagDataBridge) NumChildren() int {
	return len(node.format.Blocksizes)
}

func (node *DagDataBridge) FileSize() int64 {
	return int64(node.format.GetFilesize())
}

func (node *DagDataBridge) Commit() error {

	fileData, err := proto.Marshal(node.format)
	if err != nil {
		return err
	}

	node.dag.SetData(fileData)

	return nil
}

func (node *DagDataBridge) Type() unixfs_pb.Data_DataType {
	return node.format.GetType()
}
