/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var server string
var room string
var num int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "webrtc",
	Short: "基于chrome实现webrtc压测",
	Long:  `基于chrome实现webrtc压测`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if room == "" {
			fmt.Println("房间不能为空")
			return
		}
		var wg sync.WaitGroup
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT)

		ctx, cancel := context.WithCancel(context.Background())

		for i := 0; i < num; i++ {
			time.Sleep(time.Second)
			wg.Add(1)
			go func() {
				webrtcmock(ctx)
				wg.Done()
			}()

		}

		<-sigs
		fmt.Println("接收到 Ctrl+C 信号，程序退出中.......")
		cancel()
		wg.Wait()

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
	rootCmd.Flags().StringVarP(&room, "room", "r", "", "房价号")
	rootCmd.Flags().StringVarP(&server, "server", "s", "http://192.168.1.61", "服务器地址")
	rootCmd.Flags().IntVarP(&num, "num", "n", 5, "模拟的数量")
}

func webrtcmock(ctx context.Context) {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("unsafely-treat-insecure-origin-as-secure", server),
	}

	execCtx, _ := chromedp.NewExecAllocator(ctx, opts...)

	ctx2, _ := chromedp.NewContext(execCtx)

	err := chromedp.Run(ctx2, chromedp.Navigate(server),
		chromedp.Sleep(time.Second),
		chromedp.SendKeys("#app > div > section > main > div > div:nth-child(2) > div > div:nth-child(2) > input", room),
		chromedp.Sleep(time.Second),
		chromedp.Click("#app > div > section > main > div > div:nth-child(2) > div > button.el-button.el-button--primary > span"),
		chromedp.Sleep(5*time.Second),
		chromedp.Evaluate(`const videos = document.getElementsByTagName('video');

// 循环遍历所有的video元素，并停止它们的播放
for (let i = 0; i < videos.length; i++) {
  const video = videos[i];
  if (!video.paused) {
    video.pause();
  }
}
`, nil),
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = chromedp.Run(ctx2, chromedp.Sleep(24*time.Hour))
	if err != nil {
		fmt.Println(err)
	}

}
