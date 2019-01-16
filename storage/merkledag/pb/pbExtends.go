package merkledag_pb

import (
	"fmt"
	"io"
)

func sovMerkledag(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}

func encodeVarintMerkledag(data []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		data[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	data[offset] = uint8(v)
	return offset + 1
}

func (m *PBLink) Size() (n int) {
	var l int
	_ = l
	if m.Hash != nil {
		l = len(m.Hash)
		n += 1 + l + sovMerkledag(uint64(l))
	}

	l = len(m.Name)
	n += 1 + l + sovMerkledag(uint64(l))

	n += 1 + sovMerkledag(uint64(m.Tsize))

	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}
func (m *PBNode) Size() (n int) {
	var l int
	_ = l
	if len(m.Links) > 0 {
		for _, e := range m.Links {
			l = e.Size()
			n += 1 + l + sovMerkledag(uint64(l))
		}
	}
	if m.Data != nil {
		l = len(m.Data)
		n += 1 + l + sovMerkledag(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *PBNode) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *PBLink) MarshalTo(data []byte) (n int, err error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Hash != nil {
		data[i] = 0xa
		i++
		i = encodeVarintMerkledag(data, i, uint64(len(m.Hash)))
		i += copy(data[i:], m.Hash)
	}

	data[i] = 0x12
	i++
	i = encodeVarintMerkledag(data, i, uint64(len(m.Name)))
	i += copy(data[i:], m.Name)

	data[i] = 0x18
	i++
	i = encodeVarintMerkledag(data, i, uint64(m.Tsize))

	if m.XXX_unrecognized != nil {
		i += copy(data[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *PBNode) MarshalTo(data []byte) (n int, err error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Links) > 0 {
		for _, msg := range m.Links {
			data[i] = 0x12
			i++
			i = encodeVarintMerkledag(data, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(data[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.Data != nil {
		data[i] = 0xa
		i++
		i = encodeVarintMerkledag(data, i, uint64(len(m.Data)))
		i += copy(data[i:], m.Data)
	}
	if m.XXX_unrecognized != nil {
		i += copy(data[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *PBLink) Unmarshal(data []byte) error {
	l := len(data)
	index := 0
	for index < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if index >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[index]
			index++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = append([]byte{}, data[index:postIndex]...)
			index = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + int(stringLen)
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(data[index:postIndex])
			index = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tsize", wireType)
			}
			var v uint64
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				v |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Tsize = int64(v)
		default:
			var sizeOfWire int
			for {
				sizeOfWire++
				wire >>= 7
				if wire == 0 {
					break
				}
			}
			index -= sizeOfWire
			skippy, err := Skip(data[index:])
			if err != nil {
				return err
			}
			if (index + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, data[index:index+skippy]...)
			index += skippy
		}
	}
	return nil
}

func (m *PBNode) Unmarshal(data []byte) error {
	l := len(data)
	index := 0
	for index < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if index >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[index]
			index++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		switch fieldNum {
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Links", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Links = append(m.Links, &PBLink{})
			m.Links[len(m.Links)-1].Unmarshal(data[index:postIndex])
			index = postIndex
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			postIndex := index + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append([]byte{}, data[index:postIndex]...)
			index = postIndex
		default:
			var sizeOfWire int
			for {
				sizeOfWire++
				wire >>= 7
				if wire == 0 {
					break
				}
			}
			index -= sizeOfWire
			skippy, err := Skip(data[index:])
			if err != nil {
				return err
			}
			if (index + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, data[index:index+skippy]...)
			index += skippy
		}
	}
	return nil
}

func Skip(data []byte) (n int, err error) {
	l := len(data)
	index := 0
	for index < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if index >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := data[index]
			index++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for {
				if index >= l {
					return 0, io.ErrUnexpectedEOF
				}
				index++
				if data[index-1] < 0x80 {
					break
				}
			}
			return index, nil
		case 1:
			index += 8
			return index, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if index >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := data[index]
				index++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			index += length
			return index, nil
		case 3:
			for {
				var innerWire uint64
				var start int = index
				for shift := uint(0); ; shift += 7 {
					if index >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := data[index]
					index++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := Skip(data[start:])
				if err != nil {
					return 0, err
				}
				index = start + next
			}
			return index, nil
		case 4:
			return index, nil
		case 5:
			index += 4
			return index, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}
