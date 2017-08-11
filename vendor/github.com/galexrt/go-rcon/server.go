package rcon

import (
	"bytes"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
)

// DialFn connect to server using the options
type DialFn func(network, address string) (net.Conn, error)

// Server represents a Source engine game server.
type Server struct {
	addr string

	dial DialFn

	rconPassword string

	timeout time.Duration

	rsock           *rconSocket
	rconInitialized bool

	mu sync.Mutex
}

// ConnectOptions describes the various connections options.
type ConnectOptions struct {
	// Default will use net.Dialer.Dial. You can override the same by
	// providing your own.
	Dial DialFn

	// RCON password.
	RCONPassword string

	Timeout string
}

// Connect to the source server.
func Connect(addr string, os ...*ConnectOptions) (_ *Server, err error) {
	s := &Server{
		addr: addr,
	}
	if len(os) > 0 {
		o := os[0]
		s.dial = o.Dial
		s.rconPassword = o.RCONPassword
		s.timeout, err = time.ParseDuration(o.Timeout)
		if err != nil {
			log.WithFields(logrus.Fields{
				"err": err,
			}).Fatal("rcon: could not parse timeout duration")
		}
	}
	if s.dial == nil {
		s.dial = (&net.Dialer{
			Timeout: 1 * time.Second,
		}).Dial
	}
	if s.rconPassword == "" {
		return s, nil
	}
	if err := s.initRCON(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) String() string {
	return s.addr
}

func (s *Server) initRCON() (err error) {
	if s.addr == "" {
		return errors.New("rcon: server needs a address")
	}
	log.WithFields(logrus.Fields{
		"addr": s.addr,
	}).Debug("rcon: connecting rcon")
	if s.rsock, err = newRCONSocket(s.dial, s.addr, s.timeout); err != nil {
		log.WithFields(logrus.Fields{
			"err": err,
		}).Error("rcon: could not open tcp socket")
		return err
	}
	defer func() {
		if err != nil {
			s.rsock.close()
		}
	}()
	if err := s.authenticate(); err != nil {
		log.WithFields(logrus.Fields{
			"err": err,
		}).Error("rcon: could not authenticate")
		return err
	}
	s.rconInitialized = true
	return nil
}

func (s *Server) authenticate() error {
	log.WithFields(logrus.Fields{
		"addr": s.addr,
	}).Debug("rcon: authenticating")
	req := newRCONRequest(rrtAuth, s.rconPassword)
	data, _ := req.marshalBinary()
	if err := s.rsock.send(data); err != nil {
		return err
	}
	// Receive the empty response value
	data, err := s.rsock.receive()
	if err != nil {
		return err
	}
	log.WithFields(logrus.Fields{
		"data": data,
	}).Debug("rcon: received empty response")
	var resp rconResponse
	if err = resp.unmarshalBinary(data); err != nil {
		return err
	}
	if resp.typ != rrtRespValue || resp.id != req.id {
		return ErrInvalidResponseID
	}
	if resp.id != req.id {
		return ErrInvalidResponseType
	}
	// Receive the actual auth response
	data, err = s.rsock.receive()
	if err != nil {
		return err
	}
	if err := resp.unmarshalBinary(data); err != nil {
		return err
	}
	if resp.typ != rrtAuthResp || resp.id != req.id {
		return ErrRCONAuthFailed
	}
	log.Debug("rcon: authenticated")
	return nil
}

// Close releases the resources associated with this server.
func (s *Server) Close() {
	if s.rconInitialized {
		s.rsock.close()
	}
}

// Send RCON command to the server.
func (s *Server) Send(cmd string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.rconInitialized {
		return "", ErrRCONNotInitialized
	}
	req := newRCONRequest(rrtExecCmd, cmd)
	data, _ := req.marshalBinary()
	if err := s.rsock.send(data); err != nil {
		log.WithFields(logrus.Fields{
			"err": err,
		}).Error("rcon: sending rcon request")
		return "", err
	}
	// Send the mirror packet.
	reqMirror := newRCONRequest(rrtRespValue, "")
	data, _ = reqMirror.marshalBinary()
	if err := s.rsock.send(data); err != nil {
		log.WithFields(logrus.Fields{
			"err": err,
		}).Error("rcon: sending rcon mirror request")
		return "", err
	}
	var (
		buf       bytes.Buffer
		sawMirror bool
	)
	// Start receiving data.
	for {
		data, err := s.rsock.receive()
		if err != nil {
			log.WithFields(logrus.Fields{
				"err": err,
			}).Error("rcon: receiving rcon response")
			return "", err
		}
		var resp rconResponse
		if err = resp.unmarshalBinary(data); err != nil {
			log.WithFields(logrus.Fields{
				"err": err,
			}).Error("rcon: decoding response")
			return "", err
		}
		if resp.typ != rrtRespValue {
			return "", ErrInvalidResponseType
		}
		if !sawMirror && resp.id == reqMirror.id {
			sawMirror = true
			continue
		}
		if sawMirror {
			if bytes.Compare(resp.body, trailer) == 0 {
				break
			}
			return "", ErrInvalidResponseTrailer
		}
		if req.id != resp.id {
			return "", ErrInvalidResponseID
		}
		_, err = buf.Write(resp.body)
		if err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

var (
	trailer = []byte{0x00, 0x01, 0x00, 0x00}

	ErrRCONAuthFailed         = errors.New("rcon: authentication failed")
	ErrRCONNotInitialized     = errors.New("rcon: rcon is not initialized")
	ErrInvalidResponseType    = errors.New("rcon: invalid response type from server")
	ErrInvalidResponseID      = errors.New("rcon: invalid response id from server")
	ErrInvalidResponseTrailer = errors.New("rcon: invalid response trailer from server")
)
