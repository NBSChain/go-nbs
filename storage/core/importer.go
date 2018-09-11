package core

import (
	"github.com/NBSChain/go-nbs/storage/core/pb"
	"github.com/NBSChain/go-nbs/storage/merkledag"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
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
		batch:    merkledag.NewBatch(),
	}

	rootNode, err := adder.buildNodeLayout()
	if err != nil {
		return err
	}

	cidStr := rootNode.String()

	logger.Info("rootNode:->", cidStr)

	importer.Close()

	return nil
}

/*****************************************************************
*
*		Adder
*
*****************************************************************/

type Adder struct {
	rootNode ipld.DagNode
	tempRoot *cid.Cid
	rootDir  *Directory
	importer FileImporter
	nextData []byte
	batch    *merkledag.Batch
}

/*******************************************************************************
*                      +-------------+
*                      |   `node`    |
*                      |  (new root) |
*                      +-------------+
*                            |
*              +-------------+ - - - - - - + - - - - - - - - - - - +
*              |                           |                       |
*      +--------------+             + - - - - -  +           + - - - - -  +
*      |  (old root)  |             |  new child |           |            |
*      +--------------+             + - - - - -  +           + - - - - -  +
*              |                          |                        |
*       +------+------+             + - - + - - - +
*       |             |             |             |
*  +=========+   +=========+   + - - - - +    + - - - - +
*  | Chunk 1 |   | Chunk 2 |   | Chunk 3 |    | Chunk 4 |
*  +=========+   +=========+   + - - - - +    + - - - - +
*
*******************************************************************************/
func (adder *Adder) hasNext() bool {

	if adder.nextData != nil {
		return true
	}

	data, err := adder.importer.NextChunk()
	if err != nil {
		adder.nextData = nil
		return false
	}

	dataLen := len(data)
	if dataLen > BlockSizeLimit {
		logger.Error("object size limit exceeded")
		return false
	}

	adder.nextData = data

	return true
}

func (adder *Adder) buildNodeLayout() (ipld.DagNode, error) {

	if !adder.hasNext() {
		return nil, errors.New("don't build empty node.")
	}

	root := adder.newImportNode(TFile)
	fileSize, err := adder.leafNodeWithData(root)
	if err != nil {
		return nil, err
	}

	logger.Info("start leaf node->", root.dag.String())

	for depth := 1; adder.hasNext(); depth++ {

		newRoot := adder.newImportNode(TFile)

		logger.Info("===1===depth: ", depth, " newRoot->", newRoot.dag.String())

		newRoot.AddChild(adder, root, fileSize)

		logger.Info("===2===depth: ", depth, " newRoot->", newRoot.dag.String())

		fileSize, err = adder.fillNodeRec(newRoot, depth)
		if err != nil {
			return nil, err
		}

		root = newRoot

		logger.Info("root->", root.dag.String())
	}

	return adder.AddNodeAndClose(root)
}

func (adder *Adder) newImportNode(nodeType unixfs_pb.Data_DataType) *ImportNode {

	node := new(ImportNode)

	node.dag = new(ipld.ProtoDagNode)

	node.format = &unixfs_pb.Data{
		Type: &nodeType,
	}

	node.format.Filesize = proto.Uint64(node.format.GetFilesize())

	return node
}

func (adder *Adder) leafNodeWithData(node *ImportNode) (int64, error) {

	data := adder.nextData

	defer func() {
		adder.nextData = nil
	}()

	dataLen := int64(len(data))
	node.format.Filesize = proto.Uint64(uint64(
		int64(node.format.GetFilesize()) + int64(len(data)-len(node.format.GetData()))))
	node.format.Data = data

	err := node.Commit()

	if err != nil {
		return 0, err
	}

	return dataLen, nil
}

func (adder *Adder) fillNodeRec(node *ImportNode, depth int) (int64, error) {

	if depth < 1 {
		return 0, errors.New("attempt to fillNode at depth < 1")
	}

	var childFileSize int64
	var err error

	for node.NumChildren() < DefaultLinksPerBlock && adder.hasNext() {

		childNode := adder.newImportNode(TFile)

		if depth == 1 {
			childFileSize, err = adder.leafNodeWithData(childNode)
			if err != nil {
				if err == io.EOF {
					break
				} else {
					return 0, err
				}
			}
		} else {
			childFileSize, err = adder.fillNodeRec(childNode, depth-1)
			if err != nil {
				return 0, err
			}
		}

		err = node.AddChild(adder, childNode, childFileSize)
		if err != nil {
			return 0, err
		}

		err = adder.batch.Add(childNode.dag)
		if err != nil {
			return 0, err
		}
	}

	resultFileSize := node.FileSize()
	err = node.Commit()

	if err != nil {
		return 0, err
	}

	return resultFileSize, nil
}

func (adder *Adder) AddNodeAndClose(node *ImportNode) (ipld.DagNode, error) {

	dagNode := node.dag

	err := adder.batch.Add(dagNode)
	if err != nil {
		return nil, err
	}

	err = adder.batch.Commit()
	if err != nil {
		return nil, err
	}

	return dagNode, nil
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
