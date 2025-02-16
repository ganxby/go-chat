package handlers

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"unicode"

	t "chat/internal/tools"

	"github.com/eiannone/keyboard"
)

func KeyboardHandler(input *strings.Builder, mx *sync.Mutex) bool {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			keyboard.Close()
			log.Fatal(err)
		}

		if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
			return true
		}

		mx.Lock()
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
			mx.Unlock()
			continue
		}
		mx.Unlock()

		if key == 32 {
			mx.Lock()
			input.WriteRune(32)
			mx.Unlock()

			fmt.Print(" ")
			continue
		}

		if key == keyboard.KeyEnter {
			mx.Lock()
			tmp := strings.ReplaceAll(input.String(), " ", "")
			mx.Unlock()

			if len(tmp) == 0 {
				continue
			}

			mx.Lock()
			if input.Len() > 0 {
				mx.Unlock()
				break
			}
			mx.Unlock()
		}

		if t.IsCharInRange(char) {
			fmt.Print(string(char))
			mx.Lock()
			input.WriteRune(char)
			mx.Unlock()
		}
	}

	return false
}
