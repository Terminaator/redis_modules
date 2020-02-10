package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/mediocregopher/radix"
)

func makeClientRequests(clientUrl string) *map[string]interface{} {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(clientUrl)

	if err != nil || resp.StatusCode != 200 {
		log.Fatal("Error occured when getting response from client or statuscode is not 200")
	}
	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	return &result

}

func addValueToMap(final *map[string]interface{}, key *string, value int) {
	if finalValue, ok := (*final)[*key]; ok {
		if finalValue.(int) < value {
			(*final)[*key] = value
		}
	} else {
		(*final)[*key] = value
	}
}

func convertMap(m *map[string]interface{}) *map[string]interface{} {
	var r = make(map[string]interface{})
	for k, v := range *m {
		r[k] = int(v.(float64))
	}
	return &r
}

func addValuesToMap(final *map[string]interface{}, client *map[string]interface{}) {
	for key, value := range *client {
		switch v := value.(type) {
		case map[string]interface{}:
			m := convertMap(&v)
			if finalValue, ok := (*final)[key]; ok {
				finalSubValue := finalValue.(map[string]interface{})
				addValuesToMap(&finalSubValue, m)
			} else {
				(*final)[key] = *m
			}
		case float64:
			addValueToMap(final, &key, int(v))
		case int:
			addValueToMap(final, &key, v)
		default:
			log.Fatal("Wrong type")
		}
	}
}

func addValuesIntoRedis(p *radix.Pool, m *map[string]interface{}) *radix.Pool {
	for k, v := range *m {
		if reflect.ValueOf(v).Kind() == reflect.Map {
			for k2, v2 := range v.(map[string]interface{}) {
				if err := p.Do(radix.Cmd(nil, "HSET", k, k2, strconv.Itoa(v2.(int)))); err != nil {
					log.Fatal("failed adding value")
				}
			}
		} else {
			if err := p.Do(radix.Cmd(nil, "SET", k, strconv.Itoa(v.(int)))); err != nil {
				log.Fatal("failed adding value")
			}
		}
	}
	return p
}

func getValuesFromClients(clients []string) *map[string]interface{} {
	final := make(map[string]interface{})
	for _, client := range clients {
		addValuesToMap(&final, makeClientRequests(client))
	}
	log.Println(final)
	return &final
}
