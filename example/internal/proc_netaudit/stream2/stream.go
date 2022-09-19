package stream2

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

type ProtocolStreamFactory struct {
	Plug *Plug
}

type ProtocolStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
}

func (m *ProtocolStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {

	//init stream struct
	stm := &ProtocolStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}

	//new stream
	fmt.Println("# Start new stream:", net, transport)

	//decode packet
	go m.Plug.ResolveStream(net, transport, &(stm.r))

	return &(stm.r)
}
