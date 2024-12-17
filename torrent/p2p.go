package torrent

import (
	"errors"
	"fmt"
)

const MAX_PIECE_PROCESSING = 200

type P2P struct {
	infoHash         [20]byte
	pid              [20]byte
	peers            []Peer
	pieces           []*Piece
	processingPieces chan *Piece
	finishedPieces   []*Piece
	start            bool
}

func NewP2P(infoHash [20]byte, pid [20]byte, peers []Peer, pieces []*Piece) *P2P {
	return &P2P{infoHash: infoHash, pid: pid, peers: peers, pieces: pieces}
}

func (p2p *P2P) Download() error {
	p2p.processingPieces = make(chan *Piece, MAX_PIECE_PROCESSING)

	go p2p.enqueuePieces()

	go p2p.downloadPieces()

	finished := -1
	for len(p2p.finishedPieces) != len(p2p.pieces) {
		// wait download to finish
		if finished != len(p2p.finishedPieces) {
			finished = len(p2p.finishedPieces)
			fmt.Printf("%d / %d\n", finished, len(p2p.pieces))
		}
	}

	return nil
}

func (p2p *P2P) enqueuePieces() {
	for _, piece := range p2p.pieces {
		p2p.processingPieces <- piece
		p2p.start = true
	}
}

func (p2p *P2P) downloadPieces() {
	for !p2p.start {
		// wait enqueue a piece
	}

	for len(p2p.processingPieces) > 0 {
		piece := <-p2p.processingPieces

		go func() {
			peer, err := p2p.connectToPeer(piece.Index)
			if err != nil {
				fmt.Printf("PIECE %d: %v\n", piece.Index, err)
				p2p.processingPieces <- piece
				return
			}

			defer peer.conn.Close()

			if peer == nil {
				panic("empty peers")
			}

			peer.Piece = piece

			fmt.Printf("downloading: PIECE %d - %s:%d\n", piece.Index, string(peer.IP), peer.Port)

			if err := peer.DownloadPiece(); err != nil {
				peer.Piece.Begin = 0
				peer.Piece.Waiting = false
				peer.Piece.Status = PiecePending

				p2p.processingPieces <- peer.Piece

				fmt.Printf("error: PIECE %d - %s:%d\n", piece.Index, string(peer.IP), peer.Port)
			}

			if peer.Piece.Status == PieceFinished {
				fmt.Printf("done: PIECE %d - %s:%d\n", piece.Index, string(peer.IP), peer.Port)

				p2p.finishedPieces = append(p2p.finishedPieces, peer.Piece)
			} else {
				fmt.Printf("error: PIECE %d - %s:%d\n", piece.Index, string(peer.IP), peer.Port)
			}
		}()
	}
}

func (p2p *P2P) connectToPeer(pieceIndex int) (*Peer, error) {
	for _, peer := range p2p.peers {
		if peer.Piece != nil {
			continue
		}

		peerConn, err := NewPeerConn(&peer)
		peer.SetConnection(peerConn)

		if err != nil {
			continue
		}

		if err = peer.Handshake(p2p.infoHash, p2p.pid); err != nil {
			continue
		}

		if err = peer.RecvBitfield(); err != nil {
			continue
		}

		if !peer.Bitfield.HasPiece(uint(pieceIndex)) {
			continue
		}

		return &peer, nil
	}

	return nil, errors.New("piece not found")
}
