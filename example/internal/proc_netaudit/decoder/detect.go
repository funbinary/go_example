package decoder

import (
	"github.com/google/gopacket"
)

func DetectProtocol(data []byte, t gopacket.TransportLayer) ApplicationPortocol {
	for p, d := range decoderMap {
		if d.CanDecode(data, t) {
			return p
		}
	}
	return UnknownProtocol
}
