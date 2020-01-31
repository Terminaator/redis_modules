package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Proxy struct {
	Sentinel Sentinel
	Redis    Redis
	Clients  []string
}

type Sentinel struct {
	Service string
	Port    string
}

type Redis struct {
	Name string
}

func (s Sentinel) GetIp() string {
	return s.Service + ":" + s.Port
}

func readClientList() {
	jsonFile, err := os.Open(conf)

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(bytes, proxy)
}

func initProxy() {
	sentinelAddr := getEnv("sentinel.ip", "127.0.0.1")
	sentinelPort := getEnv("sentinel.port", "26379")
	redisName := getEnv("redis.name", "mymaster")

	BUILDING_CODE = getEnv("redis.building.command", "BUILDING_CODE")
	UTILITY_BUILDING_CODE = getEnv("redis.utility.building.command", "UTILITY_BUILDING_CODE")
	PROCEDURE_CODE = getEnv("redis.procedure.command", "PROCEDURE_CODE")
	DOCUMENT_CODE = getEnv("redis.document.command", "DOCUMENT_CODE")
	YEAR_KEY = getEnv("redis.year.key", "YEAR_KEY")
	TOKEN = getEnv("redis.token", "d29bbdfd9f7c2b46d142590330f28ef9029da92a83c947b57924504fd7f4abc092a550eb868f3ebf2d7f152690de26e975cf991e7b5d47bdeabf8990c89d09ed32ee8e18ca7ae62d13fd302cfc2683c5e39e398c38cf2b0e82f7ff764b30a8af587b651a")

	readClientList()

	proxy.Sentinel = Sentinel{sentinelAddr, sentinelPort}
	proxy.Redis = Redis{redisName}

	log.Println(proxy)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
