package socks5

import (
	"encoding/binary"
	"io"
	"log/slog"
	"net"
	"strconv"
)

type (
	Server struct {
		addr string

		log *slog.Logger
	}

	handleFunc func(conn net.Conn, target string)
)

func New(addr string, log *slog.Logger) *Server {
	return &Server{
		addr: addr,
		log:  log,
	}
}

func (s *Server) Start(handle handleFunc) {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.log.Error(err.Error())
		return
	}

	s.log.Info("SOCKS5 listening on", s.addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			s.log.Error(err.Error())
			continue
		}

		go s.handleConn(conn, handle)
	}
}

func (s *Server) handleConn(conn net.Conn, handle handleFunc) {
	defer conn.Close()

	buf := make([]byte, 256)

	_, err := io.ReadFull(conn, buf[:2])
	if err != nil {
		s.log.Error(err.Error())
		return
	}

	ver := buf[0]
	nmethods := buf[1]

	if ver != 5 {
		s.log.Warn("not SOCKS5")
		return
	}

	_, err = io.ReadFull(conn, buf[:nmethods])
	if err != nil {
		s.log.Error(err.Error())
		return
	}

	_, err = conn.Write([]byte{0x05, 0x00})
	if err != nil {
		s.log.Error(err.Error())
		return
	}

	_, err = io.ReadFull(conn, buf[:4])
	if err != nil {
		s.log.Error(err.Error())
		return
	}

	cmd := buf[1]
	atyp := buf[3]

	if cmd != 0x01 {
		s.log.Warn("not connect")
		return
	}

	var target string
	switch atyp {
	case 0x01: // IPv4
		_, err = io.ReadFull(conn, buf[:4])
		if err != nil {
			s.log.Error(err.Error())
			return
		}

		ip := net.IP(buf[:4]).String()

		_, err = io.ReadFull(conn, buf[:2])
		if err != nil {
			s.log.Error(err.Error())
			return
		}

		port := int(binary.BigEndian.Uint16(buf[:2]))

		target = ip + ":" + strconv.Itoa(port)
	case 0x03: // domain
		_, err = io.ReadFull(conn, buf[:1])
		if err != nil {
			s.log.Error(err.Error())
			return
		}
		length := int(buf[0])

		_, err = io.ReadFull(conn, buf[:length])
		if err != nil {
			s.log.Error(err.Error())
			return
		}
		domain := string(buf[:length])

		_, err = io.ReadFull(conn, buf[:2])
		if err != nil {
			s.log.Error(err.Error())
			return
		}
		port := int(binary.BigEndian.Uint16(buf[:2]))

		target = domain + ":" + strconv.Itoa(port)
	case 0x04:
		_, err = io.ReadFull(conn, buf[:16])
		if err != nil {
			s.log.Error(err.Error())
			return
		}
		ip := net.IP(buf[:16]).String()

		_, err = io.ReadFull(conn, buf[:2])
		if err != nil {
			s.log.Error(err.Error())
			return
		}
		port := int(binary.BigEndian.Uint16(buf[:2]))

		target = ip + ":" + strconv.Itoa(port)
	}

	conn.Write([]byte{
		0x05, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00,
	})

	handle(conn, target)
}
