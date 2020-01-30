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

	BUILDING_CODE         string
	UTILITY_BUILDING_CODE string
	PROCEDURE_CODE        string
	DOCUMENT_CODE         string
	YEAR_KEY              string

	REDIS_STRING_END string = "\r\n"
)

func main() {
	log.Println("Starting proxy server")

	initProxy()
	go createRoutes()
	go createSentinel()
	createProxy()

}
