package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	CONF_FILE string = "./conf/init.json"
)

type Init struct {
	Sentinel Sentinel
	Redis    Redis
	Clients  Clients
	Keys     Keys
	Ready    bool
	Timeout  bool
	Token    string
}

type Keys struct {
	BUILDING_CODE         string
	UTILITY_BUILDING_CODE string
	PROCEDURE_CODE        string
	DOCUMENT_CODE         string
}

func (i *Init) initValues() {
	log.Println("adding values")
	i.Ready = false
	i.Timeout = false
	i.Keys = Keys{
		BUILDING_CODE:         getEnv("redis.building.command", "BUILDING_CODE"),
		UTILITY_BUILDING_CODE: getEnv("redis.utility.building.command", "UTILITY_BUILDING_CODE"),
		PROCEDURE_CODE:        getEnv("redis.procedure.command", "PROCEDURE_CODE"),
		DOCUMENT_CODE:         getEnv("redis.document.command", "DOCUMENT_CODE")}

	i.Sentinel = Sentinel{
		Ip:   getEnv("sentinel.ip", "127.0.0.1"),
		Port: getEnv("sentinel.port", "26379"),
		Name: getEnv("redis.name", "mymaster")}

	i.Redis = Redis{}
	i.Clients = Clients{File: CONF_FILE}
	i.Clients.readFile()
	i.Token = getEnv("redis.token", "d29bbdfd9f7c2b46d142590330f28ef9029da92a83c947b57924504fd7f4abc092a550eb868f3ebf2d7f152690de26e975cf991e7b5d47bdeabf8990c89d09ed32ee8e18ca7ae62d13fd302cfc2683c5e39e398c38cf2b0e82f7ff764b30a8af587b651a")
}

type Clients struct {
	File    string
	Clients []string
}

func (i *Clients) readFile() {
	jsonFile, err := os.Open(i.File)

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(bytes, i)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
