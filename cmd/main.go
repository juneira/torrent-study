package main

import (
	"fmt"

	"github.com/juneira/torrent-study/torrent"
)

func main() {
	t, _ := torrent.FromFilename("tests/fixtures/debian-12.8.0-amd64-netinst.iso.torrent")

	fmt.Println(t.Announce)

	pid := [20]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 10, 11}
	t.GetPeers(pid, 80)
}
