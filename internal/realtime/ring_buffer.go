package realtime

import (
	"SIEMSEasy/internal/database"
	"runtime"
	"sync/atomic"
)

const bufferSize = 65536
const mask = bufferSize - 1

type logItem struct {
	EventData      database.LogSchema
	SequenceNumber atomic.Uint64
}

type RingBuffer struct {
	buffer [bufferSize]logItem
	front  uint64
	rear   atomic.Uint64
}

func (b *RingBuffer) Enqueue(insertionValue database.LogSchema) bool {
	for {

		ticket := b.rear.Load()
		slot := ticket & mask
		diff := int64(b.buffer[slot].SequenceNumber.Load()) - int64(ticket)

		if diff == 0 {
			if b.rear.CompareAndSwap(ticket, ticket+1) {
				b.buffer[slot].EventData = insertionValue
				b.buffer[slot].SequenceNumber.Store(ticket + 1)
				return true
			}
			runtime.Gosched()
		} else if diff < 0 {
			return false
		} else {
			runtime.Gosched()
		}

	}

}

func (b *RingBuffer) Dequeue(batchSize int) []database.LogSchema {

	batch := make([]database.LogSchema, 0, batchSize)

	for len(batch) < batchSize {
		slot := b.front & mask
		diff := int64(b.buffer[slot].SequenceNumber.Load()) - int64(b.front+1)
		if diff == 0 {
			batch = append(batch, b.buffer[slot].EventData)
			b.buffer[slot].EventData = database.LogSchema{}
			b.buffer[slot].SequenceNumber.Store(b.front + bufferSize)
			b.front++
		} else if diff < 0 {
			break
		}
	}

	return batch

}

func CreateRingBuffer() *RingBuffer {
	var buf RingBuffer
	buf.front = 0
	buf.rear.Store(0)
	for i := range buf.buffer {
		buf.buffer[i].SequenceNumber.Store(uint64(i))
	}
	return &buf
}
