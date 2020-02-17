package main

import (
	"log"
	"net"
	"time"
)

var (
	quit []byte = []byte("*1\r\n$4\r\nquit\r\n")
)

type RedisState int

const (
	Connect RedisState = iota
	Open
	Close
)

type Redis struct {
	conn  *net.TCPConn
	state RedisState
}

func (r *Redis) close() {
	r.state = Close
	r.conn.Write(quit)
}

func (r *Redis) read(in []byte, out []byte) {
	if _, err := r.conn.Read(out); err != nil {
		r.state = Connect
		r.do(in, out)
	}
}

func (r *Redis) write(in []byte, out []byte) {
	if _, err := r.conn.Write(in); err == nil {
		r.read(in, out)
	} else {
		r.state = Connect
		r.do(in, out)
	}
}

func (r *Redis) clientsDo(in []byte, out []byte) {
	r.do(in, out)
}

func (r *Redis) normalDo(in []byte, out []byte) {
	if !NORMAL {
		if CLIENTS.state != Old {
			for {
				if CLIENTS.state == Old {
					break
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
	r.do(in, out)
}

func (r *Redis) do(in []byte, out []byte) {
	if r.state == Connect {
		r.connect()
		r.do(in, out)
	} else if r.state == Open {
		r.write(in, out)
	}
}

func (r *Redis) connect() {
	if r.state != Close {
		addr, _ := net.ResolveTCPAddr("tcp", SENTINEL.ip)
		c, err := net.DialTCP("tcp", nil, addr)

		if err == nil {
			r.state = Open
			r.conn = c
		} else {
			log.Println("trying to make new redis session")
			time.Sleep(1 * time.Second)
			r.connect()
		}
	}
}
