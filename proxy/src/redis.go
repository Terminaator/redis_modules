package main

import (
	"errors"
	"log"

	"github.com/mediocregopher/radix"
	"github.com/mediocregopher/radix/resp/resp2"
)

func createRedisPool() {
	p, err := radix.NewPool("tcp", redisAddr, 10)

	if err != nil {
		log.Fatal("Failed to create a pool", err)
	} else {
		pool = p
		initRedis()
	}
}

func cmdRedis(resp interface{}, command string, args ...string) radix.CmdAction {
	return radix.Cmd(resp, command, args...)
}

func doRedis(resp interface{}, command string, args ...string) error {
	return pool.Do(cmdRedis(resp, command, args...))
}

func doRedisSafe(resp interface{}, command string, args ...string) (error, bool) {
	var redisErr resp2.Error

	err := doRedis(resp, command, args...)

	if errors.As(err, &redisErr) {
		log.Println(err)
		return err, true
	} else {
		return err, false
	}
}
