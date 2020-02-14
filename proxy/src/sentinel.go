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
	ready bool = false
)

type Sentinel struct {
	Ip   string
	Port string
	Name string
}

func (s *Sentinel) socket(conn *net.TCPConn) ([]string, error) {
	command := fmt.Sprintf("%s %s%s", "sentinel get-master-addr-by-name", s.Name, REDIS_STRING_END)

	conn.Write([]byte(fmt.Sprintf(command)))

	b := make([]byte, 256)
	_, err := conn.Read(b)

	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(b), REDIS_STRING_END)

	if len(parts) < 5 {
		err = errors.New(string(b))
	}

	return parts, err
}

func (s *Sentinel) getMaster(conn *net.TCPConn) {
	for {
		if parts, err := s.socket(conn); err == nil {
			redisInit(fmt.Sprintf("%s:%s", parts[2], parts[4]))
		} else {
			conn.Close()
			go s.start()
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func (s *Sentinel) getIp() string {
	return s.Ip + ":" + s.Port
}

func (s *Sentinel) createSentinelConnection(sAddr *net.TCPAddr) {
	log.Println("Creating sentinel connection!", sAddr)
	if conn, err := net.DialTCP("tcp", nil, sAddr); err == nil {
		go s.getMaster(conn)
	} else {
		log.Println(err)
		go s.start()
	}
}

func (s *Sentinel) start() {
	log.Println("Creating sentinel!")
	for {
		if sAddr, err := net.ResolveTCPAddr("tcp", s.getIp()); err == nil {
			go s.createSentinelConnection(sAddr)
			break
		} else {
			log.Println(err)
		}
		time.Sleep(1 * time.Second)
	}
}
