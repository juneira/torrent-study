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
	SetDeadline() error
	Close() error
}

type Peer struct {
	IP       net.IP
	Port     uint16
	Bitfield Bitfield
	Choked   bool
	Piece    *Piece
	conn     Connection
}

func (p *Peer) SetConnection(conn Connection) {
	p.conn = conn
}

func (p *Peer) DownloadPiece() error {
	if err := p.SendUnchoke(); err != nil {
		panic(err)
	}

	if err := p.SendInterested(); err != nil {
		panic(err)
	}

	if p.Piece == nil {
		return errors.New("peer has none piece")
	}

	for p.Choked {
		if err := p.ReadMessage(); err != nil {
			return err
		}
	}

	if err := p.conn.SetDeadline(); err != nil {
		return err
	}

	return p.downloadPiece(p.Piece)
}

func (p *Peer) downloadPiece(piece *Piece) error {
	for piece.Status == PiecePending {
		if err := p.SendRequest(); err != nil {
			piece.Status = PieceError

			return err
		}

		piece.Waiting = true

		for piece.Waiting {
			if err := p.ReadMessage(); err != nil {
				piece.Status = PieceError

				return err
			}
		}

		if err := piece.CheckFinished(); err != nil {
			return err
		}

		if piece.Status == PieceFinished {
			return p.SendHave()
		}
	}

	return nil
}

func (p *Peer) IsChocked() bool {
	return p.Choked
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

func (p *Peer) SendRequest() error {
	if p.Piece == nil {
		return errors.New("not exists piece to this peer")
	}

	m := FormatRequest(p.Piece.Index, p.Piece.Begin, LengthMax)

	return p.conn.Send(m.Serialize())
}

func (p *Peer) SendHave() error {
	if p.Piece == nil {
		return errors.New("not exists piece to this peer")
	}

	m := FormatHave(p.Piece.Index)

	return p.conn.Send(m.Serialize())
}

func (p *Peer) ReadMessage() error {
	m, err := readerToMessage(p.conn.GetConn())
	if err != nil {
		return err
	}

	switch m.ID {
	case MsgChoke:
		p.Choked = true
	case MsgUnchoke:
		p.Choked = false
	case MsgPiece:
		index, err := m.ParsePieceIndex()
		if err != nil {
			return err
		}

		if index != p.Piece.Index {
			return errors.New("invalid index")
		}

		downloaded, err := m.ParsePiece(index, p.Piece.Data)
		if err != nil {
			return err
		}

		p.Piece.Begin += downloaded
		p.Piece.Waiting = false
	}

	return nil
}
