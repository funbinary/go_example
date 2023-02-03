package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/funbinary/go_example/pkg/encoding/bjson"
	"net"
)

var quitSemaphore chan bool

type Version struct {
	Version string `json:"version,omitempty"`
}

type Command struct {
	Command   string     `json:"command,omitempty"`
	Arguments *Arguments `json:"arguments,omitempty"`
}

type Arguments struct {
	Iface string `json:"iface,omitempty"`
}

func main() {
	socketPath := "/kds/custom.socket"
	socket, err := NewSuricataSocket(socketPath)
	if err != nil || socket == nil {
		warn := &struct {
			Message string `json:"message"`
		}{
			Message: fmt.Sprintf("打开%s失败", socketPath),
		}
		fmt.Println(bjson.Marshal(warn))
		return
	}
	defer socket.Close()
	v := &Command{
		Command: "shutdown",
	}
	if err = socket.send(v); err != nil {
		warn := &struct {
			Message string `json:"message"`
		}{
			Message: fmt.Sprintf("send error: %v", err),
		}
		fmt.Println(bjson.Marshal(warn))
		return
	}
	s, err := socket.Read()
	if err != nil {
		warn := &struct {
			Message string `json:"message"`
		}{
			Message: fmt.Sprintf("read error: %v", err),
		}
		fmt.Println(bjson.Marshal(warn))
		return
	}
	fmt.Println(s)
}

type SuricataSocket struct {
	conn *net.UnixConn
	r    *bufio.Reader
}

func NewSuricataSocket(socketPath string) (*SuricataSocket, error) {
	unixAddr, err := net.ResolveUnixAddr("unix", socketPath)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUnix("unix", nil, unixAddr)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(conn)
	socket := &SuricataSocket{
		conn: conn,
		r:    r,
	}
	v := &Version{
		Version: "0.2",
	}
	if err := socket.send(v); err != nil {
		return nil, err
	}
	s, err := socket.Read()
	if err != nil {
		return nil, err
	}
	fmt.Println(s)
	return socket, nil
}

func (s *SuricataSocket) send(data interface{}) error {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Println("send:", string(jsonStr))
	_, err = s.conn.Write(append(jsonStr, '\n'))
	return err
}

func (s *SuricataSocket) Read() (string, error) {
	msg, _, err := s.r.ReadLine()
	return string(msg), err
}

func (s *SuricataSocket) Close() {
	s.conn.Close()
}
