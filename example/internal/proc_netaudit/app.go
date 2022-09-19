package proc_netaudit

import (
	"github.com/funbinary/go_example/example/internal/proc_netaudit/decoder"
	"github.com/funbinary/go_example/example/internal/proc_netaudit/hook"
	"github.com/funbinary/go_example/example/internal/proc_netaudit/logic"
	"github.com/funbinary/go_example/example/internal/proc_netaudit/source"
	"github.com/funbinary/go_example/example/internal/proc_netaudit/stream"

	"fmt"
	"github.com/funbinary/go_example/example/internal/proc_netaudit/config"
	"github.com/funbinary/go_example/pkg/errors"
	"github.com/funbinary/go_example/pkg/reassembly"
	"github.com/google/gopacket"
	"github.com/google/gopacket/ip4defrag"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"

	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/funbinary/go_example/pkg/bfile"
	log "github.com/funbinary/go_example/pkg/blog"
	_ "net/http/pprof"
)

type CaptureInfo struct {
	gopacket.CaptureInfo
	Data []byte
}

func Run() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
	runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪

	var level string
	if config.Config.LogLevel < 1 {
		level = "Debug"
	} else if config.Config.LogLevel == 2 {
		level = "Info"
	} else if config.Config.LogLevel == 3 {
		level = "Warn"
	}

	// 初始化日志
	InitLog(config.Config.LogPath, level)

	// 启动pprof
	go func() {
		// 启动一个 http server，注意 pprof 相关的 handler 已经自动注册过了
		if err := http.ListenAndServe(":6069", nil); err != nil {
			log.Fatalf("%v", err)
		}
		os.Exit(0)
	}()
	log.Infof("=====================================NetAudit Start============================================================")
	log.Infof("%v", config.Config)

	// 初始化pcap
	handler := source.NewPcapSource(config.Config.Interface,
		source.WithBufferSize(config.Config.BufferSize),
		source.WithTimeout(config.Config.TimeOut),
		source.WithSnaplen(config.Config.Snaplen),
		source.WithPromisc(config.Config.Promisc))
	if handler == nil {
		log.Fatal("启动包捕获器失败")
	}
	defer handler.Close()
	//"host !192.168.3.12 and port !22"
	handler.SetBPFFilter(config.Config.BPF)

	var signalChan chan os.Signal
	signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	// 注册Hook
	hook.SetHook(logic.NewHook())

	// 包解析器
	dataSource := gopacket.NewPacketSource(handler, handler.LinkType())
	// 创建分片重组器
	defrag := ip4defrag.NewIPv4Defragmenter()

	// 创建TCP重组器
	streamFactory := &stream.TCPStreamFactory{}
	streamPool := reassembly.NewStreamPool(streamFactory)
	assembler := reassembly.NewAssembler(streamPool)

	// 开启dump pcap包
	var ch chan CaptureInfo
	if config.Config.EnablePcapDump {
		path := bfile.Join(config.Config.PcapDumpPath, "capture_"+time.Now().Format("20060102150405")+".pcap")
		f, err := bfile.Create(path)
		if err != nil {
			log.Errorf("open dump pcap packet file:%s fail:err=%v", path, err)
		}
		w := pcapgo.NewWriter(f)
		w.WriteFileHeader(uint32(65536), layers.LinkTypeEthernet)
		defer f.Close()
		ch = make(chan CaptureInfo)

		go func() {
			for {
				select {
				case capture := <-ch:
					w.WritePacket(capture.CaptureInfo, capture.Data)
				}
			}
		}()
	}
	// 开启读取包协程，并从中获取已解析完传输层的包通道
	packetch := dataSource.Packets()
	ticker := time.NewTicker(5 * time.Second)
	count := 0
	for {

		select {
		case packet := <-packetch:
			count++
			if config.Config.EnablePcapDump {
				capture := CaptureInfo{
					CaptureInfo: packet.Metadata().CaptureInfo,
					Data:        packet.Data(),
				}
				ch <- capture
			}
			// 不包含IPV4不处理
			ip4Layer := packet.Layer(layers.LayerTypeIPv4)
			if ip4Layer == nil {
				continue
			}
			ip4 := ip4Layer.(*layers.IPv4)
			l := ip4.Length
			// 丢到分片重组器中
			newip4, err := defrag.DefragIPv4(ip4)
			if err != nil {
				log.Fatalf("Error while de-fragmenting", err)
			} else if newip4 == nil {
				// 存在分片
				continue
			}
			if newip4.Length != l {
				log.Errorf("Decoding re-assembled packet: %s\n", newip4.NextLayerType())
				pb, ok := packet.(gopacket.PacketBuilder)
				if !ok {
					log.Errorf("Not a PacketBuilder")
				}
				nextDecoder := newip4.NextLayerType()
				nextDecoder.Decode(newip4.Payload, pb)
			}

			// udp
			udp := packet.Layer(layers.LayerTypeUDP)
			if udp != nil {
				defer func() {
					err := recover()
					log.Errorf("%+v", err)
				}()
				udp := udp.(*layers.UDP)
				id := fmt.Sprintf("%s_%s_%s_%s_%s", ip4.NetworkFlow().Src(), udp.TransportFlow().Src(), ip4.NetworkFlow().Dst(), udp.TransportFlow().Dst(), "udp")
				if err := decoder.Decode(id, udp.Payload, udp); err != nil {
					if !errors.Is(err, decoder.ErrDecoderUnknown) {
						log.Errorf("%+v", err)
					}
				}
			}

			// tcp
			tcp := packet.Layer(layers.LayerTypeTCP)
			if tcp != nil {
				tcp := tcp.(*layers.TCP)

				err := tcp.SetNetworkLayerForChecksum(packet.NetworkLayer())
				if err != nil {
					log.Fatalf("Failed to set network layer for checksum: %s\n", err)
				}

				ctx := &stream.Context{
					CaptureInfo: packet.Metadata().CaptureInfo,
				}
				//stats.totalsz += len(tcp.Payload)
				log.Debugf("-----------start----------------")
				log.Debugf("tcp sport:%v tcp dport: %v seq:%v payloadlen%v", tcp.SrcPort, tcp.DstPort, tcp.Seq, len(tcp.Payload))
				assembler.AssembleWithContext(packet.NetworkLayer().NetworkFlow(), tcp, ctx)
				log.Debugf("-----------end----------------")
				if count%1000 == 0 {
					ref := packet.Metadata().CaptureInfo.Timestamp
					flushed, closed := assembler.FlushWithOptions(reassembly.FlushOptions{T: ref.Add(-time.Minute * 3), TC: ref.Add(-time.Minute * 5)})
					log.Debugf("Forced flush: %d flushed, %d closed (%s)", flushed, closed, ref)
				}
			}

		case <-signalChan:

			closed := assembler.FlushAll()
			log.Infof("Final flush: %d closed", closed)
			//streamFactory.WaitGoRoutines()
			//log.Infof("%s\n", assembler.Dump())
			log.Infof("接受到信号退出程序...")
			streamFactory.WaitGoRoutines()
			os.Exit(1)
		case <-ticker.C:
			log.Infof("packet chan length:%d", len(packetch))
			assembler.FlushCloseOlderThan(time.Now().Add(time.Minute * -2))

			//ref := packet.Metadata().CaptureInfo.Timestamp
			//flushed, closed := assembler.FlushWithOptions(reassembly.FlushOptions{T: ref.Add(-timeout), TC: ref.Add(-closeTimeout)})
			//Debug("Forced flush: %d flushed, %d closed (%s)", flushed, closed, ref)
		}

	}

}
