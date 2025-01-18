package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	c "chat/internal/config"
	dc "chat/internal/dataConv"
)

// TODO: ограничить максимальное количество пользователей
// TODO: ограничить максимальную длину сообщения
// TODO: добавить время сообщений
// TODO: сделать нормальное логгирование на сервере

var (
	clients   = make(map[net.Conn]string)
	broadcast = make(chan []byte)
	colors    = [8]string{"\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m", "\033[97m"}
	mutex     sync.Mutex
)

func main() {
	port := fmt.Sprintf(":%s", c.ServerPort)
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
	userColor := colors[randomNumber]

	firstMessage := dc.ServiceMessage{
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

	message := dc.ServiceMessage{
		Username:    name,
		Message:     "",
		Color:       userColor,
		MessageType: "NewUser",
	}

	jm, err = message.CreateMessage()
	if err != nil {
		fmt.Printf("\n[E] Marshalling error: %v", err)
		return
	}

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

			message := dc.ServiceMessage{
				Username:    name,
				Message:     "",
				Color:       userColor,
				MessageType: "UserDisconnected",
			}

			jm, err = message.CreateMessage()
			if err != nil {
				fmt.Printf("\n[E] Marshalling error: %v", err)
				return
			}

			broadcast <- append(jm, '\n')
			break
		}

		message := dc.ServiceMessage{
			Username:    name,
			Message:     rawMessage,
			Color:       userColor,
			MessageType: "NewMessage",
		}

		jm, err = message.CreateMessage()
		if err != nil {
			fmt.Printf("\n[E] Marshalling error: %v", err)
			return
		}

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
