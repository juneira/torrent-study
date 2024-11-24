package torrent_test

import (
	"testing"

	"github.com/juneira/torrent-study/torrent"
)

func TestFromFilename(t *testing.T) {
	tt, err := torrent.FromFilename("fixtures/debian-12.8.0-amd64-netinst.iso.torrent")
	if err != nil {
		t.Error(err)
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
