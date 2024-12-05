package main

import (
    "fmt"
    "os"
)

func CheckInputArgs(num int) {
    if len(os.Args) < num {
        fmt.Println("Usage: cyber_record_parser info <record file>")
        os.Exit(1)
    }
}

func CheckIfFileExists(filePath string) {
    _, err := os.Stat(filePath)
    if err != nil {
        fmt.Println("Input record file does not exist: ", filePath)
    }
}