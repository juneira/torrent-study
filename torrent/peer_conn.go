package torrent

import (
	"net"
	"time"
)

type PeerConn struct {
	Conn net.Conn
}

const TIMEOUT_IN_SECONDS = 3

func NewPeerConn(p *Peer) (*PeerConn, error) {
	pc := PeerConn{}
	conn, err := net.DialTimeout("tcp", p.Address(), TIMEOUT_IN_SECONDS*time.Second)
	if err != nil {
		return nil, err
	}

	pc.Conn = conn
	return &pc, nil
}
