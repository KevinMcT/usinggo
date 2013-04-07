package FifoList

type Value interface{}

// makes the rest of this package more readable
type buffer []string

type Fifo struct {
	full    chan buffer // channel providing full buffers
	empty   chan bool   // signal for buffer switching
	reading buffer      // only the reader uses this
	writing buffer      // only the writer uses this
}

func (f *Fifo) Push(item string) {
	f.writing = append(f.writing, item)
	select {
	case <-f.empty:
		f.full <- f.writing
		f.writing = make(buffer, 0, 1)
	default:
	}
}

func (f *Fifo) Pop() (item string) {
	if len(f.reading) > 0 {
		item = f.reading[0]
		f.reading = f.reading[1:]
		return
	}
	f.empty <- true
	f.reading = <-f.full
	if len(f.reading) == 0 {
		panic("I only accept non-empty buffers")
	}
	return f.Pop()
}

func NewFifo() *Fifo {
	return &Fifo{
		full:    make(chan buffer),
		empty:   make(chan bool),
		writing: make(buffer, 0, 1),
	}
}
