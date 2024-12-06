package main

import (
	"os"

	"cyber_record_parser/internal/record"
)

func InfoCommand() {
	CheckInputArgs(3)

	theRecordFilePath := os.Args[2]
	CheckIfFileExists(theRecordFilePath)

	record := record.NewRecord(theRecordFilePath)
	record.Print()
}
