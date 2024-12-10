package main

import (
	"cyber_record_parser/internal/record"

	"github.com/spf13/cobra"
)

func EchoCommand(cmd *cobra.Command, args []string) {
	theRecordFilePath := args[0]
	CheckIfFileExists(theRecordFilePath)

	record := record.NewRecord(theRecordFilePath)
	defer record.Close()

    topic := cmd.Flag("topic").Value.String()
	record.PrintTopicMsg(topic)
}
