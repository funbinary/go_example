package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"sync"
	"time"

	"github.com/shirou/gopsutil/net"
)

func main() {
	var wg sync.WaitGroup
	// Get initial IO counters
	var cpuPercent float64
	wg.Add(1)
	go func() {
		defer wg.Done()
		cpuPercents, err := cpu.Percent(time.Second, false)
		if err != nil || len(cpuPercents) < 1 {
			return
		}
		cpuPercent = cpuPercents[0]
	}()

	var ioSendPercent, ioReceivePercent uint64
	var rx, tx uint64

	wg.Add(1)
	go func() {
		defer wg.Done()

		initialCounters, err := net.IOCounters(false)
		if err != nil {
			return
		}

		// Sleep for one second
		time.Sleep(time.Second)

		// Get final IO counters
		finalCounters, err := net.IOCounters(false)
		if err != nil {
			return
		}

		// Calculate total bytes sent and received
		var totalBytesSent, totalBytesReceived uint64
		for i := 0; i < len(initialCounters); i++ {
			totalBytesSent += finalCounters[i].BytesSent - initialCounters[i].BytesSent
			totalBytesReceived += finalCounters[i].BytesRecv - initialCounters[i].BytesRecv
		}
		ioSendPercent = finalCounters[0].BytesSent - initialCounters[0].BytesSent
		ioReceivePercent = finalCounters[0].BytesRecv - initialCounters[0].BytesRecv

		rx = finalCounters[0].BytesSent
		tx = finalCounters[0].BytesRecv
	}()

	memStat, err := mem.VirtualMemory()
	if err != nil || memStat == nil {
		return
	}
	memPercent := memStat.Available / 1024 / 1024 / 16384

	wg.Wait()

	fmt.Println("cpuPercent:", cpuPercent)
	fmt.Println("memTotal:", memStat.Total/1024/1024)
	fmt.Println("memAvailable:", memStat.Available/1024/1024)
	//fmt.Println(dump.Format(memStat))
	fmt.Println("memPercent:", memPercent)
	fmt.Println("iosend:", ioSendPercent)
	fmt.Println("ioReceivePercent:", ioReceivePercent)
	fmt.Println("rx:", rx)
	fmt.Println("tx:", tx)
}
