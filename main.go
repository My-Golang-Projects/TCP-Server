package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	// an empty struct won't take memory
	quitch chan struct{}
	// byte can be anything, from protobuf to string
	msgch chan Message
}

// a constructor
func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		// buffer the channel
		msgch: make(chan Message, 10),
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

	// when the server closes, notify everyone that the channel is also closed, ciao
	close(s.msgch)
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
		fmt.Println("New Connection to the server: ", conn.RemoteAddr())
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

		// when someone sends us something we write that to channel so we can read whenever
		// we want
		s.msgch <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}
		// msg := string(buf[:n])
		// fmt.Println(msg)

		conn.Write([]byte("Thank You For Your Message!\n"))

	}
}

func main() {
	server := NewServer(":3000")
	go func() {
		// receive value from channel until it is closed (with s.msgch.Close()
		for msg := range server.msgch {
			fmt.Printf("Received message from connection (%s):(%s): \n", msg.from, msg.payload)
		}
	}()
	log.Fatal(server.Start())
}
