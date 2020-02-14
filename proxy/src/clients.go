package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
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

func addBiggestValueIntoMap(final *map[string]interface{}, key *string, value int) {
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

func addBiggestValuesIntoMap(final *map[string]interface{}, client *map[string]interface{}) {
	for key, value := range *client {
		switch v := value.(type) {
		case map[string]interface{}:
			m := convertMap(&v)
			if finalValue, ok := (*final)[key]; ok {
				finalSubValue := finalValue.(map[string]interface{})
				addBiggestValuesIntoMap(&finalSubValue, m)
			} else {
				(*final)[key] = *m
			}
		case float64:
			addBiggestValueIntoMap(final, &key, int(v))
		case int:
			addBiggestValueIntoMap(final, &key, v)
		default:
			log.Fatal("Wrong type")
		}
	}
}

func getValuesFromClients() *map[string]interface{} {
	final := make(map[string]interface{})
	for _, client := range CLIENTS.Clients {
		log.Println(client)
		addBiggestValuesIntoMap(&final, makeClientRequests(client))
	}
	log.Println(final)
	return &final
}

func addValuesIntoRedis(redis Redis, m *map[string]interface{}) {
	for k, v := range *m {
		if reflect.ValueOf(v).Kind() == reflect.Map {
			for k2, v2 := range v.(map[string]interface{}) {
				redis.doInit([]byte(
					fmt.Sprintf("*4\r\n$4\r\n%s\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n$%d\r\n%d\r\n", "HSET", len(k), k, len(k2), k2, len(strconv.Itoa(v2.(int))), v2.(int))))
			}
		} else {
			redis.doInit([]byte(
				fmt.Sprintf("*3\r\n$3\r\n%s\r\n$%d\r\n%s\r\n$%d\r\n%d\r\n", "SET", len(k), k, len(strconv.Itoa(v.(int))), v.(int))))
		}
	}
}
