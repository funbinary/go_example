package assembly

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// StreamFactory is used by assembly to create a new stream for each
// new TCP session.
type StreamFactory interface {
	// New should return a new stream for the given TCP key.
	New(netFlow, tcpFlow gopacket.Flow, tcp *layers.TCP, ac AssemblerContext) Stream
}
