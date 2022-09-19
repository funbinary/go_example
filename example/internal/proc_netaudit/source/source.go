package source

import "time"

type PcapSource struct {
	timeout    time.Duration
	snaplen    int
	promisc    bool
	bufferSize int
}

func NewDefaultPcapSource() *PcapSource {
	return &PcapSource{
		timeout:    -1 * time.Second,
		snaplen:    65536,
		promisc:    true,
		bufferSize: 1024 * 1024 * 10,
	}
}
