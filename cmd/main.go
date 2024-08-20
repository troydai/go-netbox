package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig

		log.Println("Received SIGTERM, shutting down gracefully...")
		_ = lis.Close()
		os.Exit(0)
	}()

	// next client ID, the accept is run in a loop therefore no race condition
	clientID := 0

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("failed to accept: %v", err)
		}

		workerID := fmt.Sprintf("worker-%d", clientID)
		clientID++
		log.Printf("Workder %s accepted connection from %s", workerID, conn.RemoteAddr())

		// worker
		go func() {
			for {
				conn.Write([]byte("Hello from " + workerID + "\n"))

				buf := make([]byte, 1024)
				i, err := conn.Read(buf)
				if errors.Is(err, io.EOF) {
					log.Println("Connection closed by client")
					// exit the worker
					return
				} else if err != nil {
					log.Println("Unexpected error:", err)
					// exit the worker
					return
				}

				log.Print("Received:", string(buf[:i]))
			}
		}()
	}
}
