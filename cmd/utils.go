package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/eiannone/keyboard"
)

func CheckIfFileExists(filePath string) {
	_, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("Input record file does not exist: ", filePath)
	}
}

var isPaused bool = false
var pauseChan chan bool = make(chan bool)
var stopChan chan bool = make(chan bool)

func listenForSpace() {
	wg.Add(1)
	defer wg.Done()

	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	for {
		select {
		case <-stopChan:
			return
		case event := <-keysEvents:
			if event.Err != nil {
				panic(event.Err)
			}

			if event.Key == keyboard.KeyEsc || event.Key == keyboard.KeyCtrlC || event.Key == keyboard.KeyCtrlD || event.Rune == 'q' {
				fmt.Println("\nExiting program...")
				stopChan <- true
				return
			}

			if event.Rune == ' ' || event.Key == keyboard.KeySpace {
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
