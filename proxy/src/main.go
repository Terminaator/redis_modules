package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	SENTINEL Sentinel
	PROXY    Proxy
	CLIENTS  Clients
	ROUTER   Router
	KEYS     Keys
	NORMAL   bool

	failError string = "Error occurred when getting value"
	ready     bool   = false
	CONF_FILE string = "./conf/init.json"
)

func main() {

	list, token, port, sentinel, name := addInitValues()

	log.Println("starting proxy:", port)

	SENTINEL = Sentinel{redis_name: name}

	go SENTINEL.start(sentinel)

	if !NORMAL {
		CLIENTS = Clients{state: Start, clients: list}

		go CLIENTS.start()

		ROUTER = Router{Port: ":8080"}

		go ROUTER.start(token)
	}

	PROXY = Proxy{ip: port}

	PROXY.start()
}

type Keys struct {
	BUILDING_CODE         string
	UTILITY_BUILDING_CODE string
	PROCEDURE_CODE        string
	DOCUMENT_CODE         string
}

type ClientList struct {
	File    string
	Clients []string
}

func (i *ClientList) readFile() {
	jsonFile, err := os.Open(i.File)

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(bytes, i)
}

func addInitValues() (ClientList, string, string, string, string) {
	log.Println("adding values proxy")
	n := getEnv("redis.type.normal", "false")
	if n == "true" {
		NORMAL = true
	} else if n == "false" {
		NORMAL = false
	} else {
		log.Fatal("Wrong proxy type")
	}
	var list ClientList
	if !NORMAL {
		list = ClientList{File: CONF_FILE}
		list.readFile()

		KEYS = Keys{
			BUILDING_CODE:         getEnv("redis.building.command", "BUILDING_CODE"),
			UTILITY_BUILDING_CODE: getEnv("redis.utility.building.command", "UTILITY_BUILDING_CODE"),
			PROCEDURE_CODE:        getEnv("redis.procedure.command", "PROCEDURE_CODE"),
			DOCUMENT_CODE:         getEnv("redis.document.command", "DOCUMENT_CODE")}
	}

	return list,
		getEnv("redis.token", "d29bbdfd9f7c2b46d142590330f28ef9029da92a83c947b57924504fd7f4abc092a550eb868f3ebf2d7f152690de26e975cf991e7b5d47bdeabf8990c89d09ed32ee8e18ca7ae62d13fd302cfc2683c5e39e398c38cf2b0e82f7ff764b30a8af587b651a"),
		getEnv("redis.socket.port", ":9999"),
		getEnv("sentinel.ip", "127.0.0.1") + ":" + getEnv("sentinel.port", "26379"),
		getEnv("redis.name", "mymaster")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
