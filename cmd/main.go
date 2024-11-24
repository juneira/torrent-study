package main

import (
	"fmt"

	"github.com/juneira/torrent-study/torrent"
)

func main() {
	t, _ := torrent.FromFilename("fixtures/debian-12.8.0-amd64-netinst.iso.torrent")

	fmt.Println(t.PieceHashes)
}
