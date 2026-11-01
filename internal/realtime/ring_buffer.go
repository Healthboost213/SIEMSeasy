package realtime

type RingBuffer struct {
	buffer [65536]string
	Head   uint16
	Tail   uint16
}

func (buf *RingBuffer) Enqueue(val string) bool {

	if (buf.Tail + 1) == buf.Head {
		return false
	}

	buf.buffer[buf.Tail] = val
	buf.Tail++
	return true
}

func (buf *RingBuffer) Dequeue() string {

	if buf.Head == buf.Tail {
		return "false"
	}

	valueFound := buf.buffer[buf.Head]
	buf.Head++
	return valueFound

}
