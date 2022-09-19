package stream2

import (
	"github.com/google/gopacket"
	"io"
)

type Plug struct {
	ResolveStream func(net gopacket.Flow, transport gopacket.Flow, r io.Reader)

	InternalPlugList map[string]PlugInterface
}

type PlugInterface interface {
	ResolveStream(net gopacket.Flow, transport gopacket.Flow, r io.Reader)
}

func NewPlug() *Plug {
	var p Plug
	return &p
}
