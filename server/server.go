/*
A very simple TCP server written in Go.

This is a toy project that I used to learn the fundamentals of writing
Go code and doing some really basic network stuff.

Maybe it will be fun for you to read. It's not meant to be
particularly idiomatic, or well-written for that matter.
*/
package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var addr = flag.String("addr", "", "The address to listen to; default is \"\" (all interfaces).")
var port = flag.Int("port", 8000, "The port to listen on; default is 8000.")
var clientIndex int = 0

func main() {
	flag.Parse()

	log.Println("Starting server...")

	src := *addr + ":" + strconv.Itoa(*port)
	listener, _ := net.Listen("tcp", src)
	log.Printf("Listening on %s.\n", src)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Some connection error: %s\n", err)
			continue
		}

		clientIndex++
		log.Printf("server info : client index is %v.", clientIndex)
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	log.Println("Client connected from " + remoteAddr)

	scanner := bufio.NewScanner(conn)

	for {
		ok := scanner.Scan()

		handleMessage(scanner.Text(), conn)

		if !ok {
			break
		}
	}

	log.Println("Client at " + remoteAddr + " disconnected.")
}

func handleMessage(message string, conn net.Conn) {
	log.Println(conn.RemoteAddr().String() + "> " + message)

	if len(message) > 0 && message[0] == '/' {
		switch {
		case message == "/time":
			resp := "It is " + time.Now().String() + "\n"
			log.Print(conn.RemoteAddr().String() + "< " + resp)
			conn.Write([]byte(resp))

		case message == "/quit":
			log.Println("Quitting.")
			conn.Write([]byte("I'm shutting down now.\n"))
			log.Println(conn.RemoteAddr().String() + "< " + "%quit%")
			conn.Write([]byte("%quit%\n"))
			os.Exit(0)

		default:
			conn.Write([]byte("Unrecognized command.\n"))
		}
	}
}
