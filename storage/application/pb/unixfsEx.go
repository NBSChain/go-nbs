package unixfs_pb

//TODO::change syntax = "proto2"; -->>syntax = "proto3";
import "github.com/golang/protobuf/proto"

func (format *Data) AddBlockSize(size int64) {

	format.Filesize = proto.Uint64(uint64(
		int64(format.GetFilesize()) + size))

	format.Blocksizes = append(format.Blocksizes, uint64(size))
}

func FromBytes(data []byte) (*Data, error) {
	pbdata := new(Data)
	err := proto.Unmarshal(data, pbdata)
	if err != nil {
		return nil, err
	}
	return pbdata, nil
}

func FilePBData(data []byte, totalsize uint64) []byte {
	pbfile := new(Data)
	typ := Data_File
	pbfile.Type = &typ
	pbfile.Data = data
	pbfile.Filesize = proto.Uint64(totalsize)

	data, err := proto.Marshal(pbfile)
	if err != nil {
		panic(err)
	}
	return data
}

func FolderPBData() []byte {
	pbfile := new(Data)
	typ := Data_Directory
	pbfile.Type = &typ

	data, err := proto.Marshal(pbfile)
	if err != nil {
		panic(err)
	}
	return data
}
