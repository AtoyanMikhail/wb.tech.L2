package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		log.Fatal("Usage: go-telnet [--timeout=10s] host port")
	}

	host := args[0]
	port := args[1]
	address := net.JoinHostPort(host, port)

	conn, err := net.DialTimeout("tcp", address, *timeout)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", address)

	done := make(chan struct{})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// reader 
	go func() {
		reader := bufio.NewReader(conn)
		for {
			select {
			case <-done:
				return
			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					return
				}
				fmt.Print(line)
			}
		}
	}()

	// writer
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			select {
			case <-done:
				return
			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						fmt.Println("^D")
						close(done)
						return
					}
					log.Printf("Stdin error: %v", err)
					continue
				}
				_, err = conn.Write([]byte(line))
				if err != nil {
					log.Printf("Write error: %v", err)
					close(done)
					return
				}
			}
		}
	}()

	select {
	case <-sigCh:
		fmt.Println("\nConnection closed")
	case <-done:
		fmt.Println("Connection closed")
	}
}