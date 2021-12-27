package main

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

var Data = []byte("I'm alwaysbeta!")

type ping struct {
	Addr string
	Conn net.Conn
	Data []byte
}

func main() {
	ping, err := Run("10.249.42.227", Data)

	if err != nil {
		fmt.Println(err)
	}

	ping.Ping()
}

func MarshalMsg(req int, data []byte) ([]byte, error) {
	xid, xseq := os.Getpid()&0xffff, req
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: xid, Seq: xseq,
			Data: data,
		},
	}
	return wm.Marshal(nil)
}

func Run(addr string, data []byte) (*ping, error) {
	wb, err := MarshalMsg(1, data)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return &ping{Data: wb, Addr: addr}, nil
}

func (self *ping) Dail() (err error) {
	self.Conn, err = net.Dial("ip4:icmp", self.Addr)
	if err != nil {
		return err
	}
	return nil
}

func (self *ping) Ping() {
	if err := self.Dail(); err != nil {
		fmt.Println("ICMP error:", err)
		return
	}
	fmt.Println("Start ping from ", self.Conn.LocalAddr())
	sendPingMsg(self.Conn, self.Data)
}

func sendPingMsg(c net.Conn, wb []byte) {
	if _, err := c.Write(wb); err != nil {
		print(err)
	}
}
