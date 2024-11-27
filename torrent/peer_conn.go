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

func (pc *PeerConn) Send(message []byte) error {
	_, err := pc.Conn.Write(message)
	return err
}

func (pc *PeerConn) Receive(size int) ([]byte, error) {
	buff := make([]byte, size)
	_, err := pc.Conn.Read(buff)
	if err != nil {
		return []byte{}, err
	}

	return buff, nil
}
