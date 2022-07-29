package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
	"time"

	"github.com/google/gopacket/reassembly"

	"github.com/google/gopacket/ip4defrag"
	"github.com/google/gopacket/layers"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type httpReader struct {
	ident    string
	isClient bool
	bytes    chan []byte
	data     []byte
	hexdump  bool
	parent   *tcpStream
}

func (h *httpReader) Read(p []byte) (int, error) {
	ok := true
	for ok && len(h.data) == 0 {
		h.data, ok = <-h.bytes
	}
	if !ok || len(h.data) == 0 {
		return 0, io.EOF
	}

	l := copy(p, h.data)
	h.data = h.data[l:]
	return l, nil
}

func (h *httpReader) run(wg *sync.WaitGroup) {
	defer wg.Done()
	b := bufio.NewReader(h)
	for true {
		if h.isClient {
			req, err := http.ReadRequest(b)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			} else if err != nil {
				fmt.Printf("HTTP-request", "HTTP/%s Request error: %s (%v,%+v)\n", h.ident, err, err, err)
				continue
			}
			body, err := ioutil.ReadAll(req.Body)
			s := len(body)
			if err != nil {
				fmt.Printf("HTTP-request-body", "Got body err: %s\n", err)
			} else if h.hexdump {
				fmt.Printf("Body(%d/0x%x)\n%s\n", len(body), len(body), hex.Dump(body))
			}
			req.Body.Close()
			fmt.Printf("HTTP/%s Request: %s %s (body:%d)\n", h.ident, req.Method, req.URL, s)
			h.parent.Lock()
			h.parent.urls = append(h.parent.urls, req.URL.String())
			h.parent.Unlock()
		} else {
			res, err := http.ReadResponse(b, nil)
			var req string
			h.parent.Lock()
			if len(h.parent.urls) == 0 {
				req = fmt.Sprintf("<no-request-seen>")
			} else {
				req, h.parent.urls = h.parent.urls[0], h.parent.urls[1:]
			}
			h.parent.Unlock()
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			} else if err != nil {
				fmt.Printf("HTTP-response", "HTTP/%s Response error: %s (%v,%+v)\n", h.ident, err, err, err)
				continue
			}
			body, err := ioutil.ReadAll(res.Body)
			s := len(body)
			if err != nil {
				fmt.Printf("HTTP-response-body", "HTTP/%s: failed to get body(parsed len:%d): %s\n", h.ident, s, err)
			}
			if h.hexdump {
				fmt.Printf("Body(%d/0x%x)\n%s\n", len(body), len(body), hex.Dump(body))
			}
			res.Body.Close()
			sym := ","
			if res.ContentLength > 0 && res.ContentLength != int64(s) {
				sym = "!="
			}
			contentType, ok := res.Header["Content-Type"]
			if !ok {
				contentType = []string{http.DetectContentType(body)}
			}
			encoding := res.Header["Content-Encoding"]
			fmt.Printf("HTTP/%s Response: %s URL:%s (%d%s%d%s) -> %s\n", h.ident, res.Status, req, res.ContentLength, sym, s, contentType, encoding)
			if (err == nil || false) && "" != "" {
				base := url.QueryEscape(path.Base(req))
				if err != nil {
					base = "incomplete-" + base
				}
				base = path.Join("", base)
				if len(base) > 250 {
					base = base[:250] + "..."
				}
				if base == "" {
					base = path.Join("", "noname")
				}
				target := base
				n := 0
				for true {
					_, err := os.Stat(target)
					//if os.IsNotExist(err) != nil {
					if err != nil {
						break
					}
					target = fmt.Sprintf("%s-%d", base, n)
					n++
				}
				f, err := os.Create(target)
				if err != nil {
					fmt.Printf("HTTP-create", "Cannot create %s: %s\n", target, err)
					continue
				}
				var r io.Reader
				r = bytes.NewBuffer(body)
				if len(encoding) > 0 && (encoding[0] == "gzip" || encoding[0] == "deflate") {
					r, err = gzip.NewReader(r)
					if err != nil {
						fmt.Printf("HTTP-gunzip", "Failed to gzip decode: %s", err)
					}
				}
				if err == nil {
					w, err := io.Copy(f, r)
					if _, ok := r.(*gzip.Reader); ok {
						r.(*gzip.Reader).Close()
					}
					f.Close()
					if err != nil {
						fmt.Printf("HTTP-save", "%s: failed to save %s (l:%d): %s\n", h.ident, target, w, err)
					} else {
						fmt.Printf("%s: Saved %s (l:%d)\n", h.ident, target, w)
					}
				}
			}
		}
	}
}

type Context struct {
	CaptureInfo gopacket.CaptureInfo
}

func (c *Context) GetCaptureInfo() gopacket.CaptureInfo {
	return c.CaptureInfo
}

type tcpStream struct {
	tcpstate       *reassembly.TCPSimpleFSM
	fsmerr         bool
	optchecker     reassembly.TCPOptionCheck
	net, transport gopacket.Flow
	isDNS          bool
	isHTTP         bool
	reversed       bool
	client         httpReader
	server         httpReader
	urls           []string
	ident          string
	sync.Mutex
}

func (t *tcpStream) Accept(tcp *layers.TCP, ci gopacket.CaptureInfo, dir reassembly.TCPFlowDirection, nextSeq reassembly.Sequence, start *bool, ac reassembly.AssemblerContext) bool {
	// FSM
	if !t.tcpstate.CheckState(tcp, dir) {
		fmt.Printf("FSM", "%s: Packet rejected by FSM (state:%s)\n", t.ident, t.tcpstate.String())

		if !t.fsmerr {
			t.fsmerr = true
		}

		return false
	}
	// Options
	err := t.optchecker.Accept(tcp, ci, dir, nextSeq, start)
	if err != nil {
		fmt.Printf("OptionChecker", "%s: Packet rejected by OptionChecker: %s\n", t.ident, err)
		return false
	}
	// Checksum
	accept := true
	return accept
}

func (t *tcpStream) ReassembledSG(sg reassembly.ScatterGather, ac reassembly.AssemblerContext) {
	dir, _, _, skip := sg.Info()
	length, _ := sg.Lengths()
	// update stats
	sgStats := sg.Stats()
	if sgStats.OverlapBytes != 0 && sgStats.OverlapPackets == 0 {
		fmt.Printf("bytes:%d, pkts:%d\n", sgStats.OverlapBytes, sgStats.OverlapPackets)
		log.Println("Invalid overlap")
	}

	//var ident string
	if dir == reassembly.TCPDirClientToServer {
		//ident = fmt.Sprintf("%v %v(%s): ", t.net, t.transport, dir)
	} else {
		//ident = fmt.Sprintf("%v %v(%s): ", t.net.Reverse(), t.transport.Reverse(), dir)
	}
	//fmt.Printf("%s: SG reassembled packet with %d bytes (start:%v,end:%v,skip:%d,saved:%d,nb:%d,%d,overlap:%d,%d)\n", ident, length, start, end, skip, saved, sgStats.Packets, sgStats.Chunks, sgStats.OverlapBytes, sgStats.OverlapPackets)
	if skip == -1 && false {
		// this is allowed
	} else if skip != 0 {
		// Missing bytes in stream: do not even try to parse it
		return
	}
	data := sg.Fetch(length)
	if t.isDNS {
		dns := &layers.DNS{}
		var decoded []gopacket.LayerType
		if len(data) < 2 {
			if len(data) > 0 {
				sg.KeepFrom(0)
			}
			return
		}
		dnsSize := binary.BigEndian.Uint16(data[:2])
		missing := int(dnsSize) - len(data[2:])
		fmt.Printf("dnsSize: %d, missing: %d\n", dnsSize, missing)
		if missing > 0 {
			fmt.Printf("Missing some bytes: %d\n", missing)
			sg.KeepFrom(0)
			return
		}
		p := gopacket.NewDecodingLayerParser(layers.LayerTypeDNS, dns)
		err := p.DecodeLayers(data[2:], &decoded)
		if err != nil {
			fmt.Printf("DNS-parser", "Failed to decode DNS: %v\n", err)
		} else {
			fmt.Printf("DNS: %s\n", gopacket.LayerDump(dns))
		}
		if len(data) > 2+int(dnsSize) {
			sg.KeepFrom(2 + int(dnsSize))
		}
	} else if t.isHTTP {
		if length > 0 {

			if dir == reassembly.TCPDirClientToServer && !t.reversed {
				t.client.bytes <- data
			} else {
				t.server.bytes <- data
			}
		}
	}
}

func (t *tcpStream) ReassemblyComplete(ac reassembly.AssemblerContext) bool {
	//fmt.Printf("%s: Connection closed\n", t.ident)
	if t.isHTTP {
		close(t.client.bytes)
		close(t.server.bytes)
	}
	// do not remove the connection to allow last ACK
	return false
}

type tcpStreamFactory struct {
	wg     sync.WaitGroup
	doHTTP bool
}

func (factory *tcpStreamFactory) New(net, transport gopacket.Flow, tcp *layers.TCP, ac reassembly.AssemblerContext) reassembly.Stream {
	fmt.Printf("* NEW: %s %s\n", net, transport)
	fsmOptions := reassembly.TCPSimpleFSMOptions{
		SupportMissingEstablishment: false,
	}
	stream := &tcpStream{
		net:        net,
		transport:  transport,
		isDNS:      tcp.SrcPort == 53 || tcp.DstPort == 53,
		isHTTP:     (tcp.SrcPort == 80 || tcp.DstPort == 80) && factory.doHTTP,
		reversed:   tcp.SrcPort == 80,
		tcpstate:   reassembly.NewTCPSimpleFSM(fsmOptions),
		ident:      fmt.Sprintf("%s:%s", net, transport),
		optchecker: reassembly.NewTCPOptionCheck(),
	}
	if stream.isHTTP {
		stream.client = httpReader{
			bytes:    make(chan []byte),
			ident:    fmt.Sprintf("%s %s", net, transport),
			hexdump:  false,
			parent:   stream,
			isClient: true,
		}
		stream.server = httpReader{
			bytes:   make(chan []byte),
			ident:   fmt.Sprintf("%s %s", net.Reverse(), transport.Reverse()),
			hexdump: false,
			parent:  stream,
		}
		factory.wg.Add(2)
		go stream.client.run(&factory.wg)
		go stream.server.run(&factory.wg)
	}
	return stream
}

func (factory *tcpStreamFactory) WaitGoRoutines() {
	factory.wg.Wait()
}

var (
	device       string = "\\Device\\NPF_{C410B1B0-56DE-4CD5-BC7A-5A5ACAB7619F}"
	snapshot_len int32  = 65536
	promiscuous  bool   = true
	err          error
	timeout      time.Duration = -1 * time.Second
	handle       *pcap.Handle
)

func main() {
	// 打开设备
	// 打开设备进行实时捕获
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	// 构造一个数据包源
	source := gopacket.NewPacketSource(handle, handle.LinkType())
	defragger := ip4defrag.NewIPv4Defragmenter()

	streamFactory := &tcpStreamFactory{doHTTP: !false}
	streamPool := reassembly.NewStreamPool(streamFactory)
	assembler := reassembly.NewAssembler(streamPool)

	// 读取包
	for packet := range source.Packets() {
		//fmt.Println("=========================")
		//fmt.Println(packet)
		// Process packet here
		ip4Layer := packet.Layer(layers.LayerTypeIPv4)
		if ip4Layer == nil {
			continue
		}
		ip4 := ip4Layer.(*layers.IPv4)
		l := ip4.Length
		newip4, err := defragger.DefragIPv4(ip4)
		if err != nil {
			log.Fatalln("Error while de-fragmenting", err)
		} else if newip4 == nil {
			fmt.Println("Fragment...\n")
			continue // packet fragment, we don't have whole packet yet.
		}
		if newip4.Length != l {
			//stats.ipdefrag++
			fmt.Printf("Decoding re-assembled packet: %s\n", newip4.NextLayerType())
			pb, ok := packet.(gopacket.PacketBuilder)
			if !ok {
				log.Panicln("Not a PacketBuilder")
			}
			nextDecoder := newip4.NextLayerType()
			nextDecoder.Decode(newip4.Payload, pb)
		}

		tcp := packet.Layer(layers.LayerTypeTCP)
		if tcp != nil {
			tcp := tcp.(*layers.TCP)
			c := Context{
				CaptureInfo: packet.Metadata().CaptureInfo,
			}
			assembler.AssembleWithContext(packet.NetworkLayer().NetworkFlow(), tcp, &c)
		}

	}

}
