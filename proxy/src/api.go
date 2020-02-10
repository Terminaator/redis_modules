package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	Port   string
	Auth   Auth
	Router *mux.Router
}

type Auth struct {
	Token string
}

type Respond struct {
	Respond string
}

func (a *Auth) authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")
		if token == a.Token || "/readiness" == r.URL.Path {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func (a *Auth) middleware(r *mux.Router) {
	log.Println("adding token check")
	r.Use(a.authentication)
}

func (r *Router) start(token string) {
	log.Println("starting router")
	r.Router = mux.NewRouter().StrictSlash(true)

	auth := Auth{Token: token}
	auth.middleware(r.Router)
	r.Auth = auth

	r.addingRoutes()
}

func (r *Router) addingRoutes() {
	log.Println("adding routes")
	r.Router.HandleFunc("/building", getBuilding)
	r.Router.HandleFunc("/utilitybuilding", getUtilityBuilding)
	r.Router.HandleFunc("/procedure", getProcedure)
	r.Router.HandleFunc("/document/{doty}", getDocument)
	r.Router.HandleFunc("/readiness", getReadiness)
	go http.ListenAndServe(r.Port, r.Router)
}

func getBuilding(w http.ResponseWriter, _ *http.Request) {
	get(&w, values.Keys.BUILDING_CODE)
}

func getUtilityBuilding(w http.ResponseWriter, _ *http.Request) {
	get(&w, values.Keys.UTILITY_BUILDING_CODE)
}

func getProcedure(w http.ResponseWriter, _ *http.Request) {
	get(&w, values.Keys.PROCEDURE_CODE)
}

func getDocument(w http.ResponseWriter, r *http.Request) {
	get(&w, values.Keys.DOCUMENT_CODE, mux.Vars(r)["doty"])
}

func getReadiness(w http.ResponseWriter, _ *http.Request) {
	if values.Ready {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprint(w, "")
}

func get(w *http.ResponseWriter, command string, args ...string) {
	log.Println("api", command, args)

	var redisStr string
	var respond Respond

	if err := values.Redis.doRedis(&redisStr, command, args...); err == nil {
		respond = Respond{redisStr}
	} else {
		respond = Respond{err.Error()}
		(*w).WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(*w).Encode(respond)
}
