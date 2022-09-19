package decoder

import (
	"github.com/funbinary/go_example/example/internal/proc_netaudit/decoder/nfs"
	"github.com/funbinary/go_example/example/internal/proc_netaudit/hook"
	"github.com/funbinary/go_example/pkg/errors"
	"github.com/google/gopacket"
)

var UnknownProtocol = RegisterDecoder(0, NewDecoderUnknown())

var NfsProtocol = RegisterDecoder(2049, nfs.NewDecoder())

var ErrDecoderUnknown = errors.New("Layer type not currently supported")

type DecoderUnknown struct {
	name string
}

func NewDecoderUnknown() *DecoderUnknown {
	return &DecoderUnknown{
		name: "unknown",
	}
}

func (self *DecoderUnknown) Name() string { return self.name }

func (self *DecoderUnknown) CanDecode(data []byte, t gopacket.TransportLayer) bool {
	return false
}
func (self *DecoderUnknown) DetectPacketLength(data []byte, t gopacket.TransportLayer) (length uint64, off uint64, err error) {
	return 0, 0, ErrDecoderUnknown
}

func (self *DecoderUnknown) Decode(data []byte, t gopacket.TransportLayer) (app hook.Protocol, err error) {
	return nil, ErrDecoderUnknown
}

func (self *DecoderUnknown) String() string { return self.Name() }
