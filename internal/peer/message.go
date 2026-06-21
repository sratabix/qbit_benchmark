package peer

import (
	"encoding/binary"
	"io"
)

const (
	msgUnchoke    = 1
	msgInterested = 2
	msgBitfield   = 5
	msgRequest    = 6
	msgPiece      = 7
)

const protocolName = "BitTorrent protocol"

type message struct {
	id      byte
	payload []byte
}

func readMessage(r io.Reader) (*message, error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lenBuf[:])
	if length == 0 {
		return nil, nil
	}
	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return &message{id: buf[0], payload: buf[1:]}, nil
}

func writeMessage(w io.Writer, id byte, payload []byte) error {
	hdr := make([]byte, 5, 5+len(payload))
	binary.BigEndian.PutUint32(hdr[:4], uint32(1+len(payload)))
	hdr[4] = id
	_, err := w.Write(append(hdr, payload...))
	return err
}
