package main

import (
	"encoding/json"
	"flag"
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

func readFile() {
	jsonFile, err := os.Open(conf)

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(bytes, proxy)
}

func initProxy() {
	sentinelAddr := getEnv("service", "127.0.0.1")
	sentinelPort := getEnv("port", "26379")
	redisName := getEnv("name", "mymaster")

	flag.Parse()
	readFile()

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
