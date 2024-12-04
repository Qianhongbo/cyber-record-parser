package record

import (
	"fmt"
	"path/filepath"
)

type Record struct {
	Filename      string
	MajorVersion  uint32
	MinorVersion  uint32
	Size          uint64
	MessageNumber uint64
	StartTime     uint64
	EndTime       uint64
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

	return &Record{
		Filename:      filepath.Base(recordpath),
		MajorVersion:  *header.MajorVersion,
		MinorVersion:  *header.MinorVersion,
		Size:          *header.Size,
		MessageNumber: *header.MessageNumber,
		StartTime:     *header.BeginTime,
		EndTime:       *header.EndTime,
	}
}

// PrintRecord prints the record information
func (r *Record) Print() {
	fmt.Println("Record file path: ", r.Filename)
	fmt.Println("Major version: ", r.MajorVersion)
	fmt.Println("Minor version: ", r.MinorVersion)
	fmt.Println("Size: ", r.Size)
	fmt.Println("Message number: ", r.MessageNumber)
	fmt.Println("Start time: ", r.StartTime)
	fmt.Println("End time: ", r.EndTime)
}
