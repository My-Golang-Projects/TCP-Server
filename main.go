package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	// an empty struct won't take memory
	quitch chan struct{}
}

// a constructor
func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	// close the listener if the listener is opened
	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()

	// wait for the quitch channel
	// if the quitch channel is closed we can defer the listener
	<-s.quitch
	return nil
}

func (s *Server) acceptLoop() {
	for {
		// accept a connection
		// Accept waits for and returns the next connection to the listener.
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Accept error: ", err)
			continue
		}
		fmt.Println("New Connection to the server: ", conn.RemoteAddr().String())
		// each time we accept, we spin up a new goroutine so it's not blocking
		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		// read into the buffer
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read Error: ", err)
			// there might be malformed msg, we could drop here also by returning error
			continue
		}

		msg := string(buf[:n])
		fmt.Println(msg)

	}
}

func main() {
	server := NewServer(":3000")
	log.Fatal(server.Start())
}
