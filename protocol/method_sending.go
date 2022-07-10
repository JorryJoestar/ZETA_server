package protocol

import (
	"ZETA_server/network"
	"encoding/json"
)

//behavior
func (n *Node) Send_Greeting(toAddr string) {
	newGreeting := Greeting{
		FromAddr: n.SelfAddr,
		ToAddr:   toAddr,
		NewAddr:  n.SelfAddr,
	}

	bytes, _ := json.Marshal(newGreeting)
	if n.Mode == "test" {
		session := network.NewSession(n.SimulatorAddr, network.GREETING, bytes)
		n.ResponseChannel <- session
	} else if n.Mode == "released" {
		session := network.NewSession(toAddr, network.GREETING, bytes)
		n.ResponseChannel <- session
	}

}
func (n *Node) Send_Heartbeat(toAddr string, termId uint32, syncId uint32, fromSyncId uint32, nodeSyncId uint32, fromNodeSyncId uint32, sqlRecords []SqlRecord, addresses []string, dropAddresses []string) {
	newHeartbeat := Heartbeat{
		FromAddr:       n.SelfAddr,
		ToAddr:         toAddr,
		TermId:         termId,
		SyncId:         syncId,
		FromSyncId:     fromSyncId,
		NodeSyncId:     nodeSyncId,
		FromNodeSyncId: fromNodeSyncId,
		SqlRecords:     sqlRecords,
		Addresses:      addresses,
		DropAddresses:  dropAddresses,
	}

	bytes, _ := json.Marshal(newHeartbeat)
	if n.Mode == "test" {
		session := network.NewSession(n.SimulatorAddr, network.HEARTBEAT, bytes)
		n.ResponseChannel <- session
	} else if n.Mode == "released" {
		session := network.NewSession(toAddr, network.HEARTBEAT, bytes)
		n.ResponseChannel <- session
	}

}
func (n *Node) Send_HeartbeatResponse(toAddr string, termId uint32, ackSyncId uint32, ackNodeSyncId uint32, sqlRecords []SqlRecord, addresses []string, dropAddresses []string) {
	newHeartbeatResponse := HeartbeatResponse{
		FromAddr:      n.SelfAddr,
		ToAddr:        toAddr,
		TermId:        termId,
		AckSyncId:     ackSyncId,
		AckNodeSyncId: ackNodeSyncId,
		SqlRecords:    sqlRecords,
		Addresses:     addresses,
		DropAddresses: dropAddresses,
	}

	bytes, _ := json.Marshal(newHeartbeatResponse)
	if n.Mode == "test" {
		session := network.NewSession(n.SimulatorAddr, network.HEARTBEATRESPONSE, bytes)
		n.ResponseChannel <- session
	} else if n.Mode == "released" {
		session := network.NewSession(toAddr, network.HEARTBEATRESPONSE, bytes)
		n.ResponseChannel <- session
	}

}
func (n *Node) Send_VoteRequest(toAddr string, termId uint32) {
	newVoteRequest := VoteRequest{
		FromAddr: n.SelfAddr,
		ToAddr:   toAddr,
		TermId:   termId,
	}

	bytes, _ := json.Marshal(newVoteRequest)
	if n.Mode == "test" {
		session := network.NewSession(n.SimulatorAddr, network.VOTEREQUEST, bytes)
		n.ResponseChannel <- session
	} else if n.Mode == "released" {
		session := network.NewSession(toAddr, network.VOTEREQUEST, bytes)
		n.ResponseChannel <- session
	}

}
func (n *Node) Send_Vote(toAddr string, termId uint32) {
	newVote := Vote{
		FromAddr: n.SelfAddr,
		ToAddr:   toAddr,
		TermId:   termId,
	}

	bytes, _ := json.Marshal(newVote)
	if n.Mode == "test" {
		session := network.NewSession(n.SimulatorAddr, network.VOTE, bytes)
		n.ResponseChannel <- session
	} else if n.Mode == "released" {
		session := network.NewSession(toAddr, network.VOTE, bytes)
		n.ResponseChannel <- session
	}

}
func (n *Node) Send_ClientResponse(toAddr string, stateCode int32, result string) {
	newClientResponse := ClientResponse{
		FromAddr:  n.SelfAddr,
		ToAddr:    toAddr,
		StateCode: stateCode,
		Result:    result,
	}

	bytes, _ := json.Marshal(newClientResponse)
	if n.Mode == "test" {
		session := network.NewSession(n.SimulatorAddr, network.CLIENTRESPONSE, bytes)
		n.ResponseChannel <- session
	} else if n.Mode == "released" {
		session := network.NewSession(toAddr, network.CLIENTRESPONSE, bytes)
		n.ResponseChannel <- session
	}

}
func (n *Node) Send_Farewell(toAddr string) {
	newFarewell := Farewell{
		FromAddr: n.SelfAddr,
		ToAddr:   toAddr,
	}

	bytes, _ := json.Marshal(newFarewell)
	if n.Mode == "test" {
		session := network.NewSession(n.SimulatorAddr, network.FAREWELL, bytes)
		n.ResponseChannel <- session
	} else if n.Mode == "released" {
		session := network.NewSession(toAddr, network.FAREWELL, bytes)
		n.ResponseChannel <- session
	}

}
func (n *Node) Send_FarewellResponse(toAddr string) {
	newFarewellResponse := FarewellResponse{
		FromAddr: n.SelfAddr,
		ToAddr:   toAddr,
	}

	bytes, _ := json.Marshal(newFarewellResponse)
	if n.Mode == "test" {
		session := network.NewSession(n.SimulatorAddr, network.FAREWELLRESPONSE, bytes)
		n.ResponseChannel <- session
	} else if n.Mode == "released" {
		session := network.NewSession(toAddr, network.FAREWELLRESPONSE, bytes)
		n.ResponseChannel <- session
	}

}
