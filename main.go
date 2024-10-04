package tcp

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	from string
	body []byte
}

type Server struct {
	listener     net.Listener
	listenerAddr string
	quitch       chan struct{}
	messages     chan Message
}

func NewServer(listenerAddr string) *Server {
	return &Server{
		listenerAddr: listenerAddr,
		quitch:       make(chan struct{}),
		messages:     make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.listenerAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	s.listener = listener

	go s.acceptLoop()

	<-s.quitch
	close(s.messages)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Accepted connection from:", conn.RemoteAddr())

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		s.messages <- Message{
			from: conn.RemoteAddr().String(),
			body: buf[:n],
		}
	}
}

func main() {
	server := NewServer(":3000")

	go func() {
		for msg := range server.messages {
			fmt.Printf("(%s): %s\n", msg.from, msg.body)
		}
	}()

	log.Fatal(server.Start())
}
