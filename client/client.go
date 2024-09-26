package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
)

// TODO: обработка пустого ввода
// TODO: попробовать сделать input в main и работать с ним через указатели

func main() {
	var input strings.Builder

	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Printf("[!] Connected to server localhost:8080\n")

	go messagesHandler(conn, &input)

	for {
		keyboardHandler(&input)

		conn.Write([]byte(input.String() + "\n"))
		input.Reset()
	}
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

		if key == keyboard.KeyEnter {
			if input.Len() > 0 {
				break
			}
		}

		fmt.Print(string(char))
		input.WriteRune(char)
	}
}

func messagesHandler(conn net.Conn, input *strings.Builder) {
	messagesCounter := 0

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("[E] Something went wrong: %v\n", err)
			keyboard.Close()
			os.Exit(0)
		}

		if messagesCounter == 0 {
			message = message[:len(message)-1]
		}

		if messagesCounter > 1 {
			fmt.Print("\033[2K")
			fmt.Print("\033[0G")
		}

		/*
			Разделение введенного никнейма и первого сообщения от сервера;
			Далее необходимо убрать это путем введения спецсимволов от сервера и их обработки
		*/
		if messagesCounter == 1 {
			fmt.Print("\n")
		}

		fmt.Print(message)

		messagesCounter += 1

		if messagesCounter > 1 {
			fmt.Print("--> " + input.String())
		}
	}
}
