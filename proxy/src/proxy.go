package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"strings"

	"github.com/mediocregopher/radix/resp/resp2"
)

func closeClient(proxyClient *io.ReadWriteCloser) {
	(*proxyClient).Close()
}

func closeClientError(proxyClient *io.ReadWriteCloser) {
	(*proxyClient).Write([]byte("-Error occured\r\n"))
	(*proxyClient).Close()
}

func createProxy() {
	laddr, err := net.ResolveTCPAddr("tcp", localAddr)
	if err != nil {
		log.Fatal("Failed to resolve local address: ", err)
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		if ready {
			if conn, err := listener.AcceptTCP(); err != nil {
				log.Println(err)
			} else {
				log.Println("new connection from: ", conn.RemoteAddr().(*net.TCPAddr).IP)
				go inPipe(conn)
			}
		}
	}
}

func checkPart(part *string) bool {
	switch strings.ToUpper(*part) {
	case
		"PING",
		"GET",
		"QUIT",
		"EVAL",
		"HGETALL",
		DOCUMENT_CODE,
		PROCEDURE_CODE,
		BUILDING_CODE,
		UTILITY_BUILDING_CODE:
		return true
	}
	return false

}

func makeRequest(command string, args ...string) (*resp2.RawMessage, error) {
	var raw resp2.RawMessage
	var err error
	var b bool
	for index := 0; index < 10; index++ {
		if err, b = doRedisSafe(&raw, command, args...); err == nil || b {
			return &raw, nil
		}
	}
	log.Println(raw, err, "error")
	return nil, err
}

func checkMessage(message *[]byte) (*resp2.RawMessage, error) {
	parts := strings.Split(string(*message), REDIS_STRING_END)
	if len(parts) > 2 && checkPart(&parts[2]) {
		if parts[0] == "*1" {
			return makeRequest(parts[2])
		} else if parts[0] == "*2" {
			return makeRequest(parts[2], parts[4])
		} else if parts[0] == "*3" {
			return makeRequest(parts[2], parts[4], parts[6])
		}
	}
	log.Println(parts)
	return nil, errors.New("error occured when making redis request")
}

func outPipe(respond *resp2.RawMessage, proxyClient *io.ReadWriteCloser) {
	if err := respond.MarshalRESP(*proxyClient); err != nil {
		go closeClient(proxyClient)
	}
}

func redisPipe(proxyClient *io.ReadWriteCloser, message []byte) {
	if respond, err := checkMessage(&message); err == nil {
		go outPipe(respond, proxyClient)
	} else {
		go closeClientError(proxyClient)
	}
}

func inPipe(proxyClient io.ReadWriteCloser) {
	message := make([]byte, 128)
	for {
		_, err := proxyClient.Read(message)

		if err != nil {
			go closeClient(&proxyClient)
			break
		}
		go redisPipe(&proxyClient, bytes.Trim(message, "\x00"))

	}
}
