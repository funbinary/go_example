package source

import (
	log "github.com/funbinary/go_example/pkg/blog"
	"github.com/google/gopacket/pcap"
)

func NewPcapSource(device string, opts ...option) *pcap.Handle {
	ih, err := pcap.NewInactiveHandle(device)
	if err != nil {
		log.Errorf("capture %s error: %v", device, err)
		return nil
	}
	s := NewDefaultPcapSource()
	for _, opt := range opts {
		opt(s)
	}
	ih.SetPromisc(s.promisc)
	ih.SetBufferSize(s.bufferSize)
	ih.SetSnapLen(s.snaplen)
	ih.SetTimeout(s.timeout)
	handler, err := ih.Activate()
	if err != nil {
		log.Errorf("activate %s error: %v", device, err)
		return nil
	}
	return handler
}
