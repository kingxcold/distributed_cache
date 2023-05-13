package main

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

	// listenAddr := flag.String("listenaddr", ":3000", "listen address of the server")
	// leaderAddr := flag.String("leaderaddr", "", "listen address of theleader ")
	// flag.Parse()
	// opts := ServerOpts{
	// 	ListenAddr: *listenAddr,
	// 	IsLeader:   len(*leaderAddr) == 0,
	// 	LeaderAddr: *leaderAddr,
	// }
	// server := NewServer(opts, cache.New())

	// log.Fatal(server.Start())
}
