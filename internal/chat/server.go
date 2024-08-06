package chat

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type History struct {
	Messages []string
}

type Server struct {
	port         string
	listener     net.Listener
	clients      map[*Client]bool
	mutex        sync.Mutex
	historyMutex sync.Mutex
	history      History
}

func NewServer(port string) *Server {
	log.Println("Creating new server instance")
	return &Server{
		port:    port,
		clients: make(map[*Client]bool),
		history: History{},
	}
}

func (s *Server) Run() error {
	log.Println("Server Run method started")
	fmt.Println("Listening on the port :", s.port)
	var err error
	s.listener, err = net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to start the server: %v", err)
	}
	defer s.listener.Close()

	log.Println("Server is now listening for connections")

	for {
		log.Println("Waiting for a new connection...")
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		log.Printf("New connection accepted from %s", conn.RemoteAddr())
		go s.HandleConnection(conn)
	}
}

func (s *Server) SendMessageHistory(client *Client) {
	s.historyMutex.Lock()
	defer s.historyMutex.Unlock()
	for _, msg := range s.history.Messages {
		client.Outbox <- msg
	}
}

func (s *Server) BroadcastMessage(message string, sender *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	log.Println(message)
	s.history.Messages = append(s.history.Messages, message)
	for client := range s.clients {
		if client.Conn != sender.Conn {
			client.Conn.Write([]byte(message))
			select {
			case client.Outbox <- message:
			default:
				log.Printf("Failed to send message to client %s: outbox full", client.Name)
			}
		}
	}
}
