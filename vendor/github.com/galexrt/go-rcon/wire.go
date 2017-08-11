package rcon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"strconv"
)

type parseError string

func (e parseError) Error() string {
	return string(e)
}

var (
	errCouldNotReadData        = errors.New("rcon: could not read data")
	errNotEnoughDataInResponse = errors.New("rcon: not enough data in response")
	errBadData                 = errors.New("rcon: bad data in response")
)

func readByte(r io.Reader) (byte, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(r, buf)
	return buf[0], err
}

func readBytes(r io.Reader, n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(r, buf)
	return buf, err
}

func readShort(r io.Reader) (int16, error) {
	var v int16
	err := binary.Read(r, binary.LittleEndian, &v)
	return v, err
}

func readLong(r io.Reader) (int32, error) {
	var v int32
	err := binary.Read(r, binary.LittleEndian, &v)
	return v, err
}

func readULong(r io.Reader) (uint32, error) {
	var v uint32
	err := binary.Read(r, binary.LittleEndian, &v)
	return v, err
}

func readLongLong(r io.Reader) (int64, error) {
	var v int64
	err := binary.Read(r, binary.LittleEndian, &v)
	return v, err
}

func readString(r io.Reader) (string, error) {
	if buf, ok := r.(*bytes.Buffer); ok {
		// See if we are being passed a bytes.Buffer.
		// Fast path.
		bytes, err := buf.ReadBytes(0)
		return string(bytes), err
	}
	var buf bytes.Buffer
	for {
		b := make([]byte, 1)
		_, err := io.ReadFull(r, b)
		if err != nil {
			return "", err
		}
		buf.WriteByte(b[0])
		if b[0] == 0 {
			break
		}
	}
	return buf.String(), nil
}

func readFloat(r io.Reader) (float32, error) {
	v, err := readULong(r)
	return math.Float32frombits(v), err
}

func toInt(v interface{}) (int, error) {
	switch v := v.(type) {
	case byte:
	case int16:
	case int32:
	case int64:
		return int(v), nil
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, errBadData
}

func writeRequestPrefix(buf *bytes.Buffer) error {
	_, err := buf.Write(requestPrefix)
	return err
}

var requestPrefix = []byte{0xFF, 0xFF, 0xFF, 0xFF}

func writeString(buf *bytes.Buffer, v string) {
	buf.WriteString(v)
	buf.WriteByte(0)
}

func writeByte(buf *bytes.Buffer, v byte) {
	buf.WriteByte(v)
}

func writeLong(buf *bytes.Buffer, v int32) error {
	return binary.Write(buf, binary.LittleEndian, v)
}

func writeNull(buf *bytes.Buffer) error {
	return buf.WriteByte(0)
}
