package stream

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/reassembly"
	"sync"
)

type TCPStreamFactory struct {
	wg sync.WaitGroup
}

func (self *TCPStreamFactory) New(netFlow, tcpFlow gopacket.Flow, tcp *layers.TCP, ac reassembly.AssemblerContext) reassembly.Stream {
	fmt.Println("New:", netFlow, " ", tcpFlow)
	sip, dip := netFlow.Endpoints()
	srcip := fmt.Sprintf("%s", sip)
	dstip := fmt.Sprintf("%s", dip)
	fmt.Println(srcip, "->", dstip)

	fsmOptions := reassembly.TCPSimpleFSMOptions{
		SupportMissingEstablishment: true,
	}
	stream := &TCPStream{
		net:        netFlow,
		transport:  tcpFlow,
		isHTTP:     tcp.SrcPort == layers.TCPPort(80) || tcp.DstPort == layers.TCPPort(80),
		reversed:   tcp.SrcPort == layers.TCPPort(80),
		tcpstate:   reassembly.NewTCPSimpleFSM(fsmOptions),
		ident:      fmt.Sprintf("%s:%s -> %s:%S", netFlow.Src(), tcpFlow.Src(), netFlow.Dst(), tcpFlow.Dst()),
		optchecker: reassembly.NewTCPOptionCheck(),
	}
	if stream.isHTTP {
		stream.client = httpReader{
			bytes:     make(chan []byte),
			timestamp: make(chan int64),
			ident:     fmt.Sprintf("%s %s", netFlow, tcpFlow),
			//parent:    stream,
			isClient: true,
			srcport:  fmt.Sprintf("%d", tcp.SrcPort),
			dstport:  fmt.Sprintf("%d", tcp.DstPort),
			srcip:    srcip,
			dstip:    dstip,
			//httpstart: 0,
		}
		stream.server = httpReader{
			bytes:     make(chan []byte),
			timestamp: make(chan int64),
			ident:     fmt.Sprintf("%s %s", netFlow.Reverse(), netFlow.Reverse()),
			//parent:    stream,
			dstport: fmt.Sprintf("%d", tcp.SrcPort),
			srcport: fmt.Sprintf("%d", tcp.DstPort),
			dstip:   srcip,
			srcip:   dstip,
			//httpstart: 0,
		}
		self.wg.Add(2)
		go func() {
			defer self.wg.Done()
			stream.client.runClient()
		}()
		go func() {
			defer self.wg.Done()
			stream.server.runServer()
		}()
	}
	return stream

}

func (self *TCPStreamFactory) WaitGoRoutines() {
	self.wg.Wait()
}
