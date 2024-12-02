package torrent

import (
	"encoding/binary"
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

func (m *Message) Serialize() []byte {
	length := uint32(len(m.Payload) + 1) // +1 for ID
	buf := make([]byte, 4+length)

	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload[:])

	return buf
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
