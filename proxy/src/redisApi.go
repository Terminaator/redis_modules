package main

import (
	"log"
)

func makeRedisRequest(val interface{}, command string, args ...string) (error, bool) {
	log.Println(command, args)
	return doRedisSafe(val, command, args...)
}
