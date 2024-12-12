package main

import (
	"fmt"
	"strings"
	"time"

	"cyber_record_parser/internal/record"

	"github.com/spf13/cobra"
)

func EchoCommand(cmd *cobra.Command, args []string) {
	theRecordFilePath := args[0]
	CheckIfFileExists(theRecordFilePath)

	record := record.NewRecord(theRecordFilePath)
	defer record.Close()

    topic := cmd.Flag("topic").Value.String()
	printTopicMsg(record, topic)
}

func printTopicMsg(record *record.Record, topic string) {
	go listenForSpace()

loop:
	for msg := range record.ReadMessages() {
		if handleControlSignals() {
			break loop
		}
		if topic != "" && msg.ChannelName != topic {
			continue
		}

		printMessage(msg)
	}
}

func printMessage(message record.Message) {
	clearScreen()

	channelName := message.ChannelName
	fmt.Print(strings.Repeat("-", 50))
	fmt.Println()
	fmt.Printf("Channel name: %s\n", channelName)
	fmt.Printf("Time nanosecond: %d\n", message.Time)
	dt := time.Unix(0, int64(message.Time))
	fmt.Printf("Time: %s\n", dt.Format("2006-01-02 15:04:05"))
	data := message.Content
	fmt.Println("\nMessage:\n", string(data))
}