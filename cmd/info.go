package main

import (
	"os"

	"cyber_record_parser/internal/record"

	"github.com/spf13/cobra"
)

func InfoCommand(cmd *cobra.Command, args []string) {
	CheckInputArgs(3)

	theRecordFilePath := os.Args[2]
	CheckIfFileExists(theRecordFilePath)

	record := record.NewRecord(theRecordFilePath)
	record.PrintRecordHeaderInfo()
}
