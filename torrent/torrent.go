package torrent

import (
	"bytes"
	"os"

	bencode "github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	Announce     string
	InfoHash     [20]byte
	PieceHashes  [][20]byte
	PiecesLength int
	Length       int
	Name         string
}

func FromFilename(filename string) (*TorrentFile, error) {
	bencodeBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(bencodeBytes)
	bto := bencodeTorrent{}

	if err := bencode.Unmarshal(r, &bto); err != nil {
		return nil, err
	}

	return bto.toTorrentFile()
}
