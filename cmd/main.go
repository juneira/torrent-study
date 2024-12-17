package main

import (
	"github.com/juneira/torrent-study/torrent"
)

func main() {
	t, _ := torrent.FromFilename("tests/fixtures/debian-12.8.0-amd64-netinst.iso.torrent")

	pid := [20]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 10, 11}
	port := uint16(80)

	peers, err := t.GetPeers(pid, port)
	if err != nil {
		panic(err)
	}

	p2p := torrent.NewP2P(t.InfoHash, pid, peers, t.GetPieces())

	err = p2p.Download()
	if err != nil {
		panic(err)
	}
}
