package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	rtmp "github.com/zhangpeihao/gortmp"
	"github.com/zhangpeihao/log"
	"os"
	"sync/atomic"
	"time"
)

var url, room, display string
var num int
var curNum int32

var rootCmd = &cobra.Command{
	Use:   "rtmp_benchmark",
	Short: "rtmp_benchmark",
	Long:  `rtmp_benchmark`,

	Run: func(cmd *cobra.Command, args []string) {
		l := log.NewLogger(".", "player", nil, 60, 3600*24, true)
		rtmp.InitLogger(l)
		defer l.Close()
		url = fmt.Sprintf("%s/%s", url, room)
		fmt.Println(url, display)
		for i := 0; i < num-1; i++ {
			go Run()
		}
		go func() {
			time.Sleep(2 * time.Second)
			fmt.Println("当前运行拉流数:", atomic.LoadInt32(&curNum))
		}()

		Run()
	},
}

func Run() {
	atomic.AddInt32(&curNum, 1)
	defer atomic.AddInt32(&curNum, -1)

	rtmpHandler := NewRtmpOutboundConnHandler()

	var obConn rtmp.OutboundConn
	var err error

	obConn, err = rtmp.Dial(url, rtmpHandler, 100)
	if err != nil {
		fmt.Println("Dial error", err)
		os.Exit(-1)
	}
	defer obConn.Close()

	err = obConn.Connect()
	if err != nil {
		fmt.Printf("Connect error: %s", err.Error())
		os.Exit(-1)
	}
	for {
		select {
		case stream := <-rtmpHandler.StreamChan:
			// Play
			fmt.Println("play")
			err = stream.Play(display, nil, nil, nil)
			if err != nil {
				fmt.Printf("Play error: %s", err.Error())
				os.Exit(-1)
			}
		}
	}
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rtmp_benchmark.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&url, "url", "u", "rtmp://192.168.3.250", "rtmp流地址")
	rootCmd.Flags().StringVarP(&room, "room", "r", "643606265f79c3425d097473", "房间地址")
	rootCmd.Flags().StringVarP(&display, "display", "d", "98b43e937320", "display")
	rootCmd.Flags().IntVarP(&num, "num", "n", 10, "运行数量")

}
