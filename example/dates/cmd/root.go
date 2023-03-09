/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dates",
	Short: "设置时间",
	Long: `For example:
dates "20210308111553"		
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		//cst := time.FixedZone("CST", 8*3600) //创建CST时区

		var format string
		for _, v := range args {
			format += v
		}
		fmt.Println("input:", format)
		layout := "20060102150405"
		cst, _ := time.LoadLocation("Asia/Shanghai") //指定CST时区
		t, err := time.ParseInLocation(layout, format, cst)
		if err != nil {
			panic(err)
		}
		fmt.Println(t)
		unixTime := t.Unix()
		fmt.Println("unixTime:", unixTime)

		//// 将时间转换为Unix时间戳
		tv := syscall.Timeval{
			Sec:  unixTime,
			Usec: 0,
		}

		// 设置系统时间
		if err := syscall.Settimeofday(&tv); err != nil {
			fmt.Println("Failed to set system time:", err)
			return
		}

		// 构造系统调用需要的数据结构
		//tv := syscall.NsecToTimeval(unixTime)
		//tvptr := &tv
		//
		//// 调用 settimeofday 系统调用
		//err = syscall.Settimeofday(tvptr)
		//if err != nil {
		//	panic(err)
		//}
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dates.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
