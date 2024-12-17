package torrent_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/juneira/torrent-study/torrent"
)

type MockConnection struct {
	t               *testing.T
	expectedReceive []byte
	sendData        []byte
}

func (m *MockConnection) GetConn() io.Reader {
	return bytes.NewReader(m.sendData)
}

func (m *MockConnection) Send(message []byte) error {
	m.t.Helper()

	if !bytes.Equal(message, m.expectedReceive) {
		m.t.Errorf("result: %v, expected: %v", message, m.expectedReceive)
	}

	return nil
}

func (m *MockConnection) SetDeadline() error {
	return nil
}

func (m *MockConnection) Close() error {
	return nil
}

func TestPeerHandshake(t *testing.T) {
	p := torrent.Peer{}

	infoHash := [20]byte{0, 1, 2, 3, 4, 5, 6, 7, 8}
	peerID := [20]byte{0, 1, 2, 3, 4, 5, 6, 7, 8}

	message := []byte{19, 66, 105, 116, 84, 111, 114, 114, 101, 110, 116, 32, 112, 114, 111, 116, 111, 99, 111, 108, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0}

	mockConn := MockConnection{t: t, sendData: message, expectedReceive: message}
	p.SetConnection(&mockConn)

	if err := p.Handshake(infoHash, peerID); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPeerRecvBitfield(t *testing.T) {
	p := torrent.Peer{}

	message := []byte{0, 0, 0, 5, byte(torrent.MsgBitfield), 255, 255, 255, 255}
	expected := torrent.Bitfield([]byte{255, 255, 255, 255})

	mockConn := MockConnection{t: t, sendData: message, expectedReceive: nil}
	p.SetConnection(&mockConn)

	if err := p.RecvBitfield(); err != nil {
		t.Fatal(err)
	}

	result := p.Bitfield

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("result: %v, expected: %v", result, expected)
	}
}

func TestPeerSendUnchoke(t *testing.T) {
	p := torrent.Peer{}

	mockConn := MockConnection{t: t, expectedReceive: []byte{0, 0, 0, 1, 1}}
	p.SetConnection(&mockConn)

	p.SendUnchoke()
}

func TestPeerSendInterested(t *testing.T) {
	p := torrent.Peer{}

	mockConn := MockConnection{t: t, expectedReceive: []byte{0, 0, 0, 1, 2}}
	p.SetConnection(&mockConn)

	p.SendInterested()
}

func TestPeerSendRequest(t *testing.T) {
	p := torrent.Peer{}
	piece := torrent.Piece{Index: 1}

	p.Piece = &piece

	mockConn := MockConnection{t: t, expectedReceive: []byte{0, 0, 0, 13, 6, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 64, 0}}

	p.SetConnection(&mockConn)
	p.SendRequest()
}

func TestPeerReadMessage(t *testing.T) {
	p := torrent.Peer{}

	msgChoke := torrent.Message{ID: torrent.MsgChoke}
	msgUnchoke := torrent.Message{ID: torrent.MsgUnchoke}
	msgPiece := torrent.Message{ID: torrent.MsgPiece, Payload: []byte{0, 0, 0, 1, 0, 0, 0, 0, 1, 2, 3}}

	mockConn := MockConnection{t: t, sendData: msgChoke.Serialize()}
	p.SetConnection(&mockConn)

	if err := p.ReadMessage(); err != nil {
		t.Fatal(err)
	}

	result := p.IsChocked()
	expected := true

	if result != expected {
		t.Errorf("result: %t, expected: %t", result, expected)
	}

	mockConn = MockConnection{t: t, sendData: msgUnchoke.Serialize()}
	p.SetConnection(&mockConn)

	if err := p.ReadMessage(); err != nil {
		t.Fatal(err)
	}

	result = p.IsChocked()
	expected = false

	if result != expected {
		t.Errorf("result: %t, expected: %t", result, expected)
	}

	piece := torrent.Piece{Index: 1}
	piece.Data = make([]byte, 1024)

	data := [1024]byte{1, 2, 3}
	expectedPiece := torrent.Piece{Index: 1, Data: data[:]}

	p.Piece = &piece

	mockConn = MockConnection{t: t, sendData: msgPiece.Serialize()}
	p.SetConnection(&mockConn)

	if err := p.ReadMessage(); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(p.Piece, p.Piece) {
		t.Errorf("result: %v, expected: %v", p.Piece, expectedPiece)
	}
}
