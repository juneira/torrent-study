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
	IP       net.IP
	Port     uint16
	Bitfield Bitfield
	choked   bool
	pieces   []*Piece
	conn     Connection
}

func (p *Peer) SetConnection(conn Connection) {
	p.conn = conn
}

func (p *Peer) AddPiece(piece *Piece) {
	p.pieces = append(p.pieces, piece)
}

func (p *Peer) Pieces() []*Piece {
	return p.pieces
}

func (p *Peer) IsChocked() bool {
	return p.choked
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

func (p *Peer) RecvBitfield() error {
	m, err := readerToMessage(p.conn.GetConn())
	if err != nil {
		return err
	}

	if m.ID != MsgBitfield {
		return errors.New("Invalid Message ID")
	}

	p.Bitfield = m.Payload

	return nil
}

func (p *Peer) SendUnchoke() error {
	m := Message{ID: MsgUnchoke}

	return p.conn.Send(m.Serialize())
}

func (p *Peer) SendInterested() error {
	m := Message{ID: MsgInterested}

	return p.conn.Send(m.Serialize())
}

func (p *Peer) SendRequest(index int) error {
	piece := p.getPieceByIndex(index)

	if piece == nil {
		return errors.New("not exists piece to this peer")
	}

	m := FormatRequest(index, piece.Begin, LengthMax)

	return p.conn.Send(m.Serialize())
}

func (p *Peer) ReadMessage() error {
	m, err := readerToMessage(p.conn.GetConn())
	if err != nil {
		return err
	}

	switch m.ID {
	case MsgChoke:
		p.choked = true
	case MsgUnchoke:
		p.choked = false
	case MsgPiece:
		index, err := m.ParsePieceIndex()
		if err != nil {
			return err
		}

		piece := p.getPieceByIndex(index)
		_, err = m.ParsePiece(index, piece.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Peer) getPieceByIndex(index int) *Piece {
	for _, piece := range p.pieces {
		if piece.Index == index {
			return piece
		}
	}

	return nil
}
