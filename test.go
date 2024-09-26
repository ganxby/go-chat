package main

import (
	"fmt"
	"log"

	"github.com/eiannone/keyboard"
)

func main() {
	// Инициализация клавиатурного ввода
	err := keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close() // Закрываем клавиатурный ввод при выходе

	fmt.Println("Нажмите клавиши (нажмите ESC для выхода):")
	var message string
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if key == keyboard.KeyEsc {
			break
		}

		if key == keyboard.KeyEnter {
			break
		}

		fmt.Print(string(char))
		message += string(char)
	}

	fmt.Println("\n" + message)
}
