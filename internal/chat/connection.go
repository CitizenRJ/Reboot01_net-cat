package chat

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("New connection established from %s", conn.RemoteAddr())

	// set a maximum of 10 connections only
	if len(s.clients) >= 10 {
		_, err := conn.Write([]byte("Sorry, the chat server is at maximum capacity. Please try again later.\n"))
		if err != nil {
			log.Printf("Error sending 'at capacity' message: %v", err)
		}
		return
	}

	var name string
	reader := bufio.NewReader(conn)
	welcomeMsg := GetWelcomeMessage()
	_, err := conn.Write([]byte(welcomeMsg))
	if err != nil {
		log.Printf("Error sending welcome message: %v", err)
		return
	}
	for {
		name, err = reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading client name: %v", err)
			return
		}
		name = strings.TrimSpace(name)

		if name != "" {
			break
		}

		_, err = conn.Write([]byte("Name cannot be empty. Please try again.\n"))
		if err != nil {
			log.Printf("Error sending empty name message: %v", err)
			return
		}
	}

	client := &Client{
		Name:   name,
		Conn:   conn,
		Outbox: make(chan string, 1000),
	}

	s.mutex.Lock()
	s.clients[client] = true
	joinMsg := fmt.Sprintf("%s has joined our chat...\n", client.Name)
	for c := range s.clients {
		c.Outbox <- joinMsg
	}
	s.mutex.Unlock()

	s.SendMessageHistory(client)

	go s.HandleIncomingMessages(client)
	s.HandleOutgoingMessages(client)
}
