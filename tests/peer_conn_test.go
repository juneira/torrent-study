package torrent_test

import (
	"net"
	"testing"

	"github.com/juneira/torrent-study/torrent"
)

const ADDR_TEST = "127.0.0.1:5555"

func createServer(t *testing.T) (serverConn net.Conn) {
	ln, err := net.Listen("tcp", ADDR_TEST)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		defer ln.Close()
		serverConn, err = ln.Accept()
		if err != nil {
			panic(err)
		}
	}()

	return serverConn
}

func TestNewPeerConn(t *testing.T) {
	createServer(t)

	p := torrent.Peer{IP: net.IP("127.0.0.1"), Port: uint16(5555)}

	pc, err := torrent.NewPeerConn(&p)
	if err != nil {
		t.Fatal(err)
	}

	result := pc.Conn.RemoteAddr().String()
	expected := ADDR_TEST

	if result != expected {
		t.Errorf("result: %s, expected: %s", result, expected)
	}
}
