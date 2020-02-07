package main

import "log"

var (
	values Init
	router Router
)

func main() {
	log.Println("Starting proxy server")

	values = Init{}
	values.initValues()

	log.Println(values)

	router = Router{Port: ":8080"}
	router.start(values.Token)

	values.Sentinel.start()

	createProxy()
}
