package rcon

import (
	"bytes"
	"io"
	"net"
	"time"

	"github.com/Sirupsen/logrus"
)

type rconSocket struct {
	conn    net.Conn
	timeout time.Duration
}

func newRCONSocket(dial DialFn, addr string, timeout time.Duration) (*rconSocket, error) {
	conn, err := dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &rconSocket{conn, timeout}, nil
}

func (s *rconSocket) close() {
	s.conn.Close()
}

func (s *rconSocket) send(p []byte) error {
	if err := s.conn.SetWriteDeadline(time.Now().Add(s.timeout)); err != nil {
		return err
	}
	_, err := s.conn.Write(p)
	if err != nil {
		return err
	}
	return nil
}

func (s *rconSocket) receive() (_ []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	buf := new(bytes.Buffer)
	tr := io.TeeReader(s.conn, buf)
	long, err := readLong(tr)
	if err != nil {
		return nil, err
	}
	total := int(long)
	log.WithFields(logrus.Fields{
		"total": total + 4,
	}).Debug("rcon: reading packet")
	for total > 0 {
		log.WithFields(logrus.Fields{
			"bytes": total,
		}).Debug("rcon: reading")
		b := make([]byte, total)
		if err = s.conn.SetReadDeadline(time.Now().Add(s.timeout)); err != nil {
			return nil, err
		}
		n, err := s.conn.Read(b)
		if n > 0 {
			log.WithFields(logrus.Fields{
				"bytes": n,
			}).Debug("rcon: read")
			if _, err = buf.Write(b); err != nil {
				return nil, err
			}
			total -= n
		}
		if err != nil {
			if err == io.EOF {
				log.WithFields(logrus.Fields{
					"size": buf.Len(),
				}).Debug("rcon: read EOF")
				break
			}
			log.WithFields(logrus.Fields{
				"err": err,
			}).Error("rcon: could not receive data")
			return nil, err
		}
		log.WithFields(logrus.Fields{
			"bytes": total,
		}).Debug("rcon: remaining")
	}
	log.WithFields(logrus.Fields{
		"size": buf.Len(),
	}).Debug("rcon: read packet")
	return buf.Bytes(), nil
}
