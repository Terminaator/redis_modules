package main

import (
	"log"

	"github.com/mediocregopher/radix"
)

type Redis struct {
	adr  string
	Pool *radix.Pool
}

func (r *Redis) init() {
	resp := radix.MaybeNil{}
	if err := r.doRedis(&resp, "GET", values.Keys.YEAR_KEY); err == nil && resp.Nil {
		log.Println("init needed")
		addValuesIntoRedis(getValuesFromClients(values.Clients.Clients))
		values.Ready = true
	} else if err != nil {
		log.Fatal("Redis init failed!")
	} else {
		log.Println("init not needed")
		values.Ready = true
	}
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
