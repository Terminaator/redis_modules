package main

import (
	"log"

	"github.com/mediocregopher/radix"
)

type Redis struct {
	adr  string
	Pool *radix.Pool
}

func (r *Redis) init(p *radix.Pool) *radix.Pool {
	log.Println("adding values")
	return addValuesIntoRedis(p, getValuesFromClients(values.Clients.Clients))
}

func (r *Redis) createPool() {
	log.Println("creating pool")
	if p, err := radix.NewPool("tcp", r.adr, 10); err == nil {
		r.Pool = r.init(p)
		values.Ready = true
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
