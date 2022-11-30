package sniffer

import (
	log "github.com/funbinary/go_example/pkg/blog"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"
	"io"
	"time"
)

type fileHandler struct {
	pcapHandle *pcap.Handle
	file       string

	loopCount, maxLoopCount int

	topSpeed bool
	lastTS   time.Time
}

func newFileHandler(file string, topSpeed bool, maxLoopCount int) (*fileHandler, error) {
	h := &fileHandler{
		file:         file,
		topSpeed:     topSpeed,
		maxLoopCount: maxLoopCount,
	}
	if err := h.open(); err != nil {
		return nil, err
	}

	return h, nil
}

func (h *fileHandler) open() (err error) {
	h.pcapHandle, err = pcap.OpenOffline(h.file)
	return
}

func (h *fileHandler) ReadPacketData() ([]byte, gopacket.CaptureInfo, error) {
	data, ci, err := h.pcapHandle.ReadPacketData()
	if err != nil {
		if err != io.EOF { //nolint:errorlint // io.EOF should never be wrapped.
			return data, ci, err
		}

		h.pcapHandle.Close()
		h.pcapHandle = nil

		h.loopCount++
		if h.loopCount >= h.maxLoopCount {
			return data, ci, err
		}

		log.Debugf("Reopening the file:%v", h.file)
		if err = h.open(); err != nil {
			return nil, ci, errors.Errorf("failed to reopen file: %w", err)
		}

		data, ci, err = h.pcapHandle.ReadPacketData()
		h.lastTS = ci.Timestamp
		return data, ci, err
	}

	if h.topSpeed {
		return data, ci, nil
	}

	if !h.lastTS.IsZero() {
		sleep := ci.Timestamp.Sub(h.lastTS)
		if sleep > 0 {
			time.Sleep(sleep)
		} else {
			log.Warnf("Time in pcap went backwards: %d", sleep)
		}
	}

	h.lastTS = ci.Timestamp
	ci.Timestamp = time.Now()
	return data, ci, nil
}

func (h *fileHandler) LinkType() layers.LinkType {
	return h.pcapHandle.LinkType()
}

func (h *fileHandler) Close() {
	if h.pcapHandle != nil {
		h.pcapHandle.Close()
		h.pcapHandle = nil
	}
}
