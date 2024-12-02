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
