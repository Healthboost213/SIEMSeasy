package tcpserver

import (
	"SIEMSEasy/internal/database"
	"SIEMSEasy/internal/network/serializer"
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

// SIEMSeasy Log Transport Protocol (SLTP)
// SLTP Wire Format:
// [ 4 Byte : Magic Bytes ] [ 2 Bytes : Flags ] [ 4 Bytes : Content Length ] [ Remainder : Payload ]

// Flags & Signals
// Bit 0 : Message from Server when SET. Else from Worker.
// Bit 1 : Content Toggle - Logs if SET / Rules if UNSET
// Bits 2 - 9 : Reserved for future use
// Bits 10 : Heartbeat
// Bits 11 : Backpressure Signal
// Bits 12 : Resume Signal
// Bits 13 - 15 : Reserved for future use

var magicBytes = []byte{0x53, 0x4c, 0x54, 0x50}

// Writing Wire Data

const headerSize = 10

const (
	FromServer uint16 = 1 << iota
	transferringLog
	_
	_
	_
	_
	_
	_
	_
	_
	HeartbeatSignal
	BackpressureSignal
	ResumeSignal
	_
	_
	_
)

func EncodeLogSLTP(logStruct database.LogSchema) ([]byte, error) {

	protoData := serializer.LogSchema{
		LogID:           logStruct.LogID[:],
		DeviceID:        logStruct.DeviceID[:],
		SourceAddress:   logStruct.SourceAddress,
		ProductName:     logStruct.ProductName,
		ServiceName:     logStruct.ServiceName,
		Message:         logStruct.Message,
		MessageHash:     logStruct.MessageHash[:],
		ScrapeTimestamp: logStruct.ScrapeTimestamp.UnixNano(),
		RecvTimestamp:   logStruct.RecvTimestamp.UnixNano(),
	}

	wireLength := proto.Size(&protoData)
	messageFrame := make([]byte, headerSize, headerSize+wireLength)
	copy(messageFrame, magicBytes)

	flags := transferringLog
	binary.BigEndian.PutUint16(messageFrame[4:6], flags)
	binary.BigEndian.PutUint32(messageFrame[6:10], uint32(wireLength))

	marshalledData, err := proto.MarshalOptions{}.MarshalAppend(messageFrame, &protoData)
	if err != nil {
		log.Println("Protobuf Error:", err)
	}

	return marshalledData, err

}

// Reading Wire Data

func parseTCPStream(streamBuffer []byte, partial []byte) ([][]byte, []byte) {

	var indexes [][]int
	var messages [][]byte
	var startDelimiter int

	streamBuffer = append(partial, streamBuffer...)

	for i := 0; i < len(streamBuffer)-3; i++ {
		if bytes.Equal(streamBuffer[i:i+4], magicBytes) {
			startDelimiter = i
			indexes = append(indexes, []int{startDelimiter, i + 4})
		}
	}

	for i := 0; i < len(indexes)-1; i++ {
		messages = append(messages, streamBuffer[indexes[i][0]:indexes[i+1][0]])
	}

	remaining := append([]byte(nil), streamBuffer[indexes[len(indexes)-1][0]:]...)

	return messages, remaining

}

func ConvertFromProtobuf(protoValue []byte) database.LogSchema {

	var decodedLog serializer.LogSchema // protobuf struct

	err := proto.Unmarshal(protoValue, &decodedLog)
	if err != nil {
		log.Println("Protobuf Error:", err)
	}

	var convertedLog database.LogSchema
	convertedLog.LogID, _ = uuid.FromBytes(decodedLog.LogID)
	convertedLog.DeviceID, _ = uuid.FromBytes(decodedLog.DeviceID)
	convertedLog.SourceAddress = net.IP(decodedLog.SourceAddress)
	convertedLog.ProductName = decodedLog.ProductName
	convertedLog.ServiceName = decodedLog.ServiceName
	convertedLog.Message = decodedLog.Message
	convertedLog.MessageHash = [32]byte(decodedLog.MessageHash)
	convertedLog.ScrapeTimestamp = time.Unix(0, decodedLog.ScrapeTimestamp)
	convertedLog.RecvTimestamp = time.Unix(0, decodedLog.RecvTimestamp)

	return convertedLog

}
