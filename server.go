package main

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
}

type Server struct {
	ServerOpts
	cache cacher.Cache
}

func NewServer(opts ServerOpts, c cacher.Cacher) *Server {
	return &Server{ServerOpts: opts, cache: c}
}
