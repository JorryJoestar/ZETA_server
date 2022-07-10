package main

import (
	"ZETA_server/ZetaDB"
	"ZETA_server/network"
	"ZETA_server/protocol"
	"encoding/json"
	"os"
	"strconv"
)

/*
	main function argument
		0: electionLowBound (ms)
		1: electionHighBound (ms)
		2: heartbeatTimeout (ms)
		3: self address
		4: fileLocation
		5: mode (released/test)
		6: simulator address
		7: prot (raft/vr)
		8- : knownAddresses
*/

func main() {
	// request and response session channel
	requestSessionChannel := make(chan network.Session, 16384)
	responseSessionChannel := make(chan network.Session, 16384)

	//get main args and use them to get node
	electionLowBound, _ := strconv.ParseUint(os.Args[0], 2, 32)
	electionHighBound, _ := strconv.ParseUint(os.Args[1], 2, 32)
	heartbeatTimeout, _ := strconv.ParseUint(os.Args[2], 2, 32)
	selfAddress := os.Args[3]
	fileLocation := os.Args[4]
	mode := os.Args[5]
	simulatorAddress := os.Args[6]
	prot := os.Args[7]
	knownAddresses := os.Args[8:]

	node := protocol.GetNode()

	//set node const value
	node.ElectionLowBound = uint32(electionLowBound)
	node.ElectionHighBound = uint32(electionHighBound)
	node.HeartbeatTimeout = uint32(heartbeatTimeout)
	node.SelfAddr = selfAddress
	node.Mode = mode
	node.SimulatorAddr = simulatorAddress
	node.Prot = prot
	node.Db = *ZetaDB.GetDatabase(fileLocation)
	node.ResponseChannel = responseSessionChannel

	//create two syncIds map
	node.AckedSyncIds = make(map[string]uint32)
	node.AckedNodeSyncIds = make(map[string]uint32)

	//initialize node
	node.CurrentTermId = 0
	node.CurrentSyncId = 0
	node.CurrentNodeSyncId = 0
	node.CurrentNodeNum = uint32(len(node.AddressLog)) + 1
	if node.CurrentNodeNum == 1 { //if AddressLog is empty, set it to LEADER state, else set to NEWCOMER state
		node.CurrentState = protocol.LEADER
	} else {
		node.CurrentState = protocol.NEWCOMER
	}
	node.AddressLog = knownAddresses

	//create socket receive and send thread
	go network.Listen(selfAddress, requestSessionChannel)
	go network.Reply(responseSessionChannel, mode, simulatorAddress)

	//if this node is a NEWCOMER
	if node.CurrentState == protocol.NEWCOMER {
		//send a GREETING to any one of known addresses
		node.Send_Greeting(knownAddresses[0])
	}

	//main loop
	for {
		select {
		case requestSession := <-requestSessionChannel:
			switch requestSession.Type {
			case network.GREETING: //GREETING
				var greeting protocol.Greeting
				json.Unmarshal(requestSession.Bytes, greeting)
				node.Receive_Greeting(greeting)
			case network.HEARTBEAT: //HEARTBEAT
				var heartbeat protocol.Heartbeat
				json.Unmarshal(requestSession.Bytes, heartbeat)
				node.Receive_HEARTBEAT(heartbeat)
			case network.HEARTBEATRESPONSE: //HEARTBEATRESPONSE
				var heartbeatResponse protocol.HeartbeatResponse
				json.Unmarshal(requestSession.Bytes, heartbeatResponse)
				node.Receive_HEARTBEATRESPONSE(heartbeatResponse)
			case network.VOTEREQUEST: // VOTEREQUEST
				var voteRequest protocol.VoteRequest
				json.Unmarshal(requestSession.Bytes, voteRequest)
				node.Receive_VOTEREQUEST(voteRequest)
			case network.VOTE: //VOTE
				var vote protocol.Vote
				json.Unmarshal(requestSession.Bytes, vote)
				node.Receive_VOTE(vote)
			case network.CLIENTREQUEST: //CLIENTREQUEST
				var clientRequest protocol.ClientRequest
				json.Unmarshal(requestSession.Bytes, clientRequest)
				node.Receive_CLIENTREQUEST(clientRequest)
			case network.CLIENTRESPONSE: //CLIENTRESPONSE
				//just ignore, impossible for a server node to receive client response
			case network.FAREWELL:
				var farewell protocol.Farewell
				json.Unmarshal(requestSession.Bytes, farewell)
				node.Receive_FAREWELL(farewell)
			case network.FAREWELLRESPONSE: //FAREWELLRESPONSE
				var farewellResponse protocol.FarewellResponse
				json.Unmarshal(requestSession.Bytes, farewellResponse)
				node.Receive_FAREWELLRESPONSE(farewellResponse)
			case network.SIMULATORCONTROL: //SIMULATORCONTROL
				var simulatorControl protocol.SimulatorControl
				json.Unmarshal(requestSession.Bytes, simulatorControl)
				node.Receive_SIMULATORCONTROL(simulatorControl)
			}
		case <-node.ElectionTimer.C: //election timeout
			node.Reach_ElectionTimeout()
		case <-node.HeartbeatTimer.C: //heartbeat timeout
			node.Reach_HeartbeatTimeout()
		}
	}
}
