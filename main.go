package main

import (
	"distributed_cache/cache"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	opts := ServerOpts{ListenAddr: ":3000", IsLeader: true}
	server := NewServer(opts, cache.New())

	go func() {
		time.Sleep(time.Second * 2)
		conn, err := net.Dial("tcp", ":3000")
		if err != nil {
			log.Fatal(err)
		}
		conn.Write([]byte("SET hello world 50000000000"))

		time.Sleep(time.Second * 1)
		conn.Write([]byte("GET hello"))
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		fmt.Println(string(buf[:n]))
	}()

	log.Fatal(server.Start())
}
