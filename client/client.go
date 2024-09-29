package main

import (
	"bufio"
	config "chat"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"unicode"

	"github.com/eiannone/keyboard"
)

// TODO: ограничение на запуск одного клиента на устройстве

func main() {
	var input strings.Builder

	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	connectString := fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort)

	conn, err := net.Dial("tcp", connectString)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	serverName := fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort)
	fmt.Printf("[!] Connected to server %s\n", serverName)

	go messagesHandler(conn, &input)

	for {
		keyboardHandler(&input)

		conn.Write([]byte(input.String() + "\n"))
		input.Reset()
	}
}

func isCharInRange(r rune) bool {
	if unicode.Is(unicode.Latin, r) {
		return true
	}

	if unicode.Is(unicode.Cyrillic, r) {
		return true
	}

	if unicode.IsDigit(r) {
		return true
	}

	switch r {
	case ',', '*', '!', ')', '(', '#', '@', '$', '%', '&', '<', '>', '"', '\'', '\\', '|', '/', ':', ';', '^', '?', '-', '_', '=', '+', '.':
		return true
	}

	return false
}

func keyboardHandler(input *strings.Builder) {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
			keyboard.Close()
			os.Exit(1)
		}

		if key == 127 {
			if input.Len() > 0 {
				inputStr := input.String()
				input.Reset()
				input.WriteString(inputStr[:len(inputStr)-1])

				fmt.Print("\b \b")
			}
			continue
		}

		if key == 32 {
			input.WriteRune(32)
			fmt.Print(" ")
			continue
		}

		if key == keyboard.KeyEnter {
			tmp := strings.ReplaceAll(input.String(), " ", "")

			if len(tmp) == 0 {
				continue
			}

			if input.Len() > 0 {
				break
			}
		}

		if isCharInRange(char) {
			fmt.Print(string(char))
			input.WriteRune(char)
		}
	}
}

func messagesHandler(conn net.Conn, input *strings.Builder) {
	messagesCounter := 0

	for {
		rawMessage, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("[E] Something went wrong: %v\n", err)
			keyboard.Close()
			os.Exit(0)
		}

		var jm config.ServiceMessage
		err = json.Unmarshal([]byte(rawMessage), &jm)
		if err != nil {
			fmt.Printf("[E] JSON decoding error: %v\n", err)
			continue
		}

		// Очистка строки и перенос курсора в начало
		if messagesCounter > 1 {
			fmt.Print("\033[2K")
			fmt.Print("\033[0G")
		}

		var message string

		if jm.MessageType == "InputUsername" {
			message = jm.Message
		} else if jm.MessageType == "NewUser" {
			message = fmt.Sprintf("[+] User %s connected\n", jm.Username)
		} else if jm.MessageType == "UserDisconnected" {
			message = fmt.Sprintf("[-] User %s disconnected\n", jm.Username)
		} else {
			message = fmt.Sprintf("%s%s\033[0m: %s", jm.Color, jm.Username, jm.Message)
		}

		// Если сообщение второе (), то добавить перенос строки
		if messagesCounter == 1 {
			fmt.Print("\n" + message)
		} else {
			fmt.Print(message)
		}

		messagesCounter += 1

		if messagesCounter > 1 {
			fmt.Print("--> " + input.String())
		}
	}
}
