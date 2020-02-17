package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type ClientsState int

const (
	Start ClientsState = iota
	Old
	New
	Init
)

type Clients struct {
	clients ClientList
	state   ClientsState
}

func (c *Clients) makeClientRequests(clientUrl string) *map[string]interface{} {
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

func (c *Clients) addBiggestValueIntoMap(final *map[string]interface{}, key *string, value int) {
	if finalValue, ok := (*final)[*key]; ok {
		if finalValue.(int) < value {
			(*final)[*key] = value
		}
	} else {
		(*final)[*key] = value
	}
}

func (c *Clients) convertMap(m *map[string]interface{}) *map[string]interface{} {
	var r = make(map[string]interface{})
	for k, v := range *m {
		r[k] = int(v.(float64))
	}
	return &r
}

func (c *Clients) addBiggestValuesIntoMap(final *map[string]interface{}, client *map[string]interface{}) {
	for key, value := range *client {
		switch v := value.(type) {
		case map[string]interface{}:
			m := c.convertMap(&v)
			if finalValue, ok := (*final)[key]; ok {
				finalSubValue := finalValue.(map[string]interface{})
				c.addBiggestValuesIntoMap(&finalSubValue, m)
			} else {
				(*final)[key] = *m
			}
		case float64:
			c.addBiggestValueIntoMap(final, &key, int(v))
		case int:
			c.addBiggestValueIntoMap(final, &key, v)
		default:
			log.Fatal("Wrong type")
		}
	}
}

func (c *Clients) getValuesFromClients() *map[string]interface{} {
	final := make(map[string]interface{})
	for _, client := range c.clients.Clients {
		log.Println("getting values from", client)
		c.addBiggestValuesIntoMap(&final, c.makeClientRequests(client))
	}
	log.Println(final)
	return &final
}

func (c *Clients) changeState() {
	if c.state == Old {
		c.state = New
	}
}

func (c *Clients) addValuesIntoRedis(m *map[string]interface{}) {
	redis := Redis{}
	for k, v := range *m {
		if reflect.ValueOf(v).Kind() == reflect.Map {
			for k2, v2 := range v.(map[string]interface{}) {
				out := make([]byte, 128)
				redis.clientsDo([]byte(
					fmt.Sprintf("*4\r\n$4\r\n%s\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n$%d\r\n%d\r\n", "HSET", len(k), k, len(k2), k2, len(strconv.Itoa(v2.(int))), v2.(int))), out)
				log.Println("HSET", len(k), k, len(k2), k2, len(strconv.Itoa(v2.(int))), v2.(int))

			}
		} else {
			out := make([]byte, 128)
			redis.clientsDo([]byte(
				fmt.Sprintf("*3\r\n$3\r\n%s\r\n$%d\r\n%s\r\n$%d\r\n%d\r\n", "SET", len(k), k, len(strconv.Itoa(v.(int))), v.(int))), out)
			log.Println("SET", len(k), k, len(strconv.Itoa(v.(int))), v.(int))
		}
	}
}

func (c *Clients) start() {
	log.Println("starting clients")
	for {
		if c.state == New {
			c.state = Init
			c.addValuesIntoRedis(c.getValuesFromClients())
			c.state = Old
			ready = true
		}
		time.Sleep(1 * time.Second)
	}
}
