package bencode

import (
	"errors"
	"io"
	"strconv"
)

func Unmarshal(data []byte) (any, error) {
	d := &decoder{data: data}
	return d.value()
}

type decoder struct {
	data []byte
	pos  int
}

func (d *decoder) value() (any, error) {
	if d.pos >= len(d.data) {
		return nil, io.ErrUnexpectedEOF
	}
	switch c := d.data[d.pos]; {
	case c == 'i':
		return d.integer()
	case c == 'l':
		return d.list()
	case c == 'd':
		return d.dict()
	case c >= '0' && c <= '9':
		return d.str()
	default:
		return nil, errors.New("bencode: invalid token")
	}
}

func (d *decoder) integer() (int64, error) {
	d.pos++
	start := d.pos
	for d.pos < len(d.data) && d.data[d.pos] != 'e' {
		d.pos++
	}
	if d.pos >= len(d.data) {
		return 0, io.ErrUnexpectedEOF
	}
	n, err := strconv.ParseInt(string(d.data[start:d.pos]), 10, 64)
	if err != nil {
		return 0, err
	}
	d.pos++
	return n, nil
}

func (d *decoder) str() ([]byte, error) {
	start := d.pos
	for d.pos < len(d.data) && d.data[d.pos] != ':' {
		d.pos++
	}
	if d.pos >= len(d.data) {
		return nil, io.ErrUnexpectedEOF
	}
	n, err := strconv.Atoi(string(d.data[start:d.pos]))
	if err != nil {
		return nil, err
	}
	d.pos++
	if n < 0 || d.pos+n > len(d.data) {
		return nil, io.ErrUnexpectedEOF
	}
	s := d.data[d.pos : d.pos+n]
	d.pos += n
	return s, nil
}

func (d *decoder) list() ([]any, error) {
	d.pos++
	var out []any
	for d.pos < len(d.data) && d.data[d.pos] != 'e' {
		v, err := d.value()
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	if d.pos >= len(d.data) {
		return nil, io.ErrUnexpectedEOF
	}
	d.pos++
	return out, nil
}

func (d *decoder) dict() (map[string]any, error) {
	d.pos++
	out := make(map[string]any)
	for d.pos < len(d.data) && d.data[d.pos] != 'e' {
		k, err := d.str()
		if err != nil {
			return nil, err
		}
		v, err := d.value()
		if err != nil {
			return nil, err
		}
		out[string(k)] = v
	}
	if d.pos >= len(d.data) {
		return nil, io.ErrUnexpectedEOF
	}
	d.pos++
	return out, nil
}
