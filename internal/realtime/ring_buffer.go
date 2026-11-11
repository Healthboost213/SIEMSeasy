package realtime

import (
	"runtime"
	"sync/atomic"
)

const bufferSize = 65536
const mask = bufferSize - 1

type logItem struct {
	EventData      []byte
	SequenceNumber atomic.Uint64
}

type RingBuffer struct {
	buffer [bufferSize]logItem
	Front  uint64
	Rear   atomic.Uint64
}

func (b *RingBuffer) Enqueue(insertionValue []byte) bool {
	for {

		ticket := b.Rear.Load()
		slot := ticket & mask
		diff := int64(b.buffer[slot].SequenceNumber.Load()) - int64(ticket)

		if diff == 0 {
			if b.Rear.CompareAndSwap(ticket, ticket+1) {
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

func (b *RingBuffer) Dequeue(batchSize int) [][]byte {

	batch := make([][]byte, 0, batchSize)

	for len(batch) < batchSize {
		slot := b.Front & mask
		diff := int64(b.buffer[slot].SequenceNumber.Load()) - int64(b.Front+1)
		if diff == 0 {
			batch = append(batch, b.buffer[slot].EventData)
			b.buffer[slot].EventData = nil
			b.buffer[slot].SequenceNumber.Store(b.Front + bufferSize)
			b.Front++
		} else if diff < 0 {
			break
		}
	}

	return batch

}

func CreateRingBuffer() *RingBuffer {
	var buf RingBuffer
	buf.Front = 0
	buf.Rear.Store(0)
	for i := range buf.buffer {
		buf.buffer[i].SequenceNumber.Store(uint64(i))
	}
	return &buf
}
