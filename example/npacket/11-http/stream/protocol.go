package stream

import (
	"bufio"
	"bytes"
	"net/textproto"
	"strings"
)

type HTTPDirect int

// Value are not really useful
const (
	HTTPDirectUnknown  HTTPDirect = -1
	HTTPDirectRequest  HTTPDirect = 0
	HTTPDirectResopnse HTTPDirect = 1
)

func (dir HTTPDirect) String() string {
	switch dir {
	case HTTPDirectUnknown:
		return "unknown"
	case HTTPDirectRequest:
		return "http request"
	case HTTPDirectResopnse:
		return "http response"
	}
	return ""
}

func detectHttp(data []byte) (bool, HTTPDirect) {

	ishttp, _ := isResponse(data)
	if ishttp {
		return true, HTTPDirectRequest
	}

	ishttp, _ = isRequest(data)
	if ishttp {
		return true, HTTPDirectResopnse //request
	}

	return false, HTTPDirectUnknown
}

func isResponse(data []byte) (bool, string) {
	buf := bytes.NewBuffer(data)
	reader := bufio.NewReader(buf)
	tp := textproto.NewReader(reader)

	firstLine, _ := tp.ReadLine()
	return strings.HasPrefix(strings.TrimSpace(firstLine), "HTTP/"), firstLine
}
func isRequest(data []byte) (bool, string) {
	buf := bytes.NewBuffer(data)
	reader := bufio.NewReader(buf)
	tp := textproto.NewReader(reader)

	firstLine, _ := tp.ReadLine()
	arr := strings.Split(firstLine, " ")

	switch strings.TrimSpace(arr[0]) {
	case "GET", "HEAD", "POST", "PUT", "DELETE", "TRACE", "CONNECT":
		return true, firstLine
	default:
		return false, firstLine
	}
}
