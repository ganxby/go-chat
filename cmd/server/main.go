package main

import (
	"fmt"
	"log"
	"net"

	c "chat/internal/config"
	s "chat/internal/server"
)

// TODO: ограничить максимальное количество пользователей
// TODO: ограничить максимальную длину сообщения
// TODO: добавить время сообщений
// TODO: сделать нормальное логгирование на сервере
// TODO: handleConn разбить на простые функции

func main() {
	strAddr := fmt.Sprintf(":%s", c.ServerPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp", strAddr)
	if err != nil {
		log.Fatal(err)
	}

	server := s.NewServer(*tcpAddr)
	server.Serve()
}
