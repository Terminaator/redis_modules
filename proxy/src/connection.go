package main

import (
	"bytes"
	"io"
	"log"
	"strings"
)

type Connection struct {
	ip    string
	rwc   io.ReadWriteCloser
	redis Redis
}

func (c *Connection) close() {
	log.Println("connection closed", c.ip)
	c.redis.close()
	c.rwc.Close()
}

func (c *Connection) outPipe(out []byte) {
	log.Println("message out", c.ip, out)
	c.rwc.Write(out)
}

func (c *Connection) outCheckPipe(out []byte) {
	if !NORMAL {
		if len(out) > 0 && strings.Contains(string(out), failError) {
			log.Println("init needed")
			CLIENTS.changeState()
		}
	}

	c.outPipe(out)
}

func (c *Connection) redisPipe(message []byte) {
	log.Println("message from", c.ip, message)
	out := make([]byte, 2048)
	c.redis.normalDo(message, out)

	c.outCheckPipe(bytes.Trim(out, "\x00"))
}

func (c *Connection) inPipe() {
	for {
		buffer := make([]byte, 1024)

		if _, err := c.rwc.Read(buffer); err == nil {
			go c.redisPipe(bytes.Trim(buffer, "\x00"))
		} else {
			c.close()
			break
		}

	}
}
