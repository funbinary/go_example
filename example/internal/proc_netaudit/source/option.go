package source

import "time"

type option func(*PcapSource)

// 设置超时
func WithTimeout(seconds int) option {
	return func(s *PcapSource) {
		s.timeout = time.Duration(seconds) * time.Second
	}
}

func WithBufferSize(size int) option {
	return func(s *PcapSource) {
		s.bufferSize = size
	}
}

func WithSnaplen(len int) option {
	return func(s *PcapSource) {
		s.snaplen = len
	}
}

func WithPromisc(promisc bool) option {
	return func(s *PcapSource) {
		s.promisc = promisc
	}
}
