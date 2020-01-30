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
	sentinelAddr := getEnv("sentinel.ip", "sentinel")
	sentinelPort := getEnv("sentinel.port", "26379")
	redisName := getEnv("redis.name", "redis")

	BUILDING_CODE = getEnv("redis.building.command", "BUILDING_CODE")
	UTILITY_BUILDING_CODE = getEnv("redis.utility.building.command", "UTILITY_BUILDING_CODE")
	PROCEDURE_CODE = getEnv("redis.procedure.command", "PROCEDURE_CODE")
	DOCUMENT_CODE = getEnv("redis.document.command", "DOCUMENT_CODE")
	DOCUMENT_CODE = getEnv("redis.year.key", "YEAR_KEY")

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
