package torrent

import (
	"bytes"
	"crypto/sha1"

	"github.com/jackpal/bencode-go"
)

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

type bencodeInfo struct {
	Pieces       string `bencode:"pieces"`
	PiecesLength int    `bencode:"piece length"`
	Length       int    `bencode:"length"`
	Name         string `bencode:"name"`
}

func (bto bencodeTorrent) toTorrentFile() (*TorrentFile, error) {
	var err error

	tf := TorrentFile{}
	tf.InfoHash, err = bto.Info.hash()
	if err != nil {
		return nil, err
	}

	tf.Announce = bto.Announce
	tf.Length = bto.Info.Length
	tf.Name = bto.Info.Name
	tf.PieceHashes = bto.Info.getPiecesHashs()
	tf.PiecesLength = bto.Info.PiecesLength

	return &tf, nil
}

func (bi bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, bi)
	if err != nil {
		return [20]byte{}, err
	}

	return sha1.Sum(buf.Bytes()), nil
}

func (bi bencodeInfo) getPiecesHashs() [][20]byte {
	var pieceHashs [][20]byte

	for i := 0; i < len(bi.Pieces); i += 20 {
		end := i + 20
		if i == bi.PiecesLength-1 {
			end = len(bi.Pieces)
		}

		var piece [20]byte
		copy(piece[:], bi.Pieces[i:end])

		pieceHashs = append(pieceHashs, piece)
	}

	return pieceHashs
}
