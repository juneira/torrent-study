package torrent

import (
	"io"
	"net"
	"time"
)

type PeerConn struct {
	conn net.Conn
}

const TIMEOUT_IN_SECONDS = 3

func NewPeerConn(p *Peer) (*PeerConn, error) {
	pc := PeerConn{}
	conn, err := net.DialTimeout("tcp", p.Address(), TIMEOUT_IN_SECONDS*time.Second)
	if err != nil {
		return nil, err
	}

	pc.conn = conn
	return &pc, nil
}

func (pc *PeerConn) GetConn() io.Reader {
	return pc.conn
}

func (pc *PeerConn) Send(message []byte) error {
	_, err := pc.conn.Write(message)
	return err
}

func (pc *PeerConn) Close() error {
	return pc.conn.Close()
}

func (pc *PeerConn) SetDeadline() error {
	return pc.conn.SetDeadline(time.Now().Add(30 * time.Second))
}
