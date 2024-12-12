package torrent

import (
	"bytes"
	"crypto/sha1"
	"errors"
)

const LengthMax = 16384

type pieceState uint8

const (
	PiecePending  pieceState = 0
	PieceFinished pieceState = 1
	PieceError    pieceState = 2
)

type Piece struct {
	Index   int
	Hash    [20]byte
	Data    []byte
	Begin   int
	Waiting bool
	Status  pieceState
}

func (p *Piece) CheckFinished() error {
	if p.Begin == len(p.Data) {
		if p.checkSum() {
			p.Status = PieceFinished
			return nil
		}

		p.Status = PieceError

		return errors.New("invalid checksum")
	}

	return nil
}

func (p *Piece) checkSum() bool {
	check := sha1.Sum(p.Data[:])
	return bytes.Equal(p.Hash[:], check[:])
}
