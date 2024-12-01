package torrent

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
)

type Connection interface {
	GetConn() io.Reader
	Send(message []byte) error
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

func (p *Peer) Handshake(infoHash [20]byte, peerID [20]byte) error {
	messageHS := Handshake{Pstr: "BitTorrent protocol", InfoHash: infoHash, PeerID: peerID}
	message := messageHS.Serialize()

	if err := p.conn.Send(message); err != nil {
		return err
	}

	receiveHS, err := readerToHandshake(p.conn.GetConn())
	if err != nil {
		return err
	}

	if !bytes.Equal(receiveHS.InfoHash[:], infoHash[:]) {
		errorMessage := fmt.Sprintf("invalid infoHash: expected: %v | returned: %v", infoHash, receiveHS.InfoHash)

		return errors.New(errorMessage)
	}

	return nil
}
