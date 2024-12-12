package torrent_test

import (
	"crypto/sha1"
	"testing"

	"github.com/juneira/torrent-study/torrent"
)

func TestPieceFinished(t *testing.T) {
	data := [5]byte{1, 2, 3, 4, 5}
	hash := sha1.Sum(data[:])
	p := torrent.Piece{Status: torrent.PiecePending, Hash: hash, Data: data[:], Begin: 5}

	p.CheckFinished()
	result := p.Status
	expected := torrent.PieceFinished

	if result != expected {
		t.Errorf("result: %d, expected: %d", result, expected)
	}
}
