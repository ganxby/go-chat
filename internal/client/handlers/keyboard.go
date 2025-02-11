package handlers

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	t "chat/internal/tools"

	"github.com/eiannone/keyboard"
)

func KeyboardHandler(input *strings.Builder) bool {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			keyboard.Close()
			log.Fatal(err)
		}

		if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
			return true
		}

		if key == 127 && input.Len() > 0 {
			inputStr := input.String()
			input.Reset()

			runeStr := []rune(inputStr)
			if unicode.Is(unicode.Cyrillic, runeStr[len(runeStr)-1]) {
				input.WriteString(inputStr[:len(inputStr)-2])
			} else {
				input.WriteString(inputStr[:len(inputStr)-1])
			}

			fmt.Print("\b \b")
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

		if t.IsCharInRange(char) {
			fmt.Print(string(char))
			input.WriteRune(char)
		}
	}

	return false
}
