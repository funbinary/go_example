package assembly

import (
	"fmt"
	"github.com/google/gopacket"
)

const invalidSequence = -1
const uint32Max = 0xFFFFFFFF

// Sequence is a TCP sequence number.  It provides a few convenience functions
// for handling TCP wrap-around.  The sequence should always be in the range
// [0,0xFFFFFFFF]... its other bits are simply used in wrap-around calculations
// and should never be set.
type Sequence int64

// Difference defines an ordering for comparing TCP sequences that's safe for
// roll-overs.  It returns:
//    > 0 : if t comes after s
//    < 0 : if t comes before s
//      0 : if t == s
// The number returned is the sequence difference, so 4.Difference(8) will
// return 4.
//
// It handles rollovers by considering any sequence in the first quarter of the
// uint32 space to be after any sequence in the last quarter of that space, thus
// wrapping the uint32 space.
func (s Sequence) Difference(t Sequence) int {
	if s > uint32Max-uint32Max/4 && t < uint32Max/4 {
		t += uint32Max
	} else if t > uint32Max-uint32Max/4 && s < uint32Max/4 {
		s += uint32Max
	}
	return int(t - s)
}

// Add adds an integer to a sequence and returns the resulting sequence.
func (s Sequence) Add(t int) Sequence {
	return (s + Sequence(t)) & uint32Max
}

type key [2]gopacket.Flow

func (k *key) String() string {
	return fmt.Sprintf("%s:%s", k[0], k[1])
}

func (k *key) Reverse() key {
	return key{
		k[0].Reverse(),
		k[1].Reverse(),
	}
}

// TCPFlowDirection distinguish the two half-connections directions.
//
// TCPDirClientToServer is assigned to half-connection for the first received
// packet, hence might be wrong if packets are not received in order.
// It's up to the caller (e.g. in Accept()) to decide if the direction should
// be interpretted differently.
type TCPFlowDirection bool

// Value are not really useful
const (
	TCPDirClientToServer TCPFlowDirection = false
	TCPDirServerToClient TCPFlowDirection = true
)

func (dir TCPFlowDirection) String() string {
	switch dir {
	case TCPDirClientToServer:
		return "client->server"
	case TCPDirServerToClient:
		return "server->client"
	}
	return ""
}

// Reverse returns the reversed direction
func (dir TCPFlowDirection) Reverse() TCPFlowDirection {
	return !dir
}
