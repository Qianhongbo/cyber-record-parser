package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"cyber_record_parser/internal/record"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

func InfoCommand(cmd *cobra.Command, args []string) {
	theRecordFilePath := args[0]
	CheckIfFileExists(theRecordFilePath)

	record := record.NewRecord(theRecordFilePath)
	defer record.Close()

	printRecordHeaderInfo(record)
}

func printRecordHeaderInfo(record *record.Record) {
	fmt.Println()
	fmt.Println("Cyber Record information:")
	fmt.Println("----------------------------")
	fmt.Println()

	printRecordHeader(record)

	fmt.Println()
	fmt.Println("Channels information:")
	fmt.Println("----------------------------")
	fmt.Println()

	printChannelsInfo(record)
}

func printRecordHeader(record *record.Record) {
	fmt.Printf("- %-20s %s\n", "Record file path", filepath.Base(record.Filepath))

	header := record.Header
	version := fmt.Sprintf("%d.%d", *header.MajorVersion, *header.MinorVersion)
	fmt.Printf("- %-20s %s\n", "Version", version)

	size := humanize.Bytes(*header.Size)
	fmt.Printf("- %-20s %s\n", "Size", size)

	fmt.Printf("- %-20s %s\n", "Compression", header.Compress.String())

	chunkRawSize := humanize.Bytes(*header.ChunkRawSize)
	fmt.Printf("- %-20s %s\n", "Chunk raw size", chunkRawSize)

	chunkInterval := time.Duration(*header.ChunkInterval)
	fmt.Printf("- %-20s %s\n", "Chunk interval", chunkInterval)

	startTime := time.Unix(int64(*header.BeginTime/1e9), 0)
	fmt.Printf("- %-20s %s\n", "Start time", startTime)

	endTime := time.Unix(int64(*header.EndTime/1e9), 0)
	fmt.Printf("- %-20s %s\n", "End time", endTime)

	duration := endTime.Sub(startTime)
	fmt.Printf("- %-20s %s\n", "Duration", duration)

	fmt.Printf("- %-20s %d\n", "Message number", *header.MessageNumber)
	fmt.Printf("- %-20s %d\n", "Channel number", *header.ChannelNumber)
	fmt.Printf("- %-20s %t\n", "Is complete", *header.IsComplete)
}

func printChannelsInfo(record *record.Record) {
	var channelNames []string
	for name := range record.Channels {
		channelNames = append(channelNames, name)
	}

	// Sort channel names
	sort.Strings(channelNames)

	fmt.Printf("%-50s | %-7s | %s\n", "Channel name", "Count", "Type")
	for _, channelName := range channelNames {
		channel := record.Channels[channelName]
		fmt.Printf("%-50s | %-7d | %s\n", *channel.Name, *channel.MessageNumber, *channel.MessageType)
	}
}