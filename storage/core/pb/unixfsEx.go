package unixfs_pb

func (format *Data) AddBlockSize(size int64) {

	format.Filesize += size

	format.Blocksizes = append(format.Blocksizes, size)
}
