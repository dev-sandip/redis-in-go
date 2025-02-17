package main

import (
	"log/slog"
	"net"
)

const defaultListenAddress = "5001"

type Config struct {
	ListenAddr string
}
type Server struct {
	Config
	ln net.Listener
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddress
	}
	return &Server{
		Config: cfg,
	}
}

func (s *Server) start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		slog.Error("failed to listen", "error", err)
		return err
	}
	s.ln = ln
	return s.acceptLoop()

}
func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("failed to accept", "error", err)
			continue
		}
		go s.handleConn(conn)
	}
}
func (s *Server) handleConn(conn net.Conn) {

}
func main() {

}
