package protocol

import (
	"ZETA_server/node"
)

type MessageType uint16

const (
	GREETING           MessageType = 1
	HEARTBEAT          MessageType = 2
	HEARTBEAT_RESPONSE MessageType = 3
	VOTE_REQUEST       MessageType = 4
	VOTE               MessageType = 5
	CLIENT_REQUEST     MessageType = 6
	CLIENT_RESPONSE    MessageType = 7
	FAREWELL           MessageType = 8
	FAREWELLRESPONSE   MessageType = 9
)

type Message struct {
	Type       MessageType
	TermId     uint32
	SyncId     uint32
	FromSyncId uint32
	Addresses  []node.AddressRecord
	Sqls       []node.SqlRecord
	Response   string
}

func BytesToMessage(bytes []byte) Message {
	return Message{}
}

func (message Message) ToBytes() []byte {
	return nil
}
