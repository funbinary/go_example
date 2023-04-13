package main

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	// 创建一个message Printer对象，用于格式化输出
	p := message.NewPrinter(language.English)

	// 输出一个多语言字符串
	p.Printf("你好，%s！\n", "世界")
	// 设置语言为中文
	p = message.NewPrinter(language.Chinese)

	// 输出一个多语言字符串
	p.Printf("你好，%s！\n", "世界")

}
