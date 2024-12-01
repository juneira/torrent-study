package torrent_test

import (
	"net"
	"reflect"
	"testing"

	"github.com/juneira/torrent-study/torrent"
)

const ADDR_TEST = "127.0.0.1:5555"

var done chan struct{}

func createServer(t *testing.T, response []byte, expectedReceive []byte) (serverConn net.Conn) {
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

		if len(response) > 0 {
			serverConn.Write(response)
		}

		if len(expectedReceive) > 0 {
			t.Helper()

			result := make([]byte, len(expectedReceive))
			_, err = serverConn.Read(result)

			if err != nil {
				done <- struct{}{}
				t.Fatal(err)
			}

			if !reflect.DeepEqual(result, expectedReceive) {
				done <- struct{}{}
				t.Errorf("result: %v, expected: %v", result, expectedReceive)
			}

			done <- struct{}{}
		}
	}()

	return serverConn
}

func TestPeerConnSend(t *testing.T) {
	message := []byte{0, 1, 2, 3}

	createServer(t, []byte{}, message)

	p := torrent.Peer{IP: net.IP("127.0.0.1"), Port: uint16(5555)}

	done = make(chan struct{})

	pc, err := torrent.NewPeerConn(&p)
	if err != nil {
		t.Fatal(err)
	}

	if err := pc.Send(message); err != nil {
		t.Fatal(err)
	}
	<-done
}

func TestPeerConnReceive(t *testing.T) {
	message := []byte{0, 1, 2, 3}

	createServer(t, message, []byte{})

	p := torrent.Peer{IP: net.IP("127.0.0.1"), Port: uint16(5555)}

	done = make(chan struct{})

	_, err := torrent.NewPeerConn(&p)
	if err != nil {
		t.Fatal(err)
	}
}
