package cmd

import (
	"fmt"
	rtmp "github.com/zhangpeihao/gortmp"
)

type RtmpOutboundConnHandler struct {
	StreamChan    chan rtmp.OutboundStream
	VideoDataSize int64
	AudioDataSize int64

	Status uint
}

func NewRtmpOutboundConnHandler() *RtmpOutboundConnHandler {
	return &RtmpOutboundConnHandler{
		StreamChan: make(chan rtmp.OutboundStream),
	}
}
func (h *RtmpOutboundConnHandler) OnStatus(conn rtmp.OutboundConn) {
	var err error
	h.Status, err = conn.Status()
	fmt.Printf("@@@@@@@@@@@@@status: %d, err: %v\n", h.Status, err)
}

func (h *RtmpOutboundConnHandler) OnClosed(conn rtmp.Conn) {
	fmt.Printf("@@@@@@@@@@@@@Closed\n")
}

func (h *RtmpOutboundConnHandler) OnReceived(conn rtmp.Conn, message *rtmp.Message) {
	switch message.Type {
	case rtmp.VIDEO_TYPE:
		h.VideoDataSize += int64(message.Buf.Len())
	case rtmp.AUDIO_TYPE:
		h.VideoDataSize += int64(message.Buf.Len())
	}
}

func (h *RtmpOutboundConnHandler) OnReceivedRtmpCommand(conn rtmp.Conn, command *rtmp.Command) {
	fmt.Printf("ReceviedCommand: %+v\n", command)
}

func (h *RtmpOutboundConnHandler) OnStreamCreated(conn rtmp.OutboundConn, stream rtmp.OutboundStream) {
	fmt.Printf("Stream created: %d\n", stream.ID())
	h.StreamChan <- stream
}
