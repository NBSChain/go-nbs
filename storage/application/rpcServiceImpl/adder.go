package rpcServiceImpl

import (
	"errors"
	"github.com/NBSChain/go-nbs/console/pb"
	"github.com/NBSChain/go-nbs/storage/application/pb"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"github.com/NBSChain/go-nbs/storage/merkledag/ipld"
	"github.com/golang/protobuf/proto"
	"io"
	"strconv"
)

/*****************************************************************
*
*		Adder
*
*****************************************************************/

type Adder struct {
	tempRoot *cid.Cid
	rootDir  *Directory
	importer FileImporter
	nextData []byte
	batch    *Batch
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
		return nil, errors.New("don't build empty node. ")
	}

	root := adder.newImportNode(unixfs_pb.Data_File)
	fileSize, err := adder.leafNodeWithData(root)
	if err != nil {
		return nil, err
	}

	logger.Info("start leaf node->", root.dag.String())

	for depth := 1; adder.hasNext(); depth++ {

		newRoot := adder.newImportNode(unixfs_pb.Data_File)

		newRoot.AddChild(adder, root, fileSize)

		fileSize, err = adder.fillNodeRec(newRoot, depth)
		if err != nil {
			return nil, err
		}

		root = newRoot

		logger.Info("root->", root.dag.String())
	}

	return adder.AddNodeAndClose(root)
}

func (adder *Adder) newImportNode(nodeType unixfs_pb.Data_DataType) *DagDataBridge {

	node := new(DagDataBridge)

	node.dag = ipld.NewNode()

	node.format = &unixfs_pb.Data{
		Type: &nodeType,
	}

	return node
}

func (adder *Adder) leafNodeWithData(node *DagDataBridge) (int64, error) {

	data := adder.nextData

	defer func() {
		adder.nextData = nil
	}()

	dataLen := int64(len(data))
	node.format.Filesize = proto.Uint64(uint64(len(data)))
	node.format.Data = data

	err := node.Commit()

	if err != nil {
		return 0, err
	}

	return dataLen, nil
}

func (adder *Adder) fillNodeRec(node *DagDataBridge, depth int) (int64, error) {

	if depth < 1 {
		return 0, errors.New("attempt to fillNode at depth < 1")
	}

	var childFileSize int64
	var err error

	for node.NumChildren() < DefaultLinksPerBlock && adder.hasNext() {

		childNode := adder.newImportNode(unixfs_pb.Data_File)

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

		logger.Debug("===3===depth: ", depth, " childNode->", childNode.dag.String())
	}

	resultFileSize := node.FileSize()
	err = node.Commit()

	if err != nil {
		return 0, err
	}

	return resultFileSize, nil
}

func (adder *Adder) AddNodeAndClose(node *DagDataBridge) (ipld.DagNode, error) {

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

func (adder *Adder) AddNode(node ipld.DagNode, path string) error {

	if err := adder.rootDir.PutNode(path, node); err != nil {
		return err
	}

	go adder.OutputDagNode(path, node)

	return nil
}

func (adder *Adder) OutputDagNode(name string, dagNode ipld.DagNode) error {

	s, err := dagNode.Size()

	if err != nil {
		return err
	}
	//TODO:: Add process tips.
	adder.importer.ResultCh() <- &pb.AddResponse{
		Hash: dagNode.String(),
		Name: name,
		Size: strconv.FormatInt(s, 10),
	}

	return nil
}

func (adder *Adder) Finalize() (ipld.DagNode, error) {

	//TODO:: Maybe need to flush-->err := adder.rootDir.Flush()

	//TODO:: err = adder.outputDirs(name, root)

	//TODO:: republish the node to network.

	node := adder.rootDir.dirIO.(ipld.DagNode)
	return node, nil
}

func (adder *Adder) PinRoot() error {
	//TODO:: pin it.
	return nil
}
