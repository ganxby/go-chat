package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	dc "chat/internal/dataConv"

	"github.com/eiannone/keyboard"
)

func MessagesHandler(conn net.Conn, input *strings.Builder, mx *sync.Mutex) {
	var message string
	messagesCounter := 0

	for {
		rawMessage, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				fmt.Printf("\n[E] Something went wrong: %v\n", err)
			}

			keyboard.Close()
			os.Exit(0)
		}

		var jm dc.ServiceMessage
		err = json.Unmarshal([]byte(rawMessage), &jm)
		if err != nil {
			fmt.Printf("\n[E] Unmarshalling error: %v\n", err)
			continue
		}

		// Очистка строки и перенос курсора в начало
		if messagesCounter > 1 {
			fmt.Print("\033[2K")
			fmt.Print("\033[0G")
		}

		if jm.MessageType == "InputUsername" {
			message = jm.Message
		} else if jm.MessageType == "NewUser" {
			message = fmt.Sprintf("[+] User %s connected\n", jm.Username)
		} else if jm.MessageType == "UserDisconnected" {
			message = fmt.Sprintf("[-] User %s disconnected\n", jm.Username)
		} else {
			message = fmt.Sprintf("%s%s\033[0m: %s", jm.Color, jm.Username, jm.Message)
		}

		// Если сообщение второе, то добавить перенос строки
		if messagesCounter == 1 {
			fmt.Print("\n" + message)
		} else {
			fmt.Print(message)
		}

		messagesCounter += 1

		if messagesCounter > 1 {
			mx.Lock()
			fmt.Print("--> " + input.String())
			mx.Unlock()
		}
	}
}
