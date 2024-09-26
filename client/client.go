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

// TODO: исправить удаление второй строки при вводе никнейма
// TODO: вынести обработку клавиатуры в отдельную функцию
// TODO: попробовать сделать input в main и работать с ним через указатели

var input strings.Builder

func main() {
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

	go messagesHandler(conn)

	for {
		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				log.Fatal(err)
			}

			if key == keyboard.KeyEsc {
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
				break
			}

			fmt.Print(string(char))
			input.WriteRune(char)
		}

		// message, _ := input.ReadString('\n')
		conn.Write([]byte(input.String() + "\n"))
		input.Reset()
	}
}

func messagesHandler(conn net.Conn) {
	messagesCounter := 0

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("[E] Something went wrong: %v\n", err)
			os.Exit(0)
		}

		if messagesCounter == 0 {
			message = message[:len(message)-1]
		}

		fmt.Print("\033[2K")
		fmt.Print("\033[0G")

		fmt.Print(message)
		messagesCounter += 1

		if messagesCounter > 1 {
			fmt.Print("--> " + input.String())
		}
	}
}
