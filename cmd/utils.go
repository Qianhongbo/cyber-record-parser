package main

import (
    "fmt"
    "os"
)

func CheckIfFileExists(filePath string) {
    _, err := os.Stat(filePath)
    if err != nil {
        fmt.Println("Input record file does not exist: ", filePath)
    }
}