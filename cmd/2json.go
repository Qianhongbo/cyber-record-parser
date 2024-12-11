package main

import (
	"cyber_record_parser/internal/record"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var isStart bool = true

func ToJsonCommand(cmd *cobra.Command, args []string) {
	theRecordFilePath := args[0]
	CheckIfFileExists(theRecordFilePath)

	record := record.NewRecord(theRecordFilePath)
	defer record.Close()

	topic := cmd.Flag("topic").Value.String()
	output := cmd.Flag("output").Value.String()
	saveTopicMsgToJson(record, topic, output)
	fmt.Printf("Save topic (%s) messages to %s \n", topic, output)
}

// Write the opening bracket for the JSON array
func writeJsonStart(file *os.File) error {
	_, err := file.WriteString("[\n")
	if err != nil {
		return fmt.Errorf("failed to write opening bracket to json file: %w", err)
	}
	return nil
}

// Write the closing bracket for the JSON array
func writeJsonEnd(file *os.File) error {
	_, err := file.WriteString("\n]")
	if err != nil {
		return fmt.Errorf("failed to write closing bracket to json file: %w", err)
	}
	return nil
}

func saveTopicMsgToJson(record *record.Record, topic string, output string) (path string, err error) {
	// Open the file in write mode and truncate its content
	file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open json file for writing: %w", err)
	}
	defer file.Close()

	err = writeJsonStart(file)
	if err != nil {
		return "", fmt.Errorf("failed to write opening bracket to json file: %w", err)
	}

	// Loop through the messages and write them to the file
	for msg := range record.ReadMessage() {
		if topic != "" && msg.ChannelName != topic {
			continue
		}

		channelName := msg.ChannelName
		ts := msg.Time
		data := msg.Content

		// Skip if channel is not available
		if record.Channels[channelName] == nil {
			continue
		}

		channelCache := record.Channels[channelName]
		messageTypeStr := channelCache.GetMessageType()

		// Convert the message content to JSON
		contentData, err := record.ConvertMessageToJSON(messageTypeStr, data, nil)
		if err != nil {
			return "", fmt.Errorf("failed to convert message to json: %w", err)
		}
		var raw json.RawMessage = []byte(contentData)

		// Construct the JSON data
		jsonData := map[string]interface{}{
			"topic": channelName,
			"time":  ts,
			"data":  raw,
		}

		// Marshal the data into JSON format
		marshaledData, err := json.Marshal(jsonData)
		if err != nil {
			return "", fmt.Errorf("failed to marshal json data: %w", err)
		}

		// Write the JSON data to the file
		if err := writeMsgToFile(file, marshaledData); err != nil {
			return "", err
		}
	}

	err = writeJsonEnd(file)
	if err != nil {
		return "", fmt.Errorf("failed to write closing bracket to json file: %w", err)
	}

	return output, nil
}

// Write the JSON data to the file, ensuring proper formatting
func writeMsgToFile(file *os.File, jsonData []byte) error {
	var err error
	// If this is the first message, no need to add a comma
	if isStart {
		_, err = file.Write(jsonData)
		isStart = false
	} else {
		// Write a comma before appending the new message
		_, err = file.WriteString(",\n")
		if err == nil {
			_, err = file.Write(jsonData)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to write json data to file: %w", err)
	}
	return nil
}
