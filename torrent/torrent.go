package torrent

import (
	"bytes"
	"fmt"
	"io"
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

func (t *TorrentFile) GetPeers(peerID [20]byte, port uint16) {
	url, err := t.buildTrackerURL(peerID, port)
	if err != nil {
		panic(err)
	}

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	fmt.Println(string(body))
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
