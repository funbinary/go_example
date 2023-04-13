/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/binary"
	"fmt"
	"github.com/funbinary/go_example/pkg/encoding/bjson"
	"github.com/funbinary/go_example/pkg/errors"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const ntpEpochOffset = 2208988800

type packet struct {
	Settings       uint8
	Stratum        uint8
	Poll           int8
	Precision      int8
	RootDelay      uint32
	RootDispersion uint32
	ReferenceID    uint32
	RefTimeSec     uint32
	RefTimeFrac    uint32
	OrigTimeSec    uint32
	OrigTimeFrac   uint32
	RxTimeSec      uint32
	RxTimeFrac     uint32
	TxTimeSec      uint32
	TxTimeFrac     uint32
}

type Result struct {
	Error string `json:"error"`
	Time  string `json:"time"`
}

var host string
var port string
var timeout int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ntp",
	Short: "获取NTP时间",
	Long:  `获取NTP时间.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var timeLayoutStr = "2006-01-02 15:04:05"
		var err error
		var t time.Time
		defer func() {
			r := Result{
				Error: "",
				Time:  t.Format(timeLayoutStr),
			}
			if err != nil {
				r.Error = err.Error()
			}
			fmt.Print(bjson.Marshal(r))
		}()
		t, err = getremotetime(host, port, timeout)
		if err != nil {
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().StringVarP(&host, "host", "", "ntp.ntsc.ac.cn", "NTP服务器地址")
	rootCmd.Flags().StringVarP(&port, "port", "p", "123", "NTP服务器端口")
	rootCmd.Flags().IntVarP(&timeout, "timeout", "t", 2, "超时时间")
}

func getremotetime(host, port string, timeout int) (t time.Time, err error) {
	serverAddr := host + ":" + port

	var conn net.Conn
	conn, err = net.Dial("udp", serverAddr)
	if err != nil {
		return t, errors.Errorf("failed  connect: %v", err)
	}
	defer conn.Close()
	if err = conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
		return t, errors.Errorf("failed to set deadline: %v", err)
	}

	req := &packet{Settings: 0x1B}

	if err = binary.Write(conn, binary.BigEndian, req); err != nil {
		return t, errors.Errorf("failed to send request: %v", err)
	}

	rsp := &packet{}
	if err = binary.Read(conn, binary.BigEndian, rsp); err != nil {
		return t, errors.Errorf("failed to read server response: %v", err)
	}

	secs := float64(rsp.TxTimeSec) - ntpEpochOffset
	nanos := (int64(rsp.TxTimeFrac) * 1e9) >> 32

	t = time.Unix(int64(secs), nanos)

	return t, nil
}
