package torrent_test

import (
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/juneira/torrent-study/torrent"
)

func TestFromFilename(t *testing.T) {
	tt, err := torrent.FromFilename("fixtures/debian-12.8.0-amd64-netinst.iso.torrent")
	if err != nil {
		t.Fatal(err)
	}

	expected := torrent.TorrentFile{
		Name:         "debian-12.8.0-amd64-netinst.iso",
		Announce:     "http://bttracker.debian.org:6969/announce",
		PiecesLength: 262144,
		Length:       661651456,
	}

	if tt.Announce != expected.Announce {
		t.Errorf(`expected: %s, result: %s`, expected.Announce, tt.Announce)
	}

	if tt.PiecesLength != expected.PiecesLength {
		t.Errorf(`expected: %d, result: %d`, expected.PiecesLength, tt.PiecesLength)
	}

	if tt.Length != expected.Length {
		t.Errorf(`expected: %d, result: %d`, expected.PiecesLength, tt.PiecesLength)
	}

	if tt.Name != expected.Name {
		t.Errorf(`expected: %s, result: %s`, expected.Name, tt.Name)
	}
}

func TestTorrentFileGetPeers(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := []byte(
			"d8:intervali900e5:peersld2:ip13:99.104.171.974:porti52350eed2:ip12:92.190.54.114:porti51413eeee",
		)

		w.Write(response)
	}))

	tt, err := torrent.FromFilename("fixtures/debian-12.8.0-amd64-netinst.iso.torrent")
	if err != nil {
		t.Fatal(err)
	}

	tt.Announce = mockServer.URL

	pid := [20]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 10, 11}
	result, err := tt.GetPeers(pid, 80)
	if err != nil {
		t.Fatal(err)
	}

	expectedPeers := []torrent.Peer{
		{
			IP:   net.IP("99.104.171.97"),
			Port: uint16(52350),
		},
		{
			IP:   net.IP("92.190.54.11"),
			Port: uint16(51413),
		},
	}

	if !reflect.DeepEqual(expectedPeers, result) {
		t.Errorf(`expected: %v, result: %v`, expectedPeers, result)
	}
}

func TestTorrentFileGetPieces(t *testing.T) {
	tf := torrent.TorrentFile{PiecesLength: 5, Length: 13}

	// 3 pieces
	tf.PieceHashes = append(tf.PieceHashes, [20]byte{})
	tf.PieceHashes = append(tf.PieceHashes, [20]byte{})
	tf.PieceHashes = append(tf.PieceHashes, [20]byte{})

	dataA := make([]byte, 5)
	dataB := make([]byte, 3)

	expected := []*torrent.Piece{
		{Index: 0, Hash: [20]byte{}, Data: dataA},
		{Index: 1, Hash: [20]byte{}, Data: dataA},
		{Index: 2, Hash: [20]byte{}, Data: dataB},
	}

	result := tf.GetPieces()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("result: %v, expected: %v", result, expected)
	}
}
