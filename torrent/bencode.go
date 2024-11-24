package torrent

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
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
	tf := TorrentFile{}
	tf.Announce = bto.Announce
	tf.Length = bto.Info.Length
	tf.Name = bto.Info.Name
	tf.InfoHash = bto.Info.hash()
	tf.PieceHashes = bto.Info.getPiecesHashs()
	tf.PiecesLength = bto.Info.PiecesLength

	return &tf, nil
}

func (bi bencodeInfo) hash() [20]byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, bi)

	return sha1.Sum(buf.Bytes())
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
