package torrent_test

import (
	"net"
	"testing"

	"github.com/juneira/torrent-study/torrent"
)

func TestAddress(t *testing.T) {
	p := torrent.Peer{IP: net.IP("127.0.0.1"), Port: uint16(5555)}

	result := p.Address()
	expected := "127.0.0.1:5555"

	if result != expected {
		t.Errorf("result: %s, expected: %s", result, expected)
	}
}
