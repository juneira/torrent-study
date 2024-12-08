package torrent

import (
	"encoding/binary"
	"errors"
	"io"
)

type messageID uint8

const (
	MsgChoke         messageID = 0
	MsgUnchoke       messageID = 1
	MsgInterested    messageID = 2
	MsgNotInterested messageID = 3
	MsgHave          messageID = 4
	MsgBitfield      messageID = 5
	MsgRequest       messageID = 6
	MsgPiece         messageID = 7
	MsgCancel        messageID = 8
)

type Message struct {
	ID      messageID
	Payload []byte
}

func FormatRequest(index, begin, length int) *Message {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:], uint32(length))

	return &Message{ID: MsgRequest, Payload: payload}
}

func (m *Message) Serialize() []byte {
	length := uint32(len(m.Payload) + 1) // +1 for ID
	buf := make([]byte, 4+length)

	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload[:])

	return buf
}

func (m *Message) ParsePieceIndex() (int, error) {
	return int(binary.BigEndian.Uint32(m.Payload[0:4])), nil
}

func (m *Message) ParsePiece(index int, buf []byte) (int, error) {
	if m.ID != MsgPiece {
		return 0, errors.New("invalid messaage id")
	}

	if len(m.Payload) < 8 {
		return 0, errors.New("payload too short")
	}

	messageIndex := int(binary.BigEndian.Uint32(m.Payload[0:4]))
	if index != messageIndex {
		return 0, errors.New("invalid index")
	}

	begin := int(binary.BigEndian.Uint32(m.Payload[4:8]))
	if begin >= len(buf) {
		return 0, errors.New("begin offset too high")
	}

	data := m.Payload[8:]
	if begin+len(data) > len(buf) {
		return 0, errors.New("buf too small to receive data")
	}

	copy(buf[begin:], data)
	return len(data), nil
}

func readerToMessage(r io.Reader) (*Message, error) {
	m := Message{}

	lengthBuff := [4]byte{}
	_, err := io.ReadFull(r, lengthBuff[:])
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lengthBuff[:])
	messageBuff := make([]byte, length)

	_, err = io.ReadFull(r, messageBuff)
	if err != nil {
		return nil, err
	}

	m.ID = messageID(messageBuff[0])
	m.Payload = messageBuff[1:]

	return &m, nil
}
