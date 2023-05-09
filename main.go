package main

import (
	"distributed_cache/cache"
	"log"
)

func main() {
	opts := ServerOpts{ListenAddr: ":3000", IsLeader: true}
	server := NewServer(opts, cache.New())
	log.Fatal(server.Start())
}
