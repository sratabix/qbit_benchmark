package peer

import (
	"crypto/rand"
	"errors"
	"io"
)

func NewPeerID() [20]byte {
	var id [20]byte
	copy(id[:], "-QB5000-")
	_, _ = rand.Read(id[8:])
	return id
}

func writeHandshake(w io.Writer, infoHash, peerID [20]byte) error {
	buf := make([]byte, 0, 68)
	buf = append(buf, byte(len(protocolName)))
	buf = append(buf, protocolName...)
	buf = append(buf, make([]byte, 8)...)
	buf = append(buf, infoHash[:]...)
	buf = append(buf, peerID[:]...)
	_, err := w.Write(buf)
	return err
}

func readHandshake(r io.Reader) (infoHash [20]byte, err error) {
	var pstrlen [1]byte
	if _, err = io.ReadFull(r, pstrlen[:]); err != nil {
		return
	}
	buf := make([]byte, int(pstrlen[0])+48)
	if _, err = io.ReadFull(r, buf); err != nil {
		return
	}
	if string(buf[:pstrlen[0]]) != protocolName {
		err = errors.New("peer: bad protocol string")
		return
	}
	off := int(pstrlen[0]) + 8
	copy(infoHash[:], buf[off:off+20])
	return
}
