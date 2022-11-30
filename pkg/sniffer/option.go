package sniffer

import "time"

type Option func(sniffer *Sniffer)

// 设置超时
func WithTimeout(seconds int) Option {
	return func(s *Sniffer) {
		s.timeout = time.Duration(seconds) * time.Second
	}
}

func WithBufferSize(size int) Option {
	return func(s *Sniffer) {
		s.bufferSize = size
	}
}

func WithSnaplen(len int) Option {
	return func(s *Sniffer) {
		s.snaplen = len
	}
}

func WithPromisc(promisc bool) Option {
	return func(s *Sniffer) {
		s.promisc = promisc
	}
}

func WithFile(name string) Option {
	return func(s *Sniffer) {
		s.file = name
	}
}

func WithBpfFilter(filter string) Option {
	return func(s *Sniffer) {
		s.bpfFilter = filter
	}
}

func WithDumpFiile(path string) Option {
	return func(s *Sniffer) {
		s.dumpFile = path
	}
}

func WithDevice(name string) Option {
	return func(s *Sniffer) {
		s.device = name
	}
}
