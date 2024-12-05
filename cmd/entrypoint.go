package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: cyber_record_parser info <record file>")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "info":
        InfoCommand() 
    case "echo":
        EchoCommand()
    default:
        fmt.Println("Unknown command ", os.Args[1])
    }
}