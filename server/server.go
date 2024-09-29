package main

import (
	"bufio"
	config "chat"
	data "chat"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	clients   = make(map[net.Conn]string)
	broadcast = make(chan []byte)
	colors    = []string{"\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m", "\033[97m"}
	mutex     sync.Mutex
)

func main() {
	port := fmt.Sprintf(":%s", config.ServerPort)
	listener, err := net.Listen("tcp", port)
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

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	randomNumber := r.Intn(8)

	firstMessage := data.ServiceMessage{
		Username:    "",
		Message:     "[S] Enter your nickname: ",
		Color:       "",
		MessageType: "InputUsername",
	}

	jm, _ := json.Marshal(firstMessage)
	conn.Write(append(jm, '\n'))

	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Printf("\n[E] Reading client nickname error: %v", err)
		return
	}

	name = strings.TrimSpace(name)

	mutex.Lock()
	clients[conn] = name
	mutex.Unlock()

	newUserMessage := data.ServiceMessage{
		Username:    name,
		Message:     "",
		Color:       colors[randomNumber],
		MessageType: "NewUser",
	}

	jm, _ = json.Marshal(newUserMessage)

	connectionInfo := fmt.Sprintf("[+] User %s connected", name)
	fmt.Println(connectionInfo)

	broadcast <- append(jm, '\n')

	for {
		rawMessage, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()

			connectionInfo := fmt.Sprintf("[-] User %s disconnected", name)
			fmt.Println(connectionInfo)

			message := data.ServiceMessage{
				Username:    name,
				Message:     "",
				Color:       colors[randomNumber],
				MessageType: "UserDisconnected",
			}

			jm, _ := json.Marshal(message)

			broadcast <- append(jm, '\n')
			break
		}

		message := data.ServiceMessage{
			Username:    name,
			Message:     rawMessage,
			Color:       colors[randomNumber],
			MessageType: "NewMessage",
		}

		jm, _ := json.Marshal(message)

		broadcast <- append([]byte(jm), '\n')
	}
}

func handleBroadcast() {
	for {
		message := <-broadcast
		mutex.Lock()

		for conn := range clients {
			_, err := conn.Write(message)
			if err != nil {
				fmt.Printf("\n[E] Message send error: %v", err)
			}
		}

		mutex.Unlock()
	}
}
