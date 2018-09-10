package core

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/core/pb"
	"github.com/NBSChain/go-nbs/storage/merkledag"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/gogo/protobuf/proto"
	"io"
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

//TODO:: add args and optional settings.
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

	root := adder.newImportNode(TFile)

	fileSize, err := adder.fullFillLeafNode(root)

	for depth := 1; err == nil; depth++ {

		newRoot := adder.newImportNode(TFile)
		newRoot.AddChild(root, fileSize)
	}

	return adder.AddNodeAndClose(root)
}

func (adder *Adder) newImportNode(nodeType unixfs_pb.Data_DataType) *ImportNode {
	node := new(ImportNode)

	node.dag = new(ipld.ProtoDagNode)

	node.format = &unixfs_pb.Data{

		Type:     nodeType,
		Filesize: 0,
		DataLen:  0,
	}

	return node
}

func (adder *Adder) fullFillLeafNode(node *ImportNode) (int64, error) {

	data, err := adder.importer.NextChunk()
	if err != nil {
		return 0, err
	}

	dataLen := int64(len(data))
	if dataLen > BlockSizeLimit {
		return 0, fmt.Errorf("object size limit exceeded")
	}

	node.format.Data = data
	node.format.Filesize = dataLen
	node.format.DataLen = dataLen

	fileData, err := proto.Marshal(node.format)
	if err != nil {
		return 0, err
	}

	node.dag.SetData(fileData)

	return dataLen, nil
}

func (node *ImportNode) AddChild(child *ImportNode, dataSize int64) error {

	err := node.dag.AddNodeLink("", child.dag)
	if err != nil {
		return err
	}

	node.format.AddBlockSize(dataSize)

	return node.batch.Add(child.dag)
}

func (adder *Adder) AddNodeAndClose(node *ImportNode) (ipld.DagNode, error) {
	return node.dag, nil
}

type ImportNode struct {
	batch  *merkledag.Batch
	dag    *ipld.ProtoDagNode
	format *unixfs_pb.Data
}
