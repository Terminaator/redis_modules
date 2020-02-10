package main

import (
	"errors"
	"log"

	"github.com/mediocregopher/radix"
)

type Redis struct {
	adr  string
	Pool *radix.Pool
}

func (r *Redis) init() {
	log.Println("adding values")
	addValuesIntoRedis(getValuesFromClients(values.Clients.Clients))
	values.Ready = true
}

func (r *Redis) createPool() {
	log.Println("creating pool")
	if p, err := radix.NewPool("tcp", r.adr, 10); err == nil {
		r.Pool = p
		r.init()
	} else {
		p.Close()
		log.Fatal("Failed to create a pool", err)
	}
}

func (r *Redis) addAdr(adr string) {
	if len(r.adr) == 0 {
		log.Println("New master", adr)
		r.adr = adr
		r.createPool()
	} else if r.adr != adr {
		log.Println("New master", adr)
		r.Pool.Close()
		r.createPool()
	}
}

func (r *Redis) doRedis(resp interface{}, command string, args ...string) error {
	return r.Pool.Do(radix.Cmd(resp, command, args...))
}

func (r *Redis) doRedisReady(resp interface{}, command string, args ...string) error {
	if values.Ready {
		return r.Pool.Do(radix.Cmd(resp, command, args...))
	} else {
		return errors.New("Not ready")
	}
}
