package torrent_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/juneira/torrent-study/torrent"
)

func TestFormatRequest(t *testing.T) {
	result := torrent.FormatRequest(1, 2, 3)
	expected := &torrent.Message{ID: torrent.MsgRequest, Payload: []byte{0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 3}}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("result: %v, expected: %v", result, expected)
	}
}

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

func TestMessageParsePieceIndex(t *testing.T) {
	payload := []byte{0, 0, 0, 23, 0, 0, 0, 4, 1, 2, 3, 4, 5}
	m := torrent.Message{ID: torrent.MsgPiece, Payload: payload}

	result, err := m.ParsePieceIndex()
	if err != nil {
		t.Fatal(err)
	}

	expected := 23

	if result != expected {
		t.Errorf("result: %d, expected: %d", result, expected)
	}
}

func TestMessageParsePiece(t *testing.T) {
	payload := []byte{0, 0, 0, 23, 0, 0, 0, 4, 1, 2, 3, 4, 5}
	m := torrent.Message{ID: torrent.MsgPiece, Payload: payload}

	var resultPiece [9]byte

	resultLen, err := m.ParsePiece(23, resultPiece[:])
	if err != nil {
		t.Fatal(err)
	}

	expectedLen := 5
	expectedPiece := []byte{0, 0, 0, 0, 1, 2, 3, 4, 5}

	if resultLen != expectedLen {
		t.Errorf("result: %d, expected: %d", resultLen, expectedLen)
	}

	if !bytes.Equal(resultPiece[:], expectedPiece[:]) {
		t.Errorf("result: %v, expected: %v", resultPiece, expectedPiece)
	}
}
