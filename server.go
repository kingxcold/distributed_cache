package main

import (
	"distributed_cache/cache"
	proto "distributed_cache/protocol"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type ServerOpts struct {
	ListenAddr string
	LeaderAddr string
	IsLeader   bool
}

type Server struct {
	ServerOpts
	cache cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{ServerOpts: opts, cache: c}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	log.Printf("server starting on port [%s]\n", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		log.Printf("new connection [%s]\n", conn.RemoteAddr())
		if err != nil {
			log.Printf("accept error %s\n", err)
			continue
		}
		go s.HandleConn(conn)
	}
}

func (s *Server) HandleConn(conn net.Conn) {
	defer conn.Close()
	for {
		cmd, err := proto.ParseCommand(conn)
		if err != nil {
			if err != io.EOF {
				log.Println("parse command error : ", err)
			}
			break
		}
		go s.handleCommand(conn, cmd)
	}
	fmt.Println("connection closed: ", conn.RemoteAddr())
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *proto.CommandSet:
		s.handleSetCommand(conn, v)
	case *proto.CommandGet:
		s.handleGetCommand(conn, v)
	default:
	}
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *proto.CommandSet) error {
	log.Printf("SET %s %s\n", cmd.Key, cmd.Value)
	resp := proto.ResponseSet{}
	if err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL)*time.Second); err != nil {
		resp.Status = proto.StatusERR
		_, err := conn.Write(resp.Bytes())
		return err
	}
	resp.Status = proto.StatusOK
	_, err := conn.Write(resp.Bytes())
	return err
}

func (s *Server) handleGetCommand(conn net.Conn, cmd *proto.CommandGet) error {
	log.Printf("GET %s\n", cmd.Key)
	resp := proto.ResponseGet{}
	val, err := s.cache.Get(cmd.Key)
	if err != nil {
		resp.Status = proto.StatusKeyNotFound
		_, err := conn.Write(resp.Bytes())
		return err
	}
	resp.Status = proto.StatusOK
	resp.Value = val
	_, err = conn.Write(resp.Bytes())
	return err
}

func respondClient(conn net.Conn, msg proto.ResponseGet) error {
	_, err := conn.Write(msg.Bytes())
	if err != nil {
		return err
	}
	return nil
}
