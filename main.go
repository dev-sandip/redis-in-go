package main

import (
	"log"
	"log/slog"
	"net"
)

const defaultListenAddress = "5001"

type Config struct {
	ListenAddr string
}
type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	quitCh    chan struct{}
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddress
	}
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
	}
}
func (s *Server) loop() {
	for {
		select {
		case <-s.quitCh:
			return
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		}
	}
}
func (s *Server) start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		slog.Error("failed to listen", "error", err)
		return err
	}
	s.ln = ln
	go s.loop()
	slog.Info("server started", "address", s.ListenAddr)
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
	peer := NewPeer(conn)
	s.addPeerCh <- peer
	slog.Info("new peer", "address", peer.conn.RemoteAddr())
	if err := peer.readLoop(); err != nil {
		slog.Info("peer disconnected", "address", peer.conn.RemoteAddr())
	}
}
func main() {
	server := NewServer(Config{ListenAddr: ":5001"})
	log.Fatal(server.start())
}
