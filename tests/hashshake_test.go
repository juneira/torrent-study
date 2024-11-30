package torrent_test

import (
	"reflect"
	"testing"

	"github.com/juneira/torrent-study/torrent"
)

func TestHandshakeSerialize(t *testing.T) {
	infoHash := [20]byte{0, 1, 2, 3, 4, 5, 6, 7, 8}
	peerID := [20]byte{0, 1, 2, 3, 4, 5, 6, 7, 8}
	h := torrent.Handshake{Pstr: "BitTorrent protocol", InfoHash: infoHash, PeerID: peerID}

	result := h.Serialize()
	expected := []byte{19, 66, 105, 116, 84, 111, 114, 114, 101, 110, 116, 32, 112, 114, 111, 116, 111, 99, 111, 108, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("result: %v, expected: %v", result, expected)
	}
}
