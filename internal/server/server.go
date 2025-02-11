package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	dc "chat/internal/dataConv"
)

var (
	colors = [8]string{"\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m", "\033[97m"}
)

type Server struct {
	Addr      *net.TCPAddr
	clients   map[net.Conn]string
	broadcast chan []byte
	mx        sync.Mutex
}

func NewServer(addr net.TCPAddr) *Server {
	return &Server{
		Addr:      &addr,
		clients:   make(map[net.Conn]string),
		broadcast: make(chan []byte),

		// Later: logger and colors
	}
}

func (s *Server) Serve() {
	listener, err := net.ListenTCP("tcp", s.Addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	go s.handleBroadcast()

	log.Println("[!] Server has been started")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("\n[E] Client connection error: %v", err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
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

	jm, err := firstMessage.CreateMessage()
	if err != nil {
		log.Printf("\n[E] Marshalling error: %v", err)
		return
	}

	conn.Write(jm)

	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return
		}

		log.Printf("[E] Reading client nickname error: %v", err)
		return
	}

	name = strings.TrimSpace(name)

	s.mx.Lock()
	s.clients[conn] = name
	s.mx.Unlock()

	message := dc.ServiceMessage{
		Username:    name,
		Message:     "",
		Color:       userColor,
		MessageType: "NewUser",
	}

	jm, err = message.CreateMessage()
	if err != nil {
		log.Printf("\n[E] Marshalling error: %v", err)
		return
	}

	connectionInfo := fmt.Sprintf("[+] User '%s' connected", name)
	log.Println(connectionInfo)

	s.broadcast <- jm

	for {
		rawMessage, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			s.mx.Lock()
			delete(s.clients, conn)
			s.mx.Unlock()

			connectionInfo := fmt.Sprintf("[-] User '%s' disconnected", name)
			log.Println(connectionInfo)

			message := dc.ServiceMessage{
				Username:    name,
				Message:     "",
				Color:       userColor,
				MessageType: "UserDisconnected",
			}

			jm, err = message.CreateMessage()
			if err != nil {
				log.Printf("\n[E] Marshalling error: %v", err)
				return
			}

			s.broadcast <- jm
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
			log.Printf("\n[E] Marshalling error: %v", err)
			return
		}

		s.broadcast <- jm
	}
}

func (s *Server) handleBroadcast() {
	for {
		message := <-s.broadcast
		s.mx.Lock()

		for conn := range s.clients {
			_, err := conn.Write(message)
			if err != nil {
				log.Printf("\n[E] Message sending error: %v", err)
			}
		}

		s.mx.Unlock()
	}
}
