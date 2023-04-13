package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	rtmp "github.com/zhangpeihao/gortmp"
	"github.com/zhangpeihao/log"
)

const (
	programName = "RtmpPlayer"
	version     = "0.0.1"
)

var (
	url        *string = flag.String("URL", "rtmp://106.225.229.10/643606265f79c3425d097473", "The rtmp url to connect.")
	streamName *string = flag.String("Stream", "98b43e937320", "Stream name to play.")
)

type TestOutboundConnHandler struct {
}

var obConn rtmp.OutboundConn
var createStreamChan chan rtmp.OutboundStream
var videoDataSize int64
var audioDataSize int64
var status uint

func (handler *TestOutboundConnHandler) OnStatus(conn rtmp.OutboundConn) {
	var err error
	status, err = conn.Status()
	fmt.Printf("@@@@@@@@@@@@@status: %d, err: %v\n", status, err)
}

func (handler *TestOutboundConnHandler) OnClosed(conn rtmp.Conn) {
	fmt.Printf("@@@@@@@@@@@@@Closed\n")
}

func (handler *TestOutboundConnHandler) OnReceived(conn rtmp.Conn, message *rtmp.Message) {
	switch message.Type {
	case rtmp.VIDEO_TYPE:

		videoDataSize += int64(message.Buf.Len())
	case rtmp.AUDIO_TYPE:
		audioDataSize += int64(message.Buf.Len())
	}
}

func (handler *TestOutboundConnHandler) OnReceivedRtmpCommand(conn rtmp.Conn, command *rtmp.Command) {
	fmt.Printf("ReceviedCommand: %+v\n", command)
}

func (handler *TestOutboundConnHandler) OnStreamCreated(conn rtmp.OutboundConn, stream rtmp.OutboundStream) {
	fmt.Printf("Stream created: %d\n", stream.ID())
	createStreamChan <- stream
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s version[%s]\r\nUsage: %s [OPTIONS]\r\n", programName, version, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	fmt.Printf("rtmp:%s stream:%s \r\n", *url, *streamName)
	l := log.NewLogger(".", "player", nil, 60, 3600*24, true)
	rtmp.InitLogger(l)
	defer l.Close()

	createStreamChan = make(chan rtmp.OutboundStream)
	testHandler := &TestOutboundConnHandler{}
	fmt.Println("to dial")

	var err error

	obConn, err = rtmp.Dial(*url, testHandler, 100)
	/*
		conn := TryHandshakeByVLC()
		obConn, err = rtmp.NewOutbounConn(conn, *url, testHandler, 100)
	*/
	if err != nil {
		fmt.Println("Dial error", err)
		os.Exit(-1)
	}

	defer obConn.Close()
	fmt.Printf("obConn: %+v\n", obConn)
	fmt.Printf("obConn.URL(): %s\n", obConn.URL())
	fmt.Println("to connect")
	//	err = obConn.Connect("33abf6e996f80e888b33ef0ea3a32bfd", "131228035", "161114738", "play", "", "", "1368083579")
	err = obConn.Connect()
	if err != nil {
		fmt.Printf("Connect error: %s", err.Error())
		os.Exit(-1)
	}
	for {
		select {
		case stream := <-createStreamChan:
			// Play
			fmt.Printf("===================%s", *streamName)
			err = stream.Play(*streamName, nil, nil, nil)
			if err != nil {
				fmt.Printf("Play error: %s", err.Error())
				os.Exit(-1)
			}

		case <-time.After(1 * time.Second):
			fmt.Printf("Audio size: %d bytes; Vedio size: %d bytes\n", audioDataSize, videoDataSize)
		}
	}
}
