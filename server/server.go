package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	clients   = make(map[net.Conn]string)
	broadcast = make(chan string)
	mutex     sync.Mutex
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	go handleBroadcast()

	fmt.Println("[!] Server has been started")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("\n[E] Client connection error: %v", err)
			continue
		}

		go hahdleClientConn(conn)
	}
}

func hahdleClientConn(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("[S] Enter your nickname: \n"))
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Printf("\n[E] Reading client nickname error: %v", err)
		return
	}

	name = strings.TrimSpace(name)

	mutex.Lock()
	clients[conn] = name
	mutex.Unlock()

	broadcast <- fmt.Sprintf("[+] User %s connected", name)

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			broadcast <- fmt.Sprintf("[-] User %s disconnected", name)
			break
		}

		broadcast <- fmt.Sprintf("\033[32m%s\033[0m: %s", name, strings.TrimSpace(message))
	}
}

func handleBroadcast() {
	for {
		message := <-broadcast
		mutex.Lock()

		for conn := range clients {
			_, err := conn.Write([]byte(message + "\n"))
			if err != nil {
				fmt.Printf("\n[E] Message send error: %v", err)
			}
		}

		mutex.Unlock()
	}
}
