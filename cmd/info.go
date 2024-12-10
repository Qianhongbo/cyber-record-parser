package main

import (
	"cyber_record_parser/internal/record"

	"github.com/spf13/cobra"
)

func InfoCommand(cmd *cobra.Command, args []string) {
	theRecordFilePath := args[0]
	CheckIfFileExists(theRecordFilePath)

	record := record.NewRecord(theRecordFilePath)
	defer record.Close()

	record.PrintRecordHeaderInfo()
}
