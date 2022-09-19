package stream

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/reassembly"
	"strings"
	"sync"
)

type TCPStream struct {
	tcpstate       *reassembly.TCPSimpleFSM
	optchecker     reassembly.TCPOptionCheck
	net, transport gopacket.Flow
	ident          string

	isHTTP   bool
	reversed bool
	client   httpReader
	server   httpReader
	//urls     []string
	//all      []httpGroup
	//hg sync.Mutex
	sync.Mutex
}

// 每次收到包时调用
func (self *TCPStream) Accept(tcp *layers.TCP, ci gopacket.CaptureInfo, dir reassembly.TCPFlowDirection, nextSeq reassembly.Sequence, start *bool, ac reassembly.AssemblerContext) bool {
	// 检测协议
	// FSM
	if !self.tcpstate.CheckState(tcp, dir) {
		fmt.Println("FSM", self.ident, ": Packet rejected by FSM (state:", self.tcpstate.String(), ")")
		//Error("FSM", "%s: Packet rejected by FSM (state:%s)\n", t.ident, t.tcpstate.String())
		//stats.rejectFsm++
		//if !t.fsmerr {
		//	t.fsmerr = true
		//	stats.rejectConnFsm++
		//}
		//if !*ignorefsmerr {
		//	return false
		//}
		return false
	}
	var httpDirect HTTPDirect

	*start, httpDirect = detectHttp(tcp.Payload)
	if *start {
		fmt.Println("http direct:", httpDirect)
	}

	err := self.optchecker.Accept(tcp, ci, dir, nextSeq, start)
	if err != nil {
		// 重复的包，丢弃 drop
		// 调试发现此包为以前序号的包，并且出现过。
		// mss BUG,server mss通过路由拆解成mss要求的包尺寸，
		// 因此不能判断包大小大于mss为错包
		if strings.Contains(fmt.Sprintf("%s", err), " > mss ") {
			//  > mss 包 不丢弃
		} else {
			fmt.Println("OptionChecker", self.net, " -> ", self.transport, ":", "Packet rejected by option checker:", err)
			//Error("OptionChecker", "%v ->%v : Packet rejected by OptionChecker: %s\n", t.net, t.transport, err)
			//stats.rejectOpt++
			//if !*nooptcheck {
			//	return false
			//}
			return false
		}
	}

	accept := true

	// create new httpgroup,wait request+response
	//if *start {
	//	self.NewhttpGroup(isReq, ci.Timestamp.UnixNano())
	//}

	return accept

}

func (self *TCPStream) ReassembledSG(sg reassembly.ScatterGather, ac reassembly.AssemblerContext) {
	dir, start, end, skip := sg.Info()
	length, saved := sg.Lengths()

	// 更新stats
	sgStats := sg.Stats()

	var ident string
	if dir == reassembly.TCPDirClientToServer {
		ident = fmt.Sprintf("%v %v(%s): ", self.net, self.transport, dir)
	} else {
		ident = fmt.Sprintf("%v %v(%s): ", self.net.Reverse(), self.transport.Reverse(), dir)
	}
	fmt.Printf("%s: SG reassembled packet with %d bytes (start:%v,end:%v,skip:%d,saved:%d,nb:%d,%d,overlap:%d,%d)\n", ident, length, start, end, skip, saved, sgStats.Packets, sgStats.Chunks, sgStats.OverlapBytes, sgStats.OverlapPackets)
	if skip == -1 {
		// this is allowed
	} else if skip != 0 {
		// Missing bytes in stream: do not even try to parse it
		return
	}

	//use timeStamp as match flag
	timeStamp := sg.CaptureInfo(0).Timestamp.UnixNano()
	data := sg.Fetch(length)
	if self.isHTTP {
		if length > 0 {
			ok, _ := detectHttp(data)

			//if dir == reassembly.TCPDirClientToServer && !t.reversed {
			if dir == reassembly.TCPDirClientToServer {
				self.client.bytes <- data
				if ok {
					self.client.timestamp <- timeStamp
				}
			} else {
				self.server.bytes <- data
				if ok {
					self.server.timestamp <- timeStamp
				}
			}
		}
	}
}

// 连接被关闭时的操作
func (self *TCPStream) ReassemblyComplete(ac reassembly.AssemblerContext) bool {
	fmt.Printf("%s: Connection closed\n", self.ident)
	if self.isHTTP {
		close(self.client.bytes)
		close(self.server.bytes)
	}
	// do not remove the connection to allow last ACK
	return false
}
