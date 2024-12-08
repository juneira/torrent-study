package torrent

const LengthMax = 1024

type Piece struct {
	Index int
	Hash  [20]byte
	Data  []byte
	Begin int
}
