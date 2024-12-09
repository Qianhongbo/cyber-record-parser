package main

import (
	"os"

	"cyber_record_parser/internal/record"

	"github.com/spf13/cobra"
)

func EchoCommand(cmd *cobra.Command, args []string) {
	CheckInputArgs(3)
    

	theRecordFilePath := os.Args[2]
	CheckIfFileExists(theRecordFilePath)

	record := record.NewRecord(theRecordFilePath)

    topic := cmd.Flag("topic").Value.String()
	record.PrintTopicMsg(topic)
}
