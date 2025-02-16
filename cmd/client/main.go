package main

import (
	"fmt"
	"net"
	"strings"
	"sync"

	h "chat/internal/client/handlers"
	c "chat/internal/config"

	"github.com/eiannone/keyboard"
)

// TODO: ограничение на запуск одного клиента на устройстве
// TODO: сделать возможность получения инфо о чате (например, через команду "?o" выводить инфо об онлайне)
// TODO: БАГ - при нажатии ALT+ЛЮБАЯ_КЛАВИША чат закрывается

func main() {
	var mx sync.Mutex
	var input strings.Builder

	connectString := fmt.Sprintf("%s:%s", c.ServerHost, c.ServerPort)

	conn, err := net.Dial("tcp", connectString)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	serverName := fmt.Sprintf("%s:%s", c.ServerHost, c.ServerPort)
	fmt.Printf("[!] Connected to server %s\n", serverName)

	err = keyboard.Open()
	if err != nil {
		panic(err)
	}

	go h.MessagesHandler(conn, &input, &mx)

	for {
		isExit := h.KeyboardHandler(&input, &mx)

		if isExit {
			// TODO: возможно? добавить disconnect message
			break
		}

		mx.Lock()
		conn.Write([]byte(input.String() + "\n"))
		input.Reset()
		mx.Unlock()
	}

	keyboard.Close()
}
