package record

import (
	"encoding/binary"
	"fmt"
	"os"

	"google.golang.org/protobuf/proto"

	"cyber_record_parser/internal/recordproto"
)

type Section struct {
	type_ recordproto.SectionType
	size  int64
}

type ChannelCache struct {
	Name         string
	MessageType  string
	ProtoDesc    []byte
	MessageCount int
}

type Reader struct {
	file           *os.File
	chunkHeaderIdx []*recordproto.SingleIndex
	chunkBodyIdx   []*recordproto.SingleIndex
	channels       map[string]*recordproto.ChannelCache
}

func NewReader(record string) (*Reader, error) {
	file, err := os.Open(record)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	return &Reader{
		file:     file,
		channels: make(map[string]*recordproto.ChannelCache),
	}, nil
}

func (r *Reader) Close() {
	r.file.Close()
}

func (r *Reader) ReadSection(position int64) (*Section, error) {
	// Seek to the position
	r.file.Seek(position, 0)

	var section Section

	// Read the type file (4 bytes)
	err := binary.Read(r.file, binary.LittleEndian, &section.type_)
	if err != nil {
		return nil, fmt.Errorf("failed to read section type: %v", err)
	}

	// skip the reserved field (4 bytes)
	r.file.Seek(4, 1)

	// Read the size field (8 bytes)
	err = binary.Read(r.file, binary.LittleEndian, &section.size)
	if err != nil {
		return nil, fmt.Errorf("failed to read section size: %v", err)
	}

	return &section, nil
}

func (r *Reader) Read(size int64) ([]byte, error) {
	data := make([]byte, size)
	_, err := r.file.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %v", err)
	}

	return data, nil
}

func (r *Reader) ReadHeader() (*recordproto.Header, error) {
	Section, err := r.ReadSection(0)
	if err != nil {
		return nil, fmt.Errorf("failed to read section: %v", err)
	}

	if Section.type_ != recordproto.SectionType_SECTION_HEADER {
		return nil, fmt.Errorf("invalid section type: %v", Section.type_)
	}

	data, err := r.Read(Section.size)
	if err != nil {
		return nil, fmt.Errorf("failed to read header data: %v", err)
	}

	// put the data into the header
	var header recordproto.Header
	err = proto.Unmarshal(data, &header)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal header: %v", err)
	}

	return &header, nil
}

// read index
func (r *Reader) ReadIndex(position uint64) (*recordproto.Index, error) {
	Section, err := r.ReadSection(int64(position))
	if err != nil {
		return nil, fmt.Errorf("failed to read section: %v", err)
	}

	if Section.type_ != recordproto.SectionType_SECTION_INDEX {
		return nil, fmt.Errorf("invalid section type: %v", Section.type_)
	}

	data, err := r.Read(Section.size)
	if err != nil {
		return nil, fmt.Errorf("failed to read index data: %v", err)
	}

	var index recordproto.Index
	err = proto.Unmarshal(data, &index)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal index: %v", err)
	}

	// iterate over the index entries
	for _, item := range index.Indexes {
		itemType := item.GetType()
		channelCache := item.GetChannelCache()

		switch itemType {
		case recordproto.SectionType_SECTION_CHUNK_HEADER:
			r.chunkHeaderIdx = append(r.chunkHeaderIdx, item)
		case recordproto.SectionType_SECTION_CHUNK_BODY:
			r.chunkBodyIdx = append(r.chunkBodyIdx, item)
		case recordproto.SectionType_SECTION_CHANNEL:
			channelName := channelCache.GetName()
			r.channels[channelName] = channelCache
		}
	}

	return &index, nil
}
