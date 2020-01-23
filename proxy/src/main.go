package main

import (
	"log"

	"github.com/mediocregopher/radix"
)

var (
	pool  *radix.Pool
	proxy *Proxy = new(Proxy)

	redisAddr string

	ready     bool   = false
	conf      string = "./conf/init.json"
	localAddr string = ":9999"

	REDIS_STRING_END string = "\r\n"
)

func main() {
	log.Println("Starting proxy server")

	initProxy()
	go createRoutes()
	go createSentinel()
	createProxy()

}
