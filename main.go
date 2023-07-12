package main

import (
	"context"
	"distributed_cache/cache"
	"distributed_cache/client"
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	// conn, err := net.Dial("tcp", ":3000")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// _, err = conn.Write([]byte("SET hello world 50000000000"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// select {}
	// return

	// go func() {
	// 	time.Sleep(time.Second * 2)
	// 	conn, err := net.Dial("tcp", *listenAddr)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	conn.Write([]byte("SET hello world 50000000000"))

	// 	time.Sleep(time.Second * 1)
	// 	conn.Write([]byte("GET hello"))
	// 	buf := make([]byte, 1024)
	// 	n, _ := conn.Read(buf)
	// 	fmt.Println(string(buf[:n]))
	// }()

	listenAddr := flag.String("listenaddr", ":3000", "listen address of the server")
	leaderAddr := flag.String("leaderaddr", "", "listen address of theleader ")
	flag.Parse()
	opts := ServerOpts{
		ListenAddr: *listenAddr,
		IsLeader:   len(*leaderAddr) == 0,
		LeaderAddr: *leaderAddr,
	}

	go func() {
		time.Sleep(time.Second * 10)
		if opts.IsLeader {
			sendDummyData()
		}
	}()

	server := NewServer(opts, cache.New())
	log.Fatal(server.Start())
}
func sendDummyData() {
	for i := 0; i < 100; i++ {
		go func(i int) {

			key := []byte(fmt.Sprintf("key_%d", i))
			value := []byte(fmt.Sprintf("val_%d", i))

			client, err := client.New(":3000", client.Options{})
			if err != nil {
				log.Fatal(err)
			}

			err = client.Set(context.Background(), key, value, 0)
			if err != nil {
				log.Fatal(err)
			}

			time.Sleep(time.Second * 1)
			_, err = client.Get(context.Background(), key)
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Println(string(val))
			client.Close()
		}(i)
	}
}
