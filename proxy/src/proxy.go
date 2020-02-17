package main

import (
	"io"
	"log"
	"net"
	"strconv"
)

type Proxy struct {
	ip string
}

func (p *Proxy) init() *net.TCPListener {

	laddr, err := net.ResolveTCPAddr("tcp", p.ip)
	if err != nil {
		log.Fatal("Failed to resolve local address", err)
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatal(err)
	}

	return listener
}

func (p *Proxy) connection(ip string, port int, rwc io.ReadWriteCloser) {
	log.Println("new connection from:", ip, port)
	conn := Connection{ip: ip + ":" + strconv.Itoa(port), rwc: rwc, redis: Redis{state: Connect}}
	go conn.inPipe()
}

func (p *Proxy) start() {
	listener := p.init()

	for {
		if c, err := listener.AcceptTCP(); err == nil {
			go p.connection(
				c.RemoteAddr().(*net.TCPAddr).IP.String(),
				c.RemoteAddr().(*net.TCPAddr).Port,
				c)
		}
	}
}
