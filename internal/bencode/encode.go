package bencode

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
)

func Marshal(v any) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encode(buf *bytes.Buffer, v any) error {
	switch x := v.(type) {
	case string:
		writeBytes(buf, []byte(x))
	case []byte:
		writeBytes(buf, x)
	case int:
		writeInt(buf, int64(x))
	case int64:
		writeInt(buf, x)
	case []any:
		buf.WriteByte('l')
		for _, item := range x {
			if err := encode(buf, item); err != nil {
				return err
			}
		}
		buf.WriteByte('e')
	case map[string]any:
		buf.WriteByte('d')
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			writeBytes(buf, []byte(k))
			if err := encode(buf, x[k]); err != nil {
				return err
			}
		}
		buf.WriteByte('e')
	default:
		return fmt.Errorf("bencode: unsupported type %T", v)
	}
	return nil
}

func writeBytes(buf *bytes.Buffer, b []byte) {
	buf.WriteString(strconv.Itoa(len(b)))
	buf.WriteByte(':')
	buf.Write(b)
}

func writeInt(buf *bytes.Buffer, n int64) {
	buf.WriteByte('i')
	buf.WriteString(strconv.FormatInt(n, 10))
	buf.WriteByte('e')
}
