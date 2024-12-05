package main

import (
	"fmt"
	"os"

	"cyber_record_parser/internal/record"
)

func InfoCommand() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: cyber_record_parser info <record file>")
        os.Exit(1)
    }

    recordFilePath := os.Args[2]

    // check if the file exists
    _, err := os.Stat(recordFilePath)
    if err != nil {
        fmt.Println("Input record file does not exist: ", recordFilePath)
    }

    record := record.NewRecord(recordFilePath)
    record.Print()
}
