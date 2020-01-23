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
	sentinelAddr := getEnv("service", "sentinel")
	sentinelPort := getEnv("port", "26379")
	redisName := getEnv("name", "redis")

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
