package unixfs_pb

import "github.com/golang/protobuf/proto"

func (format *Data) AddBlockSize(size int64) {

	format.Filesize = proto.Uint64(uint64(
		int64(format.GetFilesize()) + size))

	format.Blocksizes = append(format.Blocksizes, uint64(size))
}
