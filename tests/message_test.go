package torrent_test

import (
	"bytes"
	"testing"

	"github.com/juneira/torrent-study/torrent"
)

func TestMessageSerialize(t *testing.T) {
	payload := []byte{1, 2, 3, 4}
	id := torrent.MsgPiece

	m := torrent.Message{ID: id, Payload: payload}
	result := m.Serialize()
	expected := []byte{0, 0, 0, 5, 7, 1, 2, 3, 4}

	if !bytes.Equal(result, expected) {
		t.Errorf("result: %v, expected: %v", result, expected)
	}
}
