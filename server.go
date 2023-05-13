package main

// import (
// 	"context"
// 	"distributed_cache/cache"
// 	"fmt"
// 	"log"
// 	"net"
// )

// type ServerOpts struct {
// 	ListenAddr string
// 	LeaderAddr string
// 	IsLeader   bool
// }

// type Server struct {
// 	ServerOpts
// 	followers map[net.Conn]struct{}
// 	cache     cache.Cacher
// }

// func NewServer(opts ServerOpts, c cache.Cacher) *Server {
// 	return &Server{ServerOpts: opts, cache: c, followers: make(map[net.Conn]struct{})}
// }

// func (s *Server) Start() error {
// 	ln, err := net.Listen("tcp", s.ListenAddr)
// 	if err != nil {
// 		return fmt.Errorf("listen error: %s", err)
// 	}

// 	log.Printf("server starting on port [%s]\n", s.ListenAddr)

// 	if !s.IsLeader {
// 		go func() {
// 			conn, err := net.Dial("tcp", s.LeaderAddr)
// 			fmt.Println("connected with leader: ", s.LeaderAddr)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			s.HandleConn(conn)
// 		}()
// 	}

// 	for {
// 		conn, err := ln.Accept()
// 		log.Printf("new connection [%s]\n", conn.RemoteAddr())
// 		if err != nil {
// 			log.Printf("accept error %s\n", err)
// 			continue
// 		}
// 		go s.HandleConn(conn)
// 	}
// }

// func (s *Server) HandleConn(conn net.Conn) {
// 	defer conn.Close()
// 	buf := make([]byte, 2048)

// 	if s.IsLeader {
// 		s.followers[conn] = struct{}{}
// 	}

// 	for {
// 		n, err := conn.Read(buf)
// 		if err != nil {
// 			log.Printf("conn read error %s\n", err)
// 			break
// 		}
// 		msg := buf[:n]
// 		go s.handleCommand(conn, msg)
// 	}
// }

// func (s *Server) handleCommand(conn net.Conn, rawCMD []byte) {
// 	msg, err := parseMessage(rawCMD)
// 	if err != nil {
// 		log.Printf("failed to parse command : %s\n", err)
// 		conn.Write([]byte(err.Error()))
// 		return
// 	}

// 	fmt.Printf("received command %s\n", msg.Cmd)

// 	switch msg.Cmd {
// 	case CMDSet:
// 		err = s.handleSetCmd(conn, msg)
// 	case CMDGet:
// 		err = s.handleGetCmd(conn, msg)
// 	}
// 	if err != nil {
// 		fmt.Printf("failed to handle command: %s\n", err)
// 		conn.Write([]byte(err.Error()))
// 	}

// }

// func (s *Server) handleSetCmd(conn net.Conn, msg *Message) error {
// 	log.Println("Handling the set message: ", msg)
// 	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
// 		return err
// 	}

// 	go s.sendToFollowers(context.TODO(), msg)

// 	return nil
// }

// func (s *Server) handleGetCmd(conn net.Conn, msg *Message) error {
// 	log.Println("Handling the get message: ", msg)
// 	data, err := s.cache.Get(msg.Key)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = conn.Write(data)
// 	return err
// }

// func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
// 	for conn := range s.followers {
// 		fmt.Println("forwarding key to follower")
// 		_, err := conn.Write(msg.ToBytes())
// 		if err != nil {
// 			log.Println("write to follower error: ", err)
// 			continue
// 		}
// 	}
// 	return nil
// }
