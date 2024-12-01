package torrent_test

import (
	"bytes"
	"io"
	"net"
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
	p := torrent.Peer{IP: net.IP("127.0.0.1"), Port: uint16(5555)}

	infoHash := [20]byte{0, 1, 2, 3, 4, 5, 6, 7, 8}
	peerID := [20]byte{0, 1, 2, 3, 4, 5, 6, 7, 8}

	message := []byte{19, 66, 105, 116, 84, 111, 114, 114, 101, 110, 116, 32, 112, 114, 111, 116, 111, 99, 111, 108, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0}

	mock := MockConnection{t: t, sendData: message, expectedReceive: message}
	p.SetConnection(&mock)

	if err := p.Handshake(infoHash, peerID); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
