package record

import (
	"fmt"
	"path/filepath"
	"time"
	"sort"

	"github.com/dustin/go-humanize"

	"cyber_record_parser/internal/recordproto"
)

type Record struct {
	Filename string
	Header   *recordproto.Header
	Index    *recordproto.Index
	Channels map[string]*recordproto.ChannelCache
}

// NewRecord creates a new record struct
func NewRecord(recordpath string) *Record {
	recordReader, err := NewReader(recordpath)
	if err != nil {
		fmt.Println("Failed to create record reader: ", err)
	}

	defer recordReader.Close()

	header, err := recordReader.ReadHeader()
	if err != nil {
		fmt.Println("Failed to read header: ", err)
	}

	index, err := recordReader.ReadIndex(*header.IndexPosition)
	if err != nil {
		fmt.Println("Failed to read index: ", err)
	}

	// get channel info map
	channels := recordReader.channels

	return &Record{
		Filename: filepath.Base(recordpath),
		Header:   header,
		Index:    index,
		Channels: channels,
	}
}

// PrintRecord prints the record information
func (r *Record) Print() {
	fmt.Println()
	fmt.Println("Cyber Record information:")
	fmt.Println("----------------------------")
	fmt.Println()

	fmt.Printf("- %-20s %s\n", "Record file path", r.Filename)

	header := r.Header
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

	fmt.Println()
	fmt.Println("Channels information:")
	fmt.Println("----------------------------")
	fmt.Println()

	// sort channels by name alphabetically
	var channelNames []string
	for name := range r.Channels {
		channelNames = append(channelNames, name)
	}
	// sort channel names
	sort.Strings(channelNames)

	fmt.Printf("%-50s | %-7s | %s\n", "Channel name", "Count", "Type")
	for _, channelName := range channelNames {
		channel := r.Channels[channelName]
		fmt.Printf("%-50s | %-7d | %s\n", *channel.Name, *channel.MessageNumber, *channel.MessageType)
	}

}
