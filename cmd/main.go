package main

import (
	"fmt"

	"github.com/juneira/torrent-study/torrent"
)

func main() {
	t, _ := torrent.FromFilename("tests/fixtures/debian-12.8.0-amd64-netinst.iso.torrent")

	pid := [20]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 10, 11}
	peers, err := t.GetPeers(pid, 80)
	if err != nil {
		panic(err)
	}

	var peer *torrent.Peer

	for _, p := range peers {
		peer = &p

		fmt.Printf("try connect to: %s:%d\n", string(peer.IP), peer.Port)

		peerConn, err := torrent.NewPeerConn(peer)
		peer.SetConnection(peerConn)

		if err != nil {
			continue
		}

		if err = peer.Handshake(t.InfoHash, pid); err != nil {
			continue
		}

		if err = peer.RecvBitfield(); err != nil {
			continue
		}

		break
	}

	pieces := t.GetPieces()

	if peer.Bitfield.HasPiece(0) {
		peer.Piece = pieces[0]
	}

	peer.DownloadPieces()

	for _, piece := range peer.Pieces() {
	}
}
