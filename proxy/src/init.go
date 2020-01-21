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
	flag.Parse()
	readFile()
}
