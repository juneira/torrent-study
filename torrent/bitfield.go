package torrent

type Bitfield []byte

func (bf Bitfield) HasPiece(index uint) bool {
	byteIndex := index / 8
	offset := index % 8
	if int(byteIndex) >= len(bf) {
		return false
	}

	return bf[int(byteIndex)]>>(7-offset)&1 != 0
}

func (bf Bitfield) SetPiece(index uint) {
	byteIndex := index / 8
	offset := index % 8

	if int(byteIndex) >= len(bf) {
		return
	}

	bf[byteIndex] |= (1 << (7 - offset))
}
