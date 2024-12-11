package record

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"

	"cyber_record_parser/internal/recordproto"
)

type Record struct {
	Filepath       string
	Reader         *Reader
	Header         *recordproto.Header
	ChunkHeaderIdx []*recordproto.SingleIndex
	ChunkBodyIdx   []*recordproto.SingleIndex
	Channels       map[string]*recordproto.ChannelCache
}

type Message struct {
	ChannelName string
	Time        uint64
	Content     []byte
}

func NewRecord(recordpath string) *Record {
	recordReader, err := NewReader(recordpath)
	if err != nil {
		fmt.Println("Failed to create record reader: ", err)
	}

	header, err := recordReader.ReadHeader()
	if err != nil {
		fmt.Println("Failed to read header: ", err)
	}

	index, err := recordReader.ReadIndex(*header.IndexPosition)
	if err != nil {
		fmt.Println("Failed to read index: ", err)
	}

	r := &Record{
		Filepath:       recordpath,
		Reader:         recordReader,
		Header:         header,
		ChunkHeaderIdx: make([]*recordproto.SingleIndex, 0),
		ChunkBodyIdx:   make([]*recordproto.SingleIndex, 0),
		Channels:       make(map[string]*recordproto.ChannelCache),
	}

	r.parseRecordIndex(index)

	return r
}

func (r *Record) Close() {
	r.Reader.Close()
}

func (r *Record) parseRecordIndex(index *recordproto.Index) {
	for _, item := range index.Indexes {
		itemType := item.GetType()
		channelCache := item.GetChannelCache()

		switch itemType {
		case recordproto.SectionType_SECTION_CHUNK_HEADER:
			r.ChunkHeaderIdx = append(r.ChunkHeaderIdx, item)
		case recordproto.SectionType_SECTION_CHUNK_BODY:
			r.ChunkBodyIdx = append(r.ChunkBodyIdx, item)
		case recordproto.SectionType_SECTION_CHANNEL:
			channelName := channelCache.GetName()
			r.Channels[channelName] = channelCache
			r.parseProtoDesc(channelCache)
		}
	}
}

func (r *Record) parseProtoDesc(channelCache *recordproto.ChannelCache) {
	bytes := channelCache.GetProtoDesc()
	var protoDesc recordproto.ProtoDesc
	err := proto.Unmarshal(bytes, &protoDesc)
	if err != nil {
		fmt.Println("Failed to unmarshal proto desc: ", err)
	}

	r.addProtoDesc(&protoDesc)
}

func (r *Record) addProtoDesc(protoDesc *recordproto.ProtoDesc) {
	// first add deps recursively
	deps := protoDesc.GetDependencies()
	for _, dep := range deps {
		r.addProtoDesc(dep)
	}

	// transform FileDescriptorProto to FileDescriptor
	descData := protoDesc.GetDesc()
	if len(descData) == 0 {
		fmt.Println("Empty descriptor data")
		return
	}

	fileDescProto := descriptorpb.FileDescriptorProto{}
	err := proto.Unmarshal(descData, &fileDescProto)
	if err != nil {
		fmt.Println("Failed to unmarshal file desc proto:", err)
		return
	}

	// register FileDescriptor
	fd, err := protodesc.NewFile(&fileDescProto, protoregistry.GlobalFiles)
	if err != nil {
		fmt.Println("Failed to create file descriptor:", err)
		return
	}

	thePath := fd.Path()
	if _, err := protoregistry.GlobalFiles.FindFileByPath(thePath); err == nil {
		return
	}

	err = protoregistry.GlobalFiles.RegisterFile(fd)
	if err != nil {
		fmt.Println("Failed to register file:", err)
		return
	}

	// register each MessageDescriptor
	for i := 0; i < fd.Messages().Len(); i++ {
		md := fd.Messages().Get(i)
		mt := dynamicpb.NewMessageType(md)
		if err = protoregistry.GlobalTypes.RegisterMessage(mt); err != nil {
			fmt.Println("Failed to register message type:", err)
			return
		}
	}
}

func (r *Record) ReadMessage() <-chan Message {
	ch := make(chan Message)

	go func() {
		defer close(ch)

		for _, chunkBodyIdx := range r.ChunkBodyIdx {
			position := chunkBodyIdx.GetPosition()
			chunk, err := r.Reader.ReadChunkBody(position)
			if err != nil {
				fmt.Println("Failed to read chunk body: ", err)
				return
			}

			for _, msg := range chunk.GetMessages() {
				channelName := msg.GetChannelName()
				data := msg.GetContent()
				ch <- Message{
					ChannelName: channelName,
					Time:        msg.GetTime(),
					Content:     data,
				}
			}
		}
	}()

	return ch
}

func (r *Record) ConvertMessageToJSON(messageTypeStr string, data []byte, options *protojson.MarshalOptions) ([]byte, error) {
	fullname := protoreflect.FullName(messageTypeStr)
	// get message type
	messageType, err := protoregistry.GlobalTypes.FindMessageByName(fullname)
	if err != nil {
		return nil, fmt.Errorf("failed to find message type: %v", err)
	}

	// create a message instance
	msg := messageType.New().Interface()

	// unmarshal the message
	err = proto.Unmarshal(data, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %v", err)
	}

	// marshal the message to json
	if options == nil {
		options = &protojson.MarshalOptions{}
	}
	jsonData, err := options.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message to json: %v", err)
	}

	return jsonData, nil
}
