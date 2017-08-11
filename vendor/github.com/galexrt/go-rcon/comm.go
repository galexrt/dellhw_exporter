package rcon

import (
	"bytes"
	"fmt"
	"math/rand"
)

type rconRequestType int32

const (
	rrtAuth      rconRequestType = 3
	rrtExecCmd   rconRequestType = 2
	rrtAuthResp  rconRequestType = 2
	rrtRespValue rconRequestType = 0
)

type rconRequest struct {
	size int32
	id   int32
	typ  rconRequestType
	body string
}

func (r *rconRequest) String() string {
	return fmt.Sprintf("%v %v %v %v", r.size, r.id, r.typ, r.body)
}

func newRCONRequest(typ rconRequestType, body string) *rconRequest {
	return &rconRequest{
		size: int32(len(body) + 10),
		id:   rand.Int31(),
		typ:  typ,
		body: body,
	}
}

func (r *rconRequest) marshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	writeLong(buf, r.size)
	writeLong(buf, r.id)
	writeLong(buf, int32(r.typ))
	buf.WriteString(r.body)
	writeNull(buf)
	writeNull(buf)
	return buf.Bytes(), nil
}

type rconResponse struct {
	size int32
	id   int32
	typ  rconRequestType
	body []byte
}

func (r *rconResponse) String() string {
	return fmt.Sprintf("%v %v %v %v", r.size, r.id, r.typ, string(r.body))
}

func (r *rconResponse) unmarshalBinary(data []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	buf := bytes.NewBuffer(data)
	if r.size, err = readLong(buf); err != nil {
		return err
	}
	if r.id, err = readLong(buf); err != nil {
		return err
	}
	typ, err := readLong(buf)
	if err != nil {
		return err
	}
	r.typ = rconRequestType(typ)
	if r.body, err = readBytes(buf, int(len(data)-14)); err != nil {
		return err
	}
	return err
}
