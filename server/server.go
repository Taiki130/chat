package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type client struct {
	conn net.Conn
	name string
}

var clients []client

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	newClient := client{conn: conn, name: name}
	clients = append(clients, newClient)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		broadcast(message, &newClient)
	}
}

func broadcast(message string, sender *client) {
	for _, client := range clients {
		if client != *sender {
			fmt.Fprintf(client.conn, "%s: %s", sender.name, message)
		}
	}
}
