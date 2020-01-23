package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Respond struct {
	Respond string
}

func createRoutes() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/building", getBuilding)
	router.HandleFunc("/utilitybuilding", getUtilityBuilding)
	router.HandleFunc("/procedure", getProcedure)
	router.HandleFunc("/document/{doty}", getDocument)
	router.HandleFunc("/readiness", getReadiness)
	go http.ListenAndServe(":8080", router)
}

func getBuilding(w http.ResponseWriter, _ *http.Request) {
	getRespond(&w, "BUILDING_CODE")
}

func getUtilityBuilding(w http.ResponseWriter, _ *http.Request) {
	getRespond(&w, "UTILITY_BUILDING_CODE")
}

func getProcedure(w http.ResponseWriter, _ *http.Request) {
	getRespond(&w, "PROCEDURE_CODE")
}

func getDocument(w http.ResponseWriter, r *http.Request) {
	getRespond(&w, "DOCUMENT_CODE", mux.Vars(r)["doty"])
}

func getReadiness(w http.ResponseWriter, _ *http.Request) {
	if ready {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprint(w, "")
}

func getRespond(w *http.ResponseWriter, command string, args ...string) {
	var redisStr string

	if err, b := makeRedisRequest(&redisStr, command, args...); err == nil || b {
		var respond Respond
		if b {
			respond = Respond{err.Error()}
			(*w).WriteHeader(http.StatusNotFound)
		} else {
			respond = Respond{redisStr}
		}
		json.NewEncoder(*w).Encode(respond)
	} else {
		(*w).WriteHeader(http.StatusInternalServerError)
	}

}
