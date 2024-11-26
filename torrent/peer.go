package torrent

import (
	"fmt"
	"net"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

func (p *Peer) Address() string {
	return fmt.Sprintf("%s:%d", string(p.IP), p.Port)
}
