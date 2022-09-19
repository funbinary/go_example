package main

import (
	"fmt"
	"github.com/funbinary/go_example/pkg/bshell"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	baseargs := []string{
		"1h", "tcpdump", "host", "192.168.3.247", "and", "port", "12011", "-w",
	}
	go func() {
		message := make(chan os.Signal, 1)

		signal.Notify(message,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGQUIT,
			syscall.SIGILL,
			syscall.SIGTRAP,
			syscall.SIGABRT,
			syscall.SIGBUS,
			syscall.SIGFPE,
			syscall.SIGKILL,
			syscall.SIGSEGV,
			syscall.SIGPIPE,
			syscall.SIGALRM,
			syscall.SIGTERM, os.Interrupt)
		<-message
		fmt.Println("程序即将退出")
		result, _ := bshell.ShellExec("ps -ef | grep tcpdump | grep -v grep | awk 'NR==1 {print $2}'")
		fmt.Println(result)
		bshell.ShellExec("kill" + result)
	}()
	go func() {
		t := time.NewTicker(6 * time.Hour)
		for {
			select {
			case <-t.C:
				fmt.Println("清理过期文件")
			}
		}
	}()
	for {
		curArgs := baseargs
		curDate := time.Now().Format("20060102150402")
		file := curDate + ".pcap"
		curArgs = append(curArgs, file)
		cmd := exec.Command("timeout", curArgs...)
		fmt.Printf("执行命令:%v\n", cmd.Args)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("执行命令:%v 出错:%+v\n", cmd.Args, err)
			break
		}
	}
}
