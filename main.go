package main

import (
	"log"
	"net"

	"golang.org/x/exp/slog"
)

const DEFAULT_LISTEN_ADDRESS string = ":5001"

type Config struct {
	ListenAddress string
}

type Server struct {
	Config
	ln         net.Listener
	addPeerCh  chan *Peer
	peers      map[*Peer]bool
	quitPeerCh chan struct{}
}

func NewServer(cfg Config) *Server {
	if cfg.ListenAddress == "" {
		cfg.ListenAddress = DEFAULT_LISTEN_ADDRESS
	}

	return &Server{
		Config:     cfg,
		addPeerCh:  make(chan *Peer),
		peers:      make(map[*Peer]bool),
		quitPeerCh: make(chan struct{}),
	}
}

func (srv *Server) Start() error {
	ln, err := net.Listen("tcp", srv.ListenAddress)
	if err != nil {
		return err
	}
	srv.ln = ln

	go srv.loop()
	return srv.AcceptLoop()
}

func (srv *Server) loop() {
	for {
		select {
		case <-srv.quitPeerCh:
			return
		case peer := <-srv.addPeerCh:
			srv.peers[peer] = true
		}
	}
}

func (srv *Server) AcceptLoop() error {
	for {
		conn, err := srv.ln.Accept()
		if err != nil {
			slog.Error("accept error", "err", err)
			continue
		}
		slog.Info("new connection accepted", "remoteAddr", conn.RemoteAddr().String())
		go srv.handleConnection(conn)
	}
}

func (srv *Server) handleConnection(conn net.Conn) {
	slog.Info("new peer called", "remoteAddr", conn.RemoteAddr().String())
	peer := NewPeer(conn)
	srv.addPeerCh <- peer
	slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr().String())
	if err := peer.readLoop(); err != nil {
		slog.Error("error read connection", "error", err)
	}
}

func main() {
	server := NewServer(Config{})
	log.Fatal(server.Start())
}