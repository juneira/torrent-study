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

	peer := peers[0]
	peerConn, err := torrent.NewPeerConn(&peer)
	if err != nil {
		panic(err)
	}

	peer.SetConnection(peerConn)

	err = peer.Handshake(t.InfoHash, pid)
	if err != nil {
		panic(err)
	}

	err = peer.RecvBitfield()
	if err != nil {
		panic(err)
	}

	fmt.Println(peer.Bitfield)
}
