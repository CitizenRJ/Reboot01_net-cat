package main

import (
	"fmt"
	"log"
	"os"

	"NetP/internal/chat"
)

func main() {
	logFile, err := os.Create("log.txt")
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	port := "8989"
	if len(os.Args) == 2 {
		port = os.Args[1]
	} else if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	s := chat.NewServer(port)
	log.Fatal(s.Run())
}
