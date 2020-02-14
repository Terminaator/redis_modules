package main

import (
	"log"
	"net"
	"time"
)

var (
	host string
)

type Redis struct {
	remote *net.TCPConn
}

func (r *Redis) close() {
	log.Println("ending session")
	r.remote.Write([]byte("*1\r\n$4\r\nQUIT\r\n"))
}

func (r *Redis) read(in []byte, out []byte) {
	_, err := r.remote.Read(out)

	if err != nil {
		r.start()
		r.write(in, out)
	}
}

func (r *Redis) write(in []byte, out []byte) {
	_, err := r.remote.Write(in)

	if err != nil {
		r.start()
		r.write(in, out)
	} else {
		r.read(in, out)
	}

}

func (r *Redis) doInit(in []byte) {
	_, err := r.remote.Write(in)
	if err != nil {
		log.Fatal("writing failed")
	}
}

func (r *Redis) do(in []byte, out []byte) {
	r.write(in, out)
}

func (r *Redis) startWithIp(ip string) {
	addr, _ := net.ResolveTCPAddr("tcp", ip)
	remote, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		log.Fatal("crash")
	} else {
		r.remote = remote
	}
}

func (r *Redis) start() {
	addr, _ := net.ResolveTCPAddr("tcp", host)

	remote, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		log.Println("failed to create redis")
		time.Sleep(1 * time.Second)
		r.start()
	} else {
		r.remote = remote
	}
}

func redisInit(ip string) {
	if !ready || ip != host {
		redis := Redis{}
		redis.startWithIp(ip)
		addValuesIntoRedis(redis, getValuesFromClients())
		redis.close()

		log.Println("new host", ip)
		host = ip
		ready = true
	}
}
