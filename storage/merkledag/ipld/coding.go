package ipld

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/merkledag/cid"
	"strings"
	"sync"
)


// DecodeBlockFunc functions decode blocks into nodes.
type DecodeBlockFunc func([]byte, *cid.Cid) (DagNode, error)

type BlockDecoder interface {
	Register(uint64, DecodeBlockFunc)
	Decode([]byte, *cid.Cid) (DagNode, error)
}
type safeBlockDecoder struct {
	// Can be replaced with an RCU if necessary.
	lock     sync.RWMutex
	decoders map[uint64]DecodeBlockFunc
}

// Register registers decoder for all blocks with the passed codec.
//
// This will silently replace any existing registered block decoders.
func (d *safeBlockDecoder) Register(code uint64, decoder DecodeBlockFunc) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.decoders[code] = decoder
}

func (d *safeBlockDecoder) Decode(block []byte, c *cid.Cid) (DagNode, error) {

	ty := c.Type()

	d.lock.RLock()
	decoder, ok := d.decoders[ty]
	d.lock.RUnlock()

	if ok {
		return decoder(block, c)
	} else {
		// TODO: get the *long* name for this format
		return nil, fmt.Errorf("unrecognized object type: %d", ty)
	}
}

var DefaultBlockDecoder BlockDecoder = &safeBlockDecoder{decoders: make(map[uint64]DecodeBlockFunc)}

// Decode decodes the given block using the default BlockDecoder.
func Decode(block []byte, c *cid.Cid) (DagNode, error) {
	return DefaultBlockDecoder.Decode(block, c)
}

// Register registers block decoders with the default BlockDecoder.
func Register(code uint64, decoder DecodeBlockFunc) {
	DefaultBlockDecoder.Register(code, decoder)
}

func DecodeProtobuf(encoded []byte) (*ProtoDagNode, error) {
	n := new(ProtoDagNode)
	err := n.unmarshal(encoded)
	if err != nil {
		return nil, fmt.Errorf("incorrectly formatted merkledag node: %s", err)
	}
	return n, nil
}

func DecodeProtobufBlock(b []byte, c *cid.Cid) (DagNode, error) {

	if c.Type() != cid.DagProtobuf {
		return nil, fmt.Errorf("this function can only decode protobuf nodes")
	}

	decnd, err := DecodeProtobuf(b)
	if err != nil {
		if strings.Contains(err.Error(), "Unmarshal failed") {
			return nil, fmt.Errorf("the block referred to by '%s' was not a valid merkledag node", c)
		}
		return nil, fmt.Errorf("failed to decode Protocol Buffers: %v", err)
	}

	decnd.cached = c
	return decnd, nil
}

