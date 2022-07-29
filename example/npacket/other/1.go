//创建和发送packet
//下面这个例子做了几个事情。首先，它将演示如何使用网络设备发送原始字节。这样，您就可以像串行连接(serial connection)一样使用它来发送数据。这对于真正的低层的数据传输很有用，但是如果你想与一个应用程序交互，你可能想建立硬件和软件都能识别的包。
//
//接下来，它将演示如何使用ethernet、IP和TCP层创建数据包。所有的东西都是默认的和空的，所以它实际上不做任何事情。
//
//为了完成它，我们创建了另一个数据包，但实际上为ethernet层填充了一些MAC地址，为IPv4填充了一些IP地址，为TCP层填充了一些端口号。您应该看到如何用它伪造数据包和模拟设备。
//
//TCP层结构具有可读取或设置的SYN, FIN, and ACK 布尔标志。这有利于控制和模糊TCP握手、会话和端口扫描。
//
//pcap库提供了一个发送字节的简单方法，但是gopacket中的layers包帮助我们为各个层创建字节结构。

package main

import (
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	device             = "eth0"
	snapshot_len int32 = 1024
	promiscuous        = false
	err          error
	timeout      = 30 * time.Second
	handle       *pcap.Handle
	buffer       gopacket.SerializeBuffer
	options      gopacket.SerializeOptions
)

func main() {
	// Open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	// Send raw bytes over wire
	rawBytes := []byte{10, 20, 30}
	err = handle.WritePacketData(rawBytes)
	if err != nil {
		log.Fatal(err)
	}
	// Create a properly formed packet, just with
	// empty details. Should fill out MAC addresses,
	// IP addresses, etc.
	buffer = gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buffer, options,
		&layers.Ethernet{},
		&layers.IPv4{},
		&layers.TCP{},
		gopacket.Payload(rawBytes),
	)
	outgoingPacket := buffer.Bytes()
	// Send our packet
	err = handle.WritePacketData(outgoingPacket)
	if err != nil {
		log.Fatal(err)
	}
	// This time lets fill out some information
	ipLayer := &layers.IPv4{
		SrcIP: net.IP{127, 0, 0, 1},
		DstIP: net.IP{8, 8, 8, 8},
	}
	ethernetLayer := &layers.Ethernet{
		SrcMAC: net.HardwareAddr{0xFF, 0xAA, 0xFA, 0xAA, 0xFF, 0xAA},
		DstMAC: net.HardwareAddr{0xBD, 0xBD, 0xBD, 0xBD, 0xBD, 0xBD},
	}
	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(4321),
		DstPort: layers.TCPPort(80),
	}
	// And create the packet with the layers
	buffer = gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buffer, options,
		ethernetLayer,
		ipLayer,
		tcpLayer,
		gopacket.Payload(rawBytes),
	)
	outgoingPacket = buffer.Bytes()
}
