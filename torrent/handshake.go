package torrent

import (
	"errors"
	"io"
)

type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

func (h *Handshake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49)
	buf[0] = byte(len(h.Pstr))
	curr := 1
	curr += copy(buf[curr:], h.Pstr)
	curr += copy(buf[curr:], make([]byte, 8))
	curr += copy(buf[curr:], h.InfoHash[:])
	curr += copy(buf[curr:], h.PeerID[:])
	return buf
}

func readerToHandshake(r io.Reader) (*Handshake, error) {
	length := [1]byte{}
	_, err := io.ReadFull(r, length[:])
	if err != nil {
		return nil, err
	}

	pstrlen := int(length[0])
	if pstrlen == 0 {
		return nil, errors.New("pstrlen should be greater than zero")
	}

	handshakeBuff := make([]byte, 48+pstrlen)
	_, err = io.ReadFull(r, handshakeBuff)
	if err != nil {
		return nil, err
	}

	var infoHash, peerID [20]byte

	copy(infoHash[:], handshakeBuff[pstrlen+8:pstrlen+8+20])
	copy(peerID[:], handshakeBuff[pstrlen+8+20:])

	h := Handshake{
		Pstr:     string(handshakeBuff[0:pstrlen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}

	return &h, nil
}
