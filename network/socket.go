package network

import (
	"log"
	"net"
)

func Listen(selfAddress string, requestSessionChannel chan Session) {
	//listen from this tcp address
	tcp_addr, _ := net.ResolveTCPAddr("tcp4", selfAddress)

	listener, _ := net.ListenTCP("tcp", tcp_addr)

	for {
		log.Println("[server] listening", tcp_addr.String())

		// wait for client connection
		conn, err := listener.Accept()
		if err != nil {
			log.Println("[server] listening error", err)
			continue
		}

		//fetch request
		buffer := make([]byte, 256)
		conn.Read(buffer)

		//first byte is used to denote message type
		headByte := buffer[0]
		buffer = buffer[1:]

		var msgT MsgType
		switch headByte {
		case byte(0):
			msgT = GREETING
		case byte(1):
			msgT = HEARTBEAT
		case byte(2):
			msgT = HEARTBEATRESPONSE
		case byte(3):
			msgT = VOTEREQUEST
		case byte(4):
			msgT = VOTE
		case byte(5):
			msgT = CLIENTREQUEST
		case byte(6):
			msgT = CLIENTRESPONSE
		case byte(7):
			msgT = FAREWELL
		case byte(8):
			msgT = FAREWELLRESPONSE
		case byte(9):
			msgT = SIMULATORCONTROL
		}

		//generate a session
		newSession := Session{
			Connection: conn,
			Type:       msgT,
			Bytes:      buffer,
		}

		//push the request into channel
		requestSessionChannel <- newSession
	}
}

func Reply(responseSessionChannel chan Session, mode string, simulatorAddr string) {

	for {
		responseSession := <-responseSessionChannel

		msgT := responseSession.Type
		var headByte byte
		switch msgT {
		case GREETING:
			headByte = byte(0)
		case HEARTBEAT:
			headByte = byte(1)
		case HEARTBEATRESPONSE:
			headByte = byte(2)
		case VOTEREQUEST:
			headByte = byte(3)
		case VOTE:
			headByte = byte(4)
		case CLIENTREQUEST:
			headByte = byte(5)
		case CLIENTRESPONSE:
			headByte = byte(6)
		case FAREWELL:
			headByte = byte(7)
		case FAREWELLRESPONSE:
			headByte = byte(8)
		case SIMULATORCONTROL:
			headByte = byte(9)
		}

		var bytes []byte
		bytes = append(bytes, headByte)
		bytes = append(bytes, responseSession.Bytes...)

		if mode == "released" {
			responseSession.Connection.Write(bytes)
			log.Println("[server] response to:", responseSession.Connection.RemoteAddr().String())
		} else if mode == "test" {
			//assign simulator address
			tcp_addr, _ := net.ResolveTCPAddr("tcp4", simulatorAddr)
			// issue connection requirement
			conn, _ := net.DialTCP("tcp", nil, tcp_addr)

			conn.Write(bytes)
			log.Println("[server] response to:", conn.RemoteAddr().String())
		}

	}
}
