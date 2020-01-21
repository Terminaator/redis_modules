package main

import (
	"fmt"
	"log"
	"net/http"
)

func createReadiness() {
	http.Handle("/ready", http.HandlerFunc(readiness))
	go http.ListenAndServe(":8080", nil)
}

func readiness(w http.ResponseWriter, _ *http.Request) {
	if ready {
		log.Println("Okay")
		w.WriteHeader(http.StatusOK)
	} else {
		log.Println("Not okay")
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprint(w, "")
}
