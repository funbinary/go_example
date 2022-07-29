package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	device       string = " \\Device\\NPF_{C410B1B0-56DE-4CD5-BC7A-5A5ACAB7619F}\n"
	snapshot_len int32  = 65536
	promiscuous  bool   = true
	err          error
	timeout      time.Duration = -1 * time.Second
	handle       *pcap.Handle
)

func main() {
	// 打开设备进行实时捕获
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	// 构造一个数据包源
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	// 读取包
	for packet := range packetSource.Packets() {
		// Process packet here
		fmt.Println(packet)
	}
}
