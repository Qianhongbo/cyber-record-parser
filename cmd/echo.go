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
	for msg := range record.ReadMessage() {
		if handleControlSignals() {
			break loop
		}
		if topic != "" && msg.ChannelName != topic {
			continue
		}

		printMessage(record, msg)
	}
}

func printMessage(record *record.Record, message record.Message) {
	clearScreen()

	channelName := message.ChannelName
	fmt.Print(strings.Repeat("-", 50))
	fmt.Println()
	fmt.Printf("Channel name: %s\n", channelName)
	fmt.Printf("Time nanosecond: %d\n", message.Time)
	dt := time.Unix(0, int64(message.Time))
	fmt.Printf("Time: %s\n", dt.Format("2006-01-02 15:04:05"))
	data := message.Content

	// get message type
	if record.Channels[channelName] == nil {
		fmt.Println("Channel not found: ", channelName)
		return
	}

	channelCache := record.Channels[channelName]
	messageTypeStr := channelCache.GetMessageType()
	jsonData, err := record.ConvertMessageToJSON(messageTypeStr, data)
	if err != nil {
		fmt.Println("Failed to marshal message to json: ", err)
		return
	}

	fmt.Println("\nMessage:\n", jsonData)
}