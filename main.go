package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ahmadfarhanstwn/goredis/client"
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
	msgCh      chan []byte
	keyVal     *KeyVal
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
		msgCh:      make(chan []byte),
		keyVal:     NewKeyVal(),
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

func (srv *Server) handleRawMessage(rawMsg []byte) error {
	cmd, err := ParseCommand(string(rawMsg))
	if err != nil {
		return err
	}

	switch v := cmd.(type) {
	case SetCommand:
		return srv.keyVal.Set(v.key, v.value)
	}

	return nil
}

func (srv *Server) loop() {
	for {
		select {
		case rawMsg := <-srv.msgCh:
			if err := srv.handleRawMessage(rawMsg); err != nil {
				slog.Error("raw message error", "error", err)
			}
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
	peer := NewPeer(conn, srv.msgCh)
	srv.addPeerCh <- peer
	slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr().String())
	if err := peer.readLoop(); err != nil {
		slog.Error("error read connection", "error", err)
	}
}

func main() {
	server := NewServer(Config{})
	go func() {
		log.Fatal(server.Start())
	}()

	time.Sleep(1 * time.Second)

	for i := 0; i < 10; i++ {
		client := client.New("localhost:5001")
		if err := client.Set(context.TODO(), fmt.Sprintf("foo%d", i), fmt.Sprintf("bar%d", i)); err != nil {
			log.Fatal(err)
		}
	}

	time.Sleep(1 * time.Second)
	fmt.Println(server.keyVal.data)
}
