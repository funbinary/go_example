package main

import (
	"fmt"
	"github.com/funbinary/go_example/example/npacket/08-ftpCapture/layer"
	"io"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	device       string = "\\Device\\NPF_{48641DC5-6BD6-4752-9CA4-5C9706829705}"
	snapshot_len int32  = 65536
	promiscuous  bool   = true
	err          error
	timeout      time.Duration = -1 * time.Second
	handle       *pcap.Handle
	// Will reuse these for each packet
	ethLayer layers.Ethernet
	ipLayer  layers.IPv4
	tcpLayer layers.TCP
	ftpLayer layer.FTPLayer
)

func main() {
	// Open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	//var filter string = "tcp and port 10021"
	//err = handle.SetBPFFilter(filter)
	//if err != nil {
	//	log.Fatal(err)
	//}

	layers.RegisterTCPPortLayerType(layers.TCPPort(21), layer.LayerTypeFTP)
	dlp := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet)
	dlp.SetDecodingLayerContainer(gopacket.DecodingLayerSparse(nil))
	//var eth layers.Ethernet
	dlp.AddDecodingLayer(&ethLayer)
	dlp.AddDecodingLayer(&ipLayer)
	dlp.AddDecodingLayer(&tcpLayer)
	dlp.AddDecodingLayer(&ftpLayer)

	// ... 添加层并照常使用 DecodingLayerParser...
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for {
		packet, err := packetSource.NextPacket()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("Error:", err)
			continue
		}
		foundLayerTypes := []gopacket.LayerType{}
		err = dlp.DecodeLayers(packet.Data(), &foundLayerTypes)
		if err != nil {
			fmt.Println("Trouble decoding layers: ", err)
		}
		for _, layerType := range foundLayerTypes {
			//if layerType == layers.LayerTypeIPv4 {
			//	fmt.Println("IPv4: ", ipLayer.SrcIP, "->", ipLayer.DstIP)
			//}
			//if layerType == layers.LayerTypeTCP {
			//	fmt.Println("TCP Port: ", tcpLayer.SrcPort, "->", tcpLayer.DstPort)
			//	fmt.Println("TCP SYN:", tcpLayer.SYN, " | ACK:", tcpLayer.ACK)
			//}
			if layerType == layer.LayerTypeFTP {
				fmt.Println(ftpLayer.Command)
			}

		}
	}

}
