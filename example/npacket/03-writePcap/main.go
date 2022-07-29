package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

var (
	device      string = "\\Device\\NPF_{C410B1B0-56DE-4CD5-BC7A-5A5ACAB7619F}"
	snaplen     int32  = 65536
	promiscuous bool   = true
	err         error
	timeout     time.Duration = -1 * time.Second
	handle      *pcap.Handle
	packetCount int = 0
)

func main() {
	// 打开一个输出的文件句柄
	f, _ := os.Create("test.pcap")
	// 创建一个writer对象
	w := pcapgo.NewWriter(f)
	// 写入文件头，必须在调用前调用
	w.WriteFileHeader(uint32(snaplen), layers.LinkTypeEthernet)
	defer f.Close()
	// 打开一个实时捕获设备
	handle, err = pcap.OpenLive(device, snaplen, promiscuous, timeout)
	if err != nil {
		fmt.Printf("Error opening device %s: %v", device, err)
		os.Exit(1)
	}
	defer handle.Close()

	// 创建包数据源
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	// 读取包
	for packet := range packetSource.Packets() {
		// 打印包
		fmt.Println(packet)
		// 写入包
		w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		packetCount++

		// Only capture 100 and then stop
		if packetCount > 100 {
			break
		}
	}
}
