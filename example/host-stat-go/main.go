//go:build linux

package main

import (
	"fmt"
	"github.com/funbinary/go_example/pkg/encoding/bjson"
	hoststat "github.com/likexian/host-stat-go"
	"github.com/shirou/gopsutil/cpu"
	"time"
)

type Stat struct {
	MemStat hoststat.MemStat  `json:"memStat"`
	CpuStat hoststat.CPUStat  `json:"cpuStat"`
	IoStat  []hoststat.IOStat `json:"ioStat"`
}

func main() {
	var s Stat
	var err error
	s.MemStat, err = hoststat.GetMemStat()
	if err != nil {
		panic(err)
	}

	s.CpuStat, err = hoststat.GetCPUStat()
	if err != nil {
		panic(err)
	}

	s.IoStat, err = hoststat.GetIOStat()
	if err != nil {
		panic(err)
	}
	fmt.Println(bjson.Marshal(s))
	fmt.Println(hoststat.GetNetStat())
	c, _ := cpu.Info()
	fmt.Println("cpu信息:", c)
	/*用户CPU时间/系统CPU时间/空闲时间。。。等等
	  用户CPU时间：就是用户的进程获得了CPU资源以后，在用户态执行的时间。
	  系统CPU时间：用户进程获得了CPU资源以后，在内核态的执行时间。
	*/
	c1, _ := cpu.Times(false)
	fmt.Println("cpu1:", c1)

	// CPU使用率，每秒刷新一次
	for {
		c2, _ := cpu.Percent(time.Duration(time.Second), true)
		fmt.Println(c2)
	}
}
