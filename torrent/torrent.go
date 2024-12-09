package torrent

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"

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

func (t *TorrentFile) GetPeers(peerID [20]byte, port uint16) ([]Peer, error) {
	url, err := t.buildTrackerURL(peerID, port)
	if err != nil {
		return []Peer{}, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return []Peer{}, err
	}

	defer resp.Body.Close()

	return t.decodePeers(resp.Body)
}

func (t *TorrentFile) GetPieces() (pieces []*Piece) {
	for index, piecehash := range t.PieceHashes {
		piece := Piece{Index: index, Hash: piecehash}
		piece.Data = make([]byte, t.calculatePieceSize(index))
		pieces = append(pieces, &piece)
	}

	return pieces
}

func (t *TorrentFile) calculateInitAndEndByPieceIndex(index int) (begin, end int) {
	begin = index * t.PiecesLength
	end = begin + t.PiecesLength
	if end > t.Length {
		end = t.Length
	}

	return begin, end
}

func (t *TorrentFile) calculatePieceSize(index int) int {
	begin, end := t.calculateInitAndEndByPieceIndex(index)
	return end - begin
}

func (t *TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}

	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"conmpact":   []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}

	base.RawQuery = params.Encode()
	return base.String(), nil
}

func (t *TorrentFile) decodePeers(bencodeIO io.Reader) ([]Peer, error) {
	var bps bencodePeers

	if err := bencode.Unmarshal(bencodeIO, &bps); err != nil {
		return nil, err
	}

	var peers []Peer

	for _, bp := range bps.Peers {
		var peer Peer

		peer.Port = uint16(bp.Port)
		peer.IP = net.IP(bp.IP)

		peers = append(peers, peer)
	}

	return peers, nil
}
