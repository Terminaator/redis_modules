package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/sparrc/go-ping"
)

func sentinelSocket(sentinel *net.TCPConn) ([]string, error) {
	command := fmt.Sprintf("%s %s%s", "sentinel get-master-addr-by-name", proxy.Redis.Name, REDIS_STRING_END)

	sentinel.Write([]byte(fmt.Sprintf(command)))

	b := make([]byte, 256)
	_, err := sentinel.Read(b)

	if err != nil {
		defer sentinel.Close()
	}

	parts := strings.Split(string(b), REDIS_STRING_END)

	if len(parts) < 5 {
		defer sentinel.Close()
		err = errors.New(string(b))
	}

	return parts, err
}

func createNewPool(stringaddr *string) {
	redisAddr = *stringaddr
	createRedisPool()
}

func getRedisMaster(sentinel *net.TCPConn) {
	for {
		if parts, err := sentinelSocket(sentinel); err == nil {
			stringaddr := fmt.Sprintf("%s:%s", parts[2], parts[4])

			if len(redisAddr) == 0 {
				log.Println("New master:", stringaddr)
				go createNewPool(&stringaddr)
			} else if stringaddr != redisAddr {
				log.Println("Master has changed! New:", stringaddr, "old:", redisAddr)
				go createNewPool(&stringaddr)
			}
		} else {
			log.Println(err)
			go createSentinel()
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func createSentinelConnection(sentinel *net.TCPAddr) {
	log.Println("Creating sentinel connection!")
	if conn, err := net.DialTCP("tcp", nil, sentinel); err == nil {
		go getRedisMaster(conn)
	} else {
		go createSentinel()
	}
}

func pingService(service string) bool {
	log.Println("Pinging: ", service)
	if _, err := ping.NewPinger(service); err == nil {
		return true
	} else {
		return false
	}
}

func createSentinel() {
	log.Println("Creating sentinel!")
	for {
		if pingService(proxy.Sentinel.Service) {
			if sAddr, err := net.ResolveTCPAddr("tcp", proxy.Sentinel.GetIp()); err == nil {
				go createSentinelConnection(sAddr)
				break
			}
		}
		time.Sleep(1 * time.Second)
	}
}
