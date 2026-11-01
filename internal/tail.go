package internal

import (
	"io"
	"log"
	"os"
	"regexp"
)

func FileTail(fileName string, regexString string) string {

	// compile regex
	re := regexp.MustCompile(regexString)

	// open file
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// overall buffer for whole event
	overallEventBuffer := make([]byte, 0)

	// reading in 4KB chunks to save memory
	var chunkSize int64 = 4096
	var currChunk int64 = 1
	currChunkBuffer := make([]byte, 4096)
	file.Seek(-(currChunk * chunkSize), io.SeekEnd)
	bytesRead, _ := file.Read(currChunkBuffer)

	// find final event
	var finalMatch string = ""
	matchArray := re.FindAllIndex(currChunkBuffer[0:bytesRead], -1)
	overallEventBuffer = append(currChunkBuffer[:bytesRead], overallEventBuffer...)

	// repeat until a match is found
	for len(matchArray) == 0 {
		currChunk += 1
		file.Seek(-(currChunk * chunkSize), io.SeekEnd)
		bytesRead, _ = file.Read(currChunkBuffer)
		matchArray = re.FindAllIndex(currChunkBuffer[0:bytesRead], -1)
		overallEventBuffer = append(currChunkBuffer[:bytesRead], overallEventBuffer...)
	}

	finalMatch = string(overallEventBuffer[matchArray[len(matchArray)-1][0]:len(overallEventBuffer)])

	return finalMatch

}
