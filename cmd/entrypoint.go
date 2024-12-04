package main

import (
    "fmt"
    "os"

    "cyber_record_parser/cmd/cyber_record_info"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: cyber_record_parser info <record file>")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "info":
        cyber_record_info.InfoCommand() 
    default:
        fmt.Println("Unknown command ", os.Args[1])
    }
}