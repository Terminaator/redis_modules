package main

import (
	"fmt"
	"net/http"
)

func readiness(w http.ResponseWriter, _ *http.Request) {
	if ready {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprint(w, "")
}
