package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var (
	ONE_COMMAND               string = "*1\r\n$%d\r\n%s\r\n"
	ONE_COMMAND_WITH_VARIABLE string = "*2\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n"
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
		log.Println(r.URL.Path, "req")
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
	respond(&w, fmt.Sprintf(ONE_COMMAND, len(KEYS.BUILDING_CODE), KEYS.BUILDING_CODE))
}

func getUtilityBuilding(w http.ResponseWriter, _ *http.Request) {
	respond(&w, fmt.Sprintf(ONE_COMMAND, len(KEYS.UTILITY_BUILDING_CODE), KEYS.UTILITY_BUILDING_CODE))
}

func getProcedure(w http.ResponseWriter, _ *http.Request) {
	respond(&w, fmt.Sprintf(ONE_COMMAND, len(KEYS.PROCEDURE_CODE), KEYS.PROCEDURE_CODE))
}

func getDocument(w http.ResponseWriter, r *http.Request) {
	respond(&w, fmt.Sprintf(ONE_COMMAND_WITH_VARIABLE, len(KEYS.DOCUMENT_CODE), KEYS.DOCUMENT_CODE, len(mux.Vars(r)["doty"]), mux.Vars(r)["doty"]))
}

func getReadiness(w http.ResponseWriter, _ *http.Request) {
	if ready {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprint(w, "")
}

func convert(out string) (string, error) {
	if len(out) > 0 {
		init := strings.Contains(out, failError)

		if init {
			log.Println("init needed")
			CLIENTS.changeState()
		} else if out[0] == ':' || out[0] == '+' {
			return strings.Trim(out[1:], "\r\n"), nil
		}

		return "", errors.New(out[1:])

	}
	return "", errors.New("sth went wrong")
}

func doRedis(command []byte) (string, error) {
	log.Println("api", command)
	redis := Redis{}
	out := make([]byte, 128)

	redis.normalDo(command, &out)
	redis.close()

	return convert(string(bytes.Trim(out, "\x00")))

}

func respond(w *http.ResponseWriter, command string) {
	var res Respond
	r, err := doRedis([]byte(command))

	if err == nil {
		res = Respond{r}
	} else {
		res = Respond{err.Error()}
		(*w).WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(*w).Encode(res)
}
