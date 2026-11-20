package tcpserver

import (
	"SIEMSEasy/internal/database"
	"SIEMSEasy/internal/network/serializer"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

var magicBytes = []byte{0x53, 0x4c, 0x54, 0x50}

// Protobuf Specific Functions

func ConvertToProtobuf(logStruct database.LogSchema) []byte {

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

	marshalledData, err := proto.Marshal(&protoData)
	if err != nil {
		log.Println("Protobuf Error:", err)
	}

	return marshalledData

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

func FrameData(protoPayload []byte) {
	baseFrame := []byte{0x53, 0x4c, 0x54, 0x50}
	tempbuf := make([]byte, 4)
	binary.BigEndian.PutUint32(tempbuf, uint32(len(protoPayload)))
	fmt.Println(tempbuf)
	baseFrame = append(baseFrame, tempbuf...)
	baseFrame = append(baseFrame, protoPayload...)
	fmt.Println(baseFrame)
}

func FrameHeader(flags uint16) []byte {
	var baseFrame = [headerSize]byte{0x53, 0x4c, 0x54, 0x50}
	binary.BigEndian.PutUint16(baseFrame[4:6], flags)
	binary.BigEndian.PutUint32(baseFrame[6:10], 0) // length
	return baseFrame[:]
}

func ReadHeader() {

}

// Get Log Struct -> Marshal To ProtoBuf -> Get Length -> Modify Length field -> send away

func SLTPEncodeLog(protoData *[]byte) {

}
