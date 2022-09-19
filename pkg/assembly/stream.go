package assembly

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Stream is implemented by the caller to handle incoming reassembled
// TCP data.  Callers create a StreamFactory, then StreamPool uses
// it to create a new Stream for every TCP stream.
//
// assembly will, in order:
//    1) Create the stream via StreamFactory.New
//    2) Call ReassembledSG 0 or more times, passing in reassembled TCP data in order
//    3) Call ReassemblyComplete one time, after which the stream is dereferenced by assembly.
type Stream interface {
	// Accept
	// Tell whether the TCP packet should be accepted, start could be modified to force a start even if no SYN have been seen
	// 告诉TCP数据包是否应该被接受，start可以修改为强制启动，即使没有看到SYN

	Accept(tcp *layers.TCP, ci gopacket.CaptureInfo, dir TCPFlowDirection, nextSeq Sequence, start bool, ac AssemblerContext) bool

	// ReassembledSG is called zero or more times.
	// ScatterGather is reused after each Reassembled call,
	// so it's important to copy anything you need out of it,
	// especially bytes (or use KeepFrom())
	ReassembledSG(sg ScatterGather, ac AssemblerContext)

	// ReassemblyComplete is called when assembly decides there is
	// no more data for this Stream, either because a FIN or RST packet
	// was seen, or because the stream has timed out without any new
	// packet data (due to a call to FlushCloseOlderThan).
	// It should return true if the connection should be removed from the pool
	// It can return false if it want to see subsequent packets with Accept(), e.g. to
	// see FIN-ACK, for deeper state-machine analysis.
	ReassemblyComplete(ac AssemblerContext) bool
}

// ScatterGather is used to pass reassembled data and metadata of reassembled
// packets to a Stream via ReassembledSG
type ScatterGather interface {
	// Returns the length of available bytes and saved bytes
	Lengths() (int, int)
	// Returns the bytes up to length (shall be <= available bytes)
	Fetch(length int) []byte
	// Tell to keep from offset
	KeepFrom(offset int)
	// Return CaptureInfo of packet corresponding to given offset
	CaptureInfo(offset int) gopacket.CaptureInfo
	// Return some info about the reassembled chunks
	Info() (direction TCPFlowDirection, start bool, end bool, skip int)
	// Return some stats regarding the state of the stream
	Stats() TCPAssemblyStats
}

// TCPAssemblyStats provides some figures for a ScatterGather
type TCPAssemblyStats struct {
	// For this ScatterGather
	Chunks  int
	Packets int
	// For the half connection, since last call to ReassembledSG()
	QueuedBytes    int
	QueuedPackets  int
	OverlapBytes   int
	OverlapPackets int
}
