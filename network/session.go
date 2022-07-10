package network

import "net"

type MsgType uint16

const (
	GREETING          MsgType = 0
	HEARTBEAT         MsgType = 1
	HEARTBEATRESPONSE MsgType = 2
	VOTEREQUEST       MsgType = 3
	VOTE              MsgType = 4
	CLIENTREQUEST     MsgType = 5
	CLIENTRESPONSE    MsgType = 6
	FAREWELL          MsgType = 7
	FAREWELLRESPONSE  MsgType = 8
	SIMULATORCONTROL  MsgType = 9
)

type Session struct {
	Connection net.Conn
	Type       MsgType
	Bytes      []byte
}

func NewSession(destAddr string, Type MsgType, Bytes []byte) Session {
	s := Session{}

	//assign server address
	tcp_addr, _ := net.ResolveTCPAddr("tcp4", destAddr)

	// issue connection requirement
	conn, _ := net.DialTCP("tcp", nil, tcp_addr)

	s.Connection = conn
	s.Type = Type
	s.Bytes = Bytes

	return s
}
