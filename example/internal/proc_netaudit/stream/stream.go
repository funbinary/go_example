package stream

import (
	"encoding/hex"
	"fmt"
	"github.com/funbinary/go_example/example/internal/proc_netaudit/decoder"
	log "github.com/funbinary/go_example/pkg/blog"
	"github.com/funbinary/go_example/pkg/errors"
	"github.com/funbinary/go_example/pkg/reassembly"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"strings"
)

type TCPStream struct {
	tcpstate       *reassembly.TCPSimpleFSM
	optchecker     reassembly.TCPOptionCheck
	net, transport gopacket.Flow
	ident          string
	tcp            *layers.TCP
	proto          decoder.ApplicationPortocol
	s2c            *Packet
	c2s            *Packet
}

func NewTCPStream(netFlow, tcpFlow gopacket.Flow, tcp *layers.TCP) *TCPStream {
	fsmOptions := reassembly.TCPSimpleFSMOptions{
		SupportMissingEstablishment: true,
	}

	return &TCPStream{
		tcpstate:   reassembly.NewTCPSimpleFSM(fsmOptions),
		optchecker: reassembly.NewTCPOptionCheck(),
		net:        netFlow,
		transport:  tcpFlow,
		ident:      fmt.Sprintf("%s_%s_%s_%s_%s", netFlow.Src(), tcpFlow.Src(), netFlow.Dst(), tcpFlow.Dst(), "tcp"),
		tcp:        tcp,
		//s2c:        make(chan []byte),
		//c2s:        make(chan []byte),
	}
}

//
func (self *TCPStream) Accept(tcp *layers.TCP,
	ci gopacket.CaptureInfo, dir reassembly.TCPFlowDirection,
	nextSeq reassembly.Sequence, start *bool, ac reassembly.AssemblerContext) bool {
	log.Debugf("==========================accept============================")

	return true
	if !self.tcpstate.CheckState(tcp, dir) {
		log.Errorf("状态机检测%s包被拒绝,FSM(state:%s)", self.ident, self.tcpstate.String())
		return false
	}

	//self.proto = protocol.DetectProtocol(tcp.Payload)
	//switch self.proto {
	//case protocol.UnknownProtocol:
	//	log.Errorf("协议检测失败，包丢弃")
	//	return false
	//}

	//if self.proto == protocol.DecoderFailed {
	//	return false
	//}

	// option检查
	err := self.optchecker.Accept(tcp, ci, dir, nextSeq, start)
	if err != nil {
		// 重复的包，丢弃 drop
		// 调试发现此包为以前序号的包，并且出现过。
		// mss BUG,server mss通过路由拆解成mss要求的包尺寸，
		// 因此不能判断包大小大于mss为错包
		if strings.Contains(fmt.Sprintf("%s", err), " > mss ") {
			//  > mss 包 不丢弃
		} else {
			log.Errorf("Option检测，包将被丢弃%v ->%v error: %s", self.net, self.transport, err)
			return false
		}
	}

	return true
}

func (self *TCPStream) ReassembledSG(sg reassembly.ScatterGather, ac reassembly.AssemblerContext) {
	log.Debugf("-------------sg--------------")

	dir, _, skip := sg.Info()
	var length int
	length, _ = sg.Lengths()
	//dir, start, end, skip := sg.Info()
	//length, saved := sg.Lengths()
	// 更新stats
	//sgStats := sg.Stats()
	//var ident = self.ident
	//if dir {
	//	ident = fmt.Sprintf("%v %v(%s): ", self.net.Reverse(), self.transport.Reverse(), dir)
	//}
	//log.Infof("%s: SG reassembled packet with %d bytes (start:%v,end:%v,skip:%d,saved:%d,nb:%d,%d,overlap:%d,%d)", ident, length, start, end, skip, saved, sgStats.Packets, sgStats.Chunks, sgStats.OverlapBytes, sgStats.OverlapPackets)

	if skip == -1 {
		// this is allowed
	} else if skip != 0 {
		// Missing bytes in stream: do not even try to parse it
		log.Debugf("MISSING BYTES in stream")
		//return
	}
	if length <= 0 {
		log.Debugf("length must be > 0")
		return
	}
	data := sg.Fetch(length)
	if self.proto == decoder.UnknownProtocol {
		// 第一次一定可以检测出协议
		self.proto = decoder.DetectProtocol(data, self.tcp)
		if self.proto == decoder.UnknownProtocol {
			return
		}
	}

	if dir == reassembly.TCPDirServerToClient {
		log.Debugf("read pkg")
		remainLen := len(data)
		for remainLen > 0 {
			log.Debugf("remainLen = %d", remainLen)
			if self.s2c == nil {
				self.s2c = NewPacket(self.proto, self.tcp)
			}
			appendLen := self.s2c.Append(data)
			remainLen = len(data) - appendLen
			data = data[appendLen:]
			if !self.decode(self.s2c) {
				continue
			}
			self.s2c = nil
		}
	} else {
		log.Debugf("read pkg")
		remainLen := len(data)
		for remainLen > 0 {
			log.Debugf("remainLen = %d", remainLen)
			if self.c2s == nil {
				self.c2s = NewPacket(self.proto, self.tcp)
			}
			appendLen := self.c2s.Append(data)
			remainLen = len(data) - appendLen
			data = data[appendLen:]
			if !self.decode(self.c2s) {
				continue
			}
			self.c2s = nil
		}
	}

}

//func (self *TCPStream) ResolveS2C() {\n\tvar curPkg *Packet\n\n\tfor {\n\t\tselect {\n\t\tcase data := <-self.s2c:\n\t\t\tlog.Debugf(\"read pkg\")\n\t\t\tremainLen := len(data)\n\t\t\tfor remainLen > 0 {\n\t\t\t\tlog.Debugf(\"remainLen = %d", remainLen)
//
//				if curPkg == nil {
//					curPkg = NewPacket(self.proto, self.tcp)
//				}
//				appendLen := curPkg.Append(data)
//				remainLen = len(data) - appendLen
//				data = data[appendLen:]
//				if !self.decode(curPkg) {
//					continue
//				}
//				curPkg = nil
//
//			}
//		}
//	}
//}
//
//func (self *TCPStream) ResolveC2S() {
//	var curPkg *Packet
//	for {
//		select {
//		case data := <-self.c2s:
//
//		}
//	}
//}

func (self *TCPStream) decode(p *Packet) bool {
	if p == nil || !p.Completed() {

		return false
	}
	data := p.Data()
	log.Debugf("hex: %s", hex.EncodeToString(data))
	if err := decoder.DecodeUseProto(self.ident, data, self.tcp, self.proto); err != nil {
		log.Errorf("%+v", err)
		if !errors.Is(err, decoder.ErrDecoderUnknown) {
			log.Errorf("%+v", err)
		}
	}
	return true
}

//func (self *TCPStream) resolvePacket(p *Packet, data []byte) (pkg []byte) {
//	if p == nil {
//		d := decoder.GetDecoder(self.proto)
//		pkgLen, off, err := d.DetectPacketLength(data, self.tcp)
//		if err != nil {
//			log.Errorf("detect packet length %+v", err)
//			return nil
//		}
//		if int(pkgLen+off) != len(data) {
//			p = NewPacket(int(pkgLen))
//			p.Append(data[off:])
//			return
//		}
//		pkg = data
//	} else {
//		appendLen := p.Append(data)
//		if appendLen < len(data) {
//			pkg = p.Data()
//			remain :=
//		}
//		if p.Completed() {
//			pkg = p.Data()
//		}
//	}
//	return pkg
//}

// 关闭连接时调用，一般会在FlushOption时调用
func (self *TCPStream) ReassemblyComplete(ac reassembly.AssemblerContext) bool {
	return false
}
