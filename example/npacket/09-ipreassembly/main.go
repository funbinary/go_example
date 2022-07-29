package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket/ip4defrag"
	"github.com/google/gopacket/layers"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

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

		udpLayer := packet.Layer(layers.LayerTypeUDP)
		if udpLayer != nil {
			udp := udpLayer.(*layers.UDP)
			if udp.DstPort == 30000 {
				fmt.Println(udp.Length)
				fmt.Println(string(udp.LayerPayload()))
				fmt.Println("=========================")
			}

		}

	}

}
