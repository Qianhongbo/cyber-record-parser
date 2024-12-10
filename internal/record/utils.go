package record

import (
	"fmt"
	"runtime"
	"time"

	"github.com/eiannone/keyboard"
)

var isPaused bool = false
var pauseChan chan bool = make(chan bool)
var stopChan chan bool = make(chan bool)

func listenForSpace() {
	err := keyboard.Open()
	if err != nil {
		fmt.Println("Error opening keyboard:", err)
		return
	}
	defer keyboard.Close()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println("Error reading key:", err)
			return
		}

		if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC || key == keyboard.KeyCtrlD || char == 'q' {
			fmt.Println("\nExiting program...")
			stopChan <- true
			return
		}

		if char == ' ' || key == keyboard.KeySpace {
			isPaused = !isPaused
			pauseChan <- isPaused
			if isPaused {
				fmt.Println("\nPaused. Press SPACE to resume or ESC / q / Ctrl+C to exit...")
			} else {
				fmt.Println("\nResumed.")
			}
		}
	}
}

func handleControlSignals() bool {
	select {
	case <-stopChan:
		return true
	case isPaused := <-pauseChan:
		for isPaused {
			time.Sleep(100 * time.Millisecond)
			select {
			case <-stopChan:
				return true
			case isPaused = <-pauseChan:
			default:
			}
		}
	default:
	}
	return false
}

func clearScreen() {
	switch os := runtime.GOOS; os {
	case "windows":
		fmt.Print("\033[H\033[2J") // Windows清屏
	default:
		fmt.Print("\033[H\033[2J") // Linux/macOS清屏
	}
}
