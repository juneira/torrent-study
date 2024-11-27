package torrent

import (
	"fmt"
	"net"
)

type Connection interface {
	Send(message []byte) error
	Receive(size int) ([]byte, error)
}

type Peer struct {
	IP   net.IP
	Port uint16
	conn Connection
}

func (p *Peer) SetConnection(conn Connection) {
	p.conn = conn
}

func (p *Peer) Address() string {
	return fmt.Sprintf("%s:%d", string(p.IP), p.Port)
}
