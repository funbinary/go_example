package main

import (
	"gopkg.in/gomail.v2"
	"strconv"
)

func main() {
	userName := "xxxx@163.com"
	authCode := "xxxx"
	host := "localhost"
	port := "2500"
	err := SendMail(userName, authCode, host, port, "byhu@1.com", "糊了", "你昨天干嘛了", "我昨天发呆了一天")
	if err != nil {
		panic(err)
	}
}

func SendMail(userName, authCode, host, portStr, mailTo, sendName string, subject, body string) error {
	port, _ := strconv.Atoi(portStr)
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(userName, sendName))
	m.SetHeader("To", mailTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(host, port, userName, authCode)
	err := d.DialAndSend(m)
	return err
}
