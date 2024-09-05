package chat

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

func (s *Server) HandleIncomingMessages(client *Client) {
	defer func() {
		s.mutex.Lock()
		delete(s.clients, client)
		s.mutex.Unlock()
		s.BroadcastMessage(fmt.Sprintf("%s has left our chat...\n", client.Name), nil)
		client.Conn.Close()
	}()

	reader := bufio.NewReader(client.Conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("Failed to read message from client %s: %v", client.Name, err)
			return
		}

		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}
		if msg == "/change" {
			s.mutex.Lock()
			delete(s.clients, client)
			s.mutex.Unlock()
			s.HandleConnection(client.Conn)	
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		formattedMsg := fmt.Sprintf("[%s][%s]:%s\n", timestamp, client.Name, msg)

		// Send to the client's own connection
		client.Conn.Write([]byte(formattedMsg))

		// Broadcast to others
		s.BroadcastMessage(formattedMsg, client)
	}
}

func (s *Server) HandleOutgoingMessages(client *Client) {
	for msg := range client.Outbox {
		_, err := client.Conn.Write([]byte(msg))
		if err != nil {
			log.Printf("Failed to send message to client %s: %v", client.Name, err)
			return
		}
	}
}
