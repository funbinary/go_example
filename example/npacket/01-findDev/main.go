package main

import (
	"fmt"
	"log"

	"github.com/funbinary/go_example/pkg/bfile"

	"github.com/google/gopacket/pcap"
)

func main() {
	// 得到所有的(网络)设备
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	// 打印设备信息
	fmt.Println("Devices found:")
	for _, device := range devices {
		if len(device.Addresses) <= 0 {
			continue
		}
		fmt.Println("\nName: ", device.Name)
		fmt.Println("Description: ", device.Description)
		fmt.Println("Devices flag: ", device.Flags)
		for _, address := range device.Addresses {
			fmt.Println("- IP address: ", address.IP)
			fmt.Println("- Subnet mask: ", address.Netmask)
			fmt.Println("- Broadaddr:  ", address.Broadaddr)
			fmt.Println("- P2P:  ", address.P2P)
		}
	}
	fmt.Println(bfile.SelfDir())
}
