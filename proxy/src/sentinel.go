package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

var (
	SENTINEL_COMMAND string = "sentinel get-master-addr-by-name %s\n"
)

type Sentinel struct {
	redis_ip   string
	redis_name string
	sentinel   *net.TCPAddr
}

func (s *Sentinel) read(conn *net.TCPConn) (string, error) {
	err := s.write(conn)

	buffer := make([]byte, 256)

	_, err = conn.Read(buffer)

	parts := strings.Split(string(buffer), "\r\n")

	if err != nil || len(parts) < 5 {
		return "", errors.New("failed to get sentinel")
	}

	return fmt.Sprintf("%s:%s", parts[2], parts[4]), err
}

func (s *Sentinel) write(conn *net.TCPConn) error {
	_, err := conn.Write([]byte(fmt.Sprintf(SENTINEL_COMMAND, s.redis_name)))
	return err
}

func (s *Sentinel) checkMaster(ip string) {
	if len(s.redis_ip) == 0 && s.redis_ip != ip {
		log.Println("new redis master", ip)
		s.redis_ip = ip

		if !NORMAL {
			CLIENTS.state = New
		}
	}
}

func (s *Sentinel) getMaster(conn *net.TCPConn) {
	for {
		if ip, err := s.read(conn); err == nil {
			s.checkMaster(ip)
		} else {
			go s.connect()
			break
		}

		time.Sleep(1 * time.Second)
	}
}

func (s *Sentinel) connect() {
	if conn, err := net.DialTCP("tcp", nil, s.sentinel); err == nil {
		s.getMaster(conn)
	} else {
		time.Sleep(1 * time.Second)
		log.Println("wailed to resolve sentinel", s.sentinel.String())
		s.init(s.sentinel.String())
		go s.connect()
	}
}

func (s *Sentinel) init(sentinel string) {
	adr, err := net.ResolveTCPAddr("tcp", sentinel)
	if err != nil {
		log.Fatal("Failed to resolve sentinel address", err)
	}

	s.sentinel = adr
}

func (s *Sentinel) start(sentinel string) {
	log.Println("starting sentinel")
	s.init(sentinel)
	s.connect()
}
