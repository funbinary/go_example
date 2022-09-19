package main

import (
	"fmt"
	"github.com/funbinary/go_example/example/npacket/11-http/stream"
	"github.com/funbinary/go_example/pkg/bfile"
	"github.com/google/gopacket"
	"github.com/google/gopacket/ip4defrag"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/reassembly"
	"log"
	"os"
	"os/signal"
	"time"
)

var exit = false

type Context struct {
	CaptureInfo gopacket.CaptureInfo
}

func (c *Context) GetCaptureInfo() gopacket.CaptureInfo {
	return c.CaptureInfo
}

func main() {
	//fp, err := bfile.Create(bfile.Join(bfile.SelfDir(), "cpuprofile.out"))
	//
	//if err != nil {
	//	log.Fatalf("could not open cpu profile file ")
	//}
	//pprof.StartCPUProfile(fp)
	//defer func() {
	//	pprof.StopCPUProfile()
	//	fp.Close()
	//
	//}()

	//f, _ := bfile.Create(bfile.Join(bfile.SelfDir(), "test.pcap"))
	//w := pcapgo.NewWriter(f)
	//w.WriteFileHeader(uint32(65536), layers.LinkTypeEthernet)
	//defer f.Close()

	var handler *pcap.Handle
	handler, err := pcap.OpenOffline(bfile.Join(bfile.SelfDir(), "test.pcap"))
	if err != nil {
		panic(err)
	}
	defer handler.Close()
	//dev := "\\Device\\NPF_{C410B1B0-56DE-4CD5-BC7A-5A5ACAB7619F}"
	//dev := "br0"
	//ih, err := pcap.NewInactiveHandle(dev)
	//if err != nil {
	//	panic(err)
	//}
	//defer ih.CleanUp()
	//if err := ih.SetTimeout(-1 * time.Second); err != nil {
	//	panic(err)
	//}
	//if err := ih.SetPromisc(true); err != nil {
	//	panic(err)
	//}
	//if err := ih.SetSnapLen(65536); err != nil {
	//	panic(err)
	//}
	//if handler, err = ih.Activate(); err != nil {
	//	panic(err)
	//}
	//defer handler.Close()

	source := gopacket.NewPacketSource(handler, handler.LinkType())
	//source.NoCopy = true

	count := 0
	bytes := int64(0)
	//start := time.Now()

	defragger := ip4defrag.NewIPv4Defragmenter()
	streamFactory := &stream.TCPStreamFactory{}
	streamPool := reassembly.NewStreamPool(streamFactory)
	assembler := reassembly.NewAssembler(streamPool)

	var signalChan chan os.Signal
	signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	go HandlerSig(signalChan)

	for packet := range source.Packets() {
		fmt.Println()
		//w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		count++
		data := packet.Data()
		bytes += int64(len(data))

		//fmt.Printf("Packet content (%d/0x%x)\n%s\n", len(data), len(data), hex.Dump(data))

		// defrag the IPv4 packet if required

		//IPV4
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
			fmt.Println("Fragment...")
			continue // packet fragment, we don't have whole packet yet.
		}
		if newip4.Length != l {

			fmt.Println("Decoding re-assembled packet: %s\n", newip4.NextLayerType())
			pb, ok := packet.(gopacket.PacketBuilder)
			if !ok {
				panic("Not a PacketBuilder")
			}
			nextDecoder := newip4.NextLayerType()
			nextDecoder.Decode(newip4.Payload, pb)
		}

		// tcp
		tcp := packet.Layer(layers.LayerTypeTCP)
		if tcp != nil {
			tcp := tcp.(*layers.TCP)

			err := tcp.SetNetworkLayerForChecksum(packet.NetworkLayer())
			if err != nil {
				log.Fatalf("Failed to set network layer for checksum: %s\n", err)
			}

			c := Context{
				CaptureInfo: packet.Metadata().CaptureInfo,
			}
			//stats.totalsz += len(tcp.Payload)

			assembler.AssembleWithContext(packet.NetworkLayer().NetworkFlow(), tcp, &c)
		}

		if count%1000 == 0 {
			ref := packet.Metadata().CaptureInfo.Timestamp
			flushed, closed := assembler.FlushWithOptions(reassembly.FlushOptions{T: ref.Add(-3 * time.Minute), TC: ref.Add(-time.Minute * 5)})
			fmt.Println("Forced flush: %d flushed, %d closed (%s)", flushed, closed, ref)
		}

		if exit {
			break
		}
	}

	closed := assembler.FlushAll()
	fmt.Println("Final flush: %d closed", closed)
	streamPool.Dump()
	streamFactory.WaitGoRoutines()
	fmt.Printf("%s\n", assembler.Dump())

}

func HandlerSig(ch chan os.Signal) {
	<-ch
	if !exit {
		exit = true
	}
	time.Sleep(3 * time.Second)
	os.Exit(1)
}
