package main

import (
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

	conn, err := lis.Accept()
	if err != nil {
		log.Fatalf("failed to accept: %v", err)
	}

	log.Println("Accepted connection from", conn.RemoteAddr())
}
