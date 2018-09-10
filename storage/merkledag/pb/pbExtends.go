package merkledag_pb

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
	if m.Name != "" {
		l = len(m.Name)
		n += 1 + l + sovMerkledag(uint64(l))
	}
	if m.Tsize != 0 {
		n += 1 + sovMerkledag(uint64(m.Tsize))
	}
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
	if m.Name != "" {
		data[i] = 0x12
		i++
		i = encodeVarintMerkledag(data, i, uint64(len(m.Name)))
		i += copy(data[i:], m.Name)
	}
	if m.Tsize != 0 {
		data[i] = 0x18
		i++
		i = encodeVarintMerkledag(data, i, uint64(m.Tsize))
	}
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
