package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/mediocregopher/radix/resp/resp2"
)

var (
	localAddr string = ":9999"
)

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
		if conn, err := listener.AcceptTCP(); err == nil {
			go inPipe(conn)
		} else {
			log.Println(err)
		}
	}
}

func validPart(command string, p []string) bool {
	command = strings.ToUpper(command)
	switch command {
	case
		"GET",
		"HGETALL",
		"PING",
		"QUIT",
		values.Keys.DOCUMENT_CODE,
		values.Keys.PROCEDURE_CODE,
		values.Keys.BUILDING_CODE,
		values.Keys.UTILITY_BUILDING_CODE:
		return true
	}

	if command == "EVAL" {
		eval := strings.ToUpper(p[0])
		if strings.Contains(eval, "RETURN REDIS.CALL('BUILDING_CODE')") {
			return true
		} else if strings.Contains(eval, "RETURN REDIS.CALL('UTILITY_BUILDING_CODE')") {
			return true
		} else if strings.Contains(eval, "RETURN REDIS.CALL('PROCEDURE_CODE')") {
			return true
		} else if strings.Contains(eval, "RETURN REDIS.CALL('DOCUMENT_CODE") {
			return true
		}
	}
	return false
}

func checkMessage(message *[]byte) (string, []string, error) {
	parts := strings.Split(string(*message), REDIS_STRING_END)
	log.Println("parts", parts)

	if len(parts) > 2 {
		x, err := strconv.Atoi(parts[0][1:])

		if err != nil {
			return "", parts, errors.New("Not valid")
		}

		var out []string

		for i := 0; i < x; i++ {
			out = append(out, parts[2+i*2])
		}

		if validPart(out[0], out[1:]) {
			return out[0], out[1:], nil
		} else {
			return "", parts, errors.New("Not valid")
		}
	} else {
		return "", parts, errors.New("Not valid")
	}
}

func outValidPipe(proxyClient *io.ReadWriteCloser, command string, parts []string) {
	var raw resp2.RawMessage
	var redisErr resp2.Error

	err := values.Redis.doRedis(&raw, command, parts...)

	if err != nil && !errors.As(err, &redisErr) {
		(*proxyClient).Write([]byte("-Error occured\r\n"))
	} else {
		raw.MarshalRESP(*proxyClient)
	}
}

func outNotValidPipe(proxyClient *io.ReadWriteCloser) {
	(*proxyClient).Write([]byte("-Not valid message\r\n"))
	(*proxyClient).Close()
}

func middlePipe(proxyClient *io.ReadWriteCloser, message []byte) {
	if command, parts, err := checkMessage(&message); err == nil {
		go outValidPipe(proxyClient, command, parts)
	} else {
		go outNotValidPipe(proxyClient)
	}
}

func inPipe(proxyClient io.ReadWriteCloser) {
	for {
		message := make([]byte, 128)
		if _, err := proxyClient.Read(message); err == nil {
			go middlePipe(&proxyClient, bytes.Trim(message, "\x00"))
		} else {
			proxyClient.Close()
			break
		}
	}
}
