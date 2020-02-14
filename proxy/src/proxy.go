package main

import (
	"bytes"
	"io"
	"log"
	"net"
)

type Proxy struct {
	adr string
}

type Connection struct {
	ip    string
	rwc   io.ReadWriteCloser
	redis Redis
}

func (c *Connection) read() ([]byte, error) {
	message := make([]byte, 128)

	if _, err := c.rwc.Read(message); err == nil {
		return bytes.Trim(message, "\x00"), err
	} else {
		return nil, err
	}
}

func (c *Connection) closeOnError(err error) {
	log.Println("closing connection", c.ip, err)
	c.rwc.Close()
	c.redis.close()
}

func (p *Proxy) outPipe(conn Connection, out []byte) {
	log.Println("out message", string(out))

	_, err := conn.rwc.Write(out)

	if err != nil {
		conn.closeOnError(err)
	}

	if len(out) > 0 && out[0] == '-' {
		log.Println("init needed")
		host = ""
	}
}

func (p *Proxy) middlePipe(conn Connection, in []byte) {
	out := make([]byte, 4096)
	conn.redis.do(in, out)
	p.outPipe(conn, bytes.Trim(out, "\x00"))
}

func (p *Proxy) inPipe(conn Connection) {
	log.Println("starting to read from", conn.ip)
	conn.redis.start()
	for {
		if ready {
			if message, err := conn.read(); err == nil {
				log.Println("from", conn.ip, "message", string(message))
				p.middlePipe(conn, message)
			} else {
				go conn.closeOnError(err)
				break
			}
		}
	}
}

func (p *Proxy) start() {
	laddr, err := net.ResolveTCPAddr("tcp", p.adr)
	if err != nil {
		log.Fatal("Failed to resolve local address: ", err)
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		if conn, err := listener.AcceptTCP(); err == nil {
			log.Println("new connection from", conn.RemoteAddr().(*net.TCPAddr).IP.String())
			go p.inPipe(Connection{ip: conn.RemoteAddr().(*net.TCPAddr).IP.String(), rwc: conn, redis: Redis{}})
		}
	}
}
