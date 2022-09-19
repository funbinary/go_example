package stream

import (
	"bufio"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

type httpReader struct {
	ident     string
	isClient  bool
	bytes     chan []byte
	timestamp chan int64
	data      []byte

	srcip   string
	dstip   string
	srcport string
	dstport string

	httpstart int
}

func (self *httpReader) Read(b []byte) (int, error) {
	ok := true
	for ok && len(self.data) == 0 {
		self.data, ok = <-self.bytes
	}
	if !ok || len(self.data) == 0 {
		return 0, io.EOF
	}

	ishttp, _ := detectHttp(self.data)

	if ishttp {

		switch self.httpstart {
		case 0: // run read,only copy
			self.httpstart = 1
			l := copy(b, self.data)
			return l, nil

		case 1: //http read
			self.httpstart = 2
			l := copy(b, self.data)
			self.data = self.data[l:]
			return l, nil

		case 2: //http read
			self.httpstart = 0
			return 0, io.EOF
		}
	}

	l := copy(b, self.data)
	self.data = self.data[l:]
	return l, nil
}

func printHeader(h http.Header) string {
	var logbuf string

	for k, v := range h {
		logbuf += fmt.Sprintf("%s :%s\n", k, v)
	}
	return logbuf
}

func (self *httpReader) runClient() {

	var p = make([]byte, 1900)

	for {

		self.httpstart = 0
		l, err := self.Read(p)
		if err == io.EOF {
			return
		}
		if l > 8 {
			isReq, firstLine := isRequest(p)
			if isReq { //start new request
				timeStamp := <-self.timestamp

				self.HandleRequest(timeStamp, firstLine)
			}
		}
	}
}

// response
func (self *httpReader) runServer() {

	var p = make([]byte, 1900)

	for {
		self.httpstart = 0
		l, err := self.Read(p)
		if err == io.EOF {
			return
		}
		if l > 8 {
			isResp, firstLine := isResponse(p)
			if isResp { //start new response
				timeStamp := <-self.timestamp

				self.HandleResponse(timeStamp, firstLine)

			}
		}
	}
}

func (h *httpReader) HandleRequest(timeStamp int64, firstline string) {

	b := bufio.NewReader(h)

	req, err := http.ReadRequest(b)
	//h.parent.UpdateReq(req, timeStamp, firstline)

	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return
	} else if err != nil {
		fmt.Println("HTTP-request", "HTTP Request error:", err)

	} else {

		req.ParseMultipartForm(defaultMaxMemory)

		r, ok := DecompressBody(req.Header, req.Body)
		if ok {
			defer r.Close()
		}
		contentType := req.Header.Get("Content-Type")
		logbuf := fmt.Sprintf("%v->%v:%v->%v\n", h.srcip, h.dstip, h.srcport, h.dstport)
		logbuf += printRequest(req)

		if strings.Contains(contentType, "application/json") {

			bodydata, err := ioutil.ReadAll(r)
			if err == nil {
				var jsonValue interface{}
				err = json.Unmarshal([]byte(bodydata), &jsonValue)
				if err == nil {
					logbuf += fmt.Sprintf("%#v\n", jsonValue)
				}
			}
		}

		fmt.Printf("%s", logbuf)
	}

}

func (self *httpReader) HandleResponse(timestamp int64, firstline string) {

	b := bufio.NewReader(self)

	resp, err := http.ReadResponse(b, nil)
	//self.parent.UpdateResp(resp, timestamp, firstline)

	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return
	} else if err != nil {
		fmt.Println("HTTP-reponse", "HTTP Response error: ", err)

	} else {

		r, ok := DecompressBody(resp.Header, resp.Body)
		if ok {
			defer r.Close()
		}
		contentType := resp.Header.Get("Content-Type")
		logbuf := fmt.Sprintf("%v->%v:%v->%v\n", self.srcip, self.dstip, self.srcport, self.dstport)
		logbuf += printResponse(resp)

		if strings.Contains(contentType, "application/json") {

			bodydata, err := ioutil.ReadAll(r)
			if err == nil {
				var jsonValue interface{}
				err = json.Unmarshal([]byte(bodydata), &jsonValue)
				if err == nil {
					logbuf += fmt.Sprintf("%#v\n", jsonValue)
				}
			}
		} else if strings.Contains(contentType, "application/javascript") {
			bodydata, err := ioutil.ReadAll(r)
			bodylen := len(bodydata)

			if err == nil {

				logbuf += fmt.Sprintf("%s\n", string(bodydata[:bodylen]))
			}
		} else if strings.Contains(contentType, "text/html") {
			bodydata, err := ioutil.ReadAll(r)
			bodylen := len(bodydata)

			if err == nil {
				logbuf += fmt.Sprintf("%s\n", string(bodydata[:bodylen]))
			}
		} else if strings.Contains(contentType, "text/plain") { //default text/plain
			bodydata, err := ioutil.ReadAll(r)
			bodylen := len(bodydata)

			if err == nil {
				logbuf += fmt.Sprintf("%s\n", string(bodydata[:bodylen]))
			}
		}

		fmt.Printf("%s", logbuf)
	}

}

func DecompressBody(header http.Header, reader io.ReadCloser) (io.ReadCloser, bool) {
	contentEncoding := header.Get("Content-Encoding")
	var nr io.ReadCloser
	var err error
	if contentEncoding == "" {
		// do nothing
		return reader, false
	} else if strings.Contains(contentEncoding, "gzip") {
		nr, err = gzip.NewReader(reader)
		if err != nil {
			return reader, false
		}
		return nr, true
	} else if strings.Contains(contentEncoding, "deflate") {
		nr, err = zlib.NewReader(reader)
		if err != nil {
			return reader, false
		}
		return nr, true
	} else {
		return reader, false
	}
}

func printRequest(req *http.Request) string {

	logbuf := fmt.Sprintf("\n")
	logbuf += fmt.Sprintf("%s\n", req.Host)
	logbuf += fmt.Sprintf("%s %s %s \n", req.Method, req.RequestURI, req.Proto)
	logbuf += printHeader(req.Header)
	logbuf += printForm(req.Form)
	logbuf += printForm(req.PostForm)
	if req.MultipartForm != nil {
		logbuf += printForm(url.Values(req.MultipartForm.Value))
	}
	logbuf += fmt.Sprintf("\n")
	return logbuf
}

func printForm(v url.Values) string {
	var logbuf string

	logbuf += fmt.Sprint("\n**************\n")
	for k, data := range v {
		logbuf += fmt.Sprint(k, ":")
		for _, d := range data {
			logbuf += fmt.Sprintf("%s", d)
		}
		logbuf += "\n"
	}
	logbuf += fmt.Sprint("**************\n")

	return logbuf
}

func printResponse(resp *http.Response) string {

	logbuf := fmt.Sprintf("\n")
	logbuf += fmt.Sprintf("%s %s\n", resp.Proto, resp.Status)
	logbuf += printHeader(resp.Header)
	logbuf += fmt.Sprintf("\n")
	return logbuf
}
