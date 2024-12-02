package torrent_test

import (
	"bytes"
	"testing"

	"github.com/juneira/torrent-study/torrent"
)

func TestBitfieldHasPiece(t *testing.T) {
	b := torrent.Bitfield([]byte{255, 2, 20})

	expected := true
	result := b.HasPiece(0)

	if result != expected {
		t.Errorf("result: %t, expected: %t", result, expected)
	}

	expected = false
	result = b.HasPiece(8)

	if result != expected {
		t.Errorf("result: %t, expected: %t", result, expected)
	}
}

func TestBitfieldSetPiece(t *testing.T) {
	b := torrent.Bitfield([]byte{255, 0, 20})

	b.SetPiece(8)

	expected := []byte{255, 128, 20}

	if !bytes.Equal([]byte(b), expected) {
		t.Errorf("result: %v, expected: %v", b, expected)
	}
}
