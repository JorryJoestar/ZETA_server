package protocol

//event
func (n *Node) Receive_Greeting(g Greeting) {

	//if this node is a FOLLOWER, add new address into buffer, send to LEADER in next heartbeatResponse
	//if this node is a LEADER, add new address into buffer, send to FOLLOWERs in next heartbeat
	if n.CurrentState == FOLLOWER || n.CurrentState == LEADER {
		n.AddressBuffer = append(n.AddressBuffer, g.NewAddr)
	}

}

func (n *Node) Receive_HEARTBEAT(h Heartbeat) {

	//this node is a FOLLOWER and the heartbeat holds the correct termId
	if n.CurrentState == FOLLOWER && n.CurrentTermId >= h.TermId {
		n.CurrentTermId = h.TermId

		//reset election timer
		n.StartElectionTimer()

		//set sendTermId
		sendTermId := n.CurrentTermId

		//set AddressBuffer
		if h.NodeSyncId == n.CurrentNodeSyncId+1 && len(n.AddressBuffer) == 0 {
			n.AddressBuffer = append(n.AddressBuffer, h.Addresses...)
		}

		//set SqlBuffer
		if h.SyncId == n.CurrentSyncId+1 && len(n.SqlBuffer) == 0 {
			n.SqlBuffer = append(n.SqlBuffer, h.SqlRecords...)
		}

		//set sendAckSyncId
		var sendAckSyncId uint32
		if len(n.SqlBuffer) == 0 {
			sendAckSyncId = n.CurrentSyncId
		} else {
			sendAckSyncId = n.CurrentSyncId + 1
		}

		//set sendAckNodeSyncId
		var sendAckNodeSyncId uint32
		if len(n.AddressBuffer) == 0 {
			sendAckNodeSyncId = n.CurrentNodeSyncId
		} else {
			sendAckNodeSyncId = n.CurrentNodeSyncId + 1
		}

		//send back heartbeatResponse
		n.Send_HeartbeatResponse(h.FromAddr, sendTermId, sendAckSyncId, sendAckNodeSyncId, n.SqlInMessage, n.AddressInMessage, n.DropAddressInMessage)

		//clear SqlInMessage
		n.SqlInMessage = nil

		//clear AddressInMessage
		n.AddressInMessage = nil
	}

	//this node is a NEWCOMER, it would be woke by a heartbeat
	if n.CurrentState == NEWCOMER {
		n.CurrentState = FOLLOWER
		n.CurrentTermId = h.TermId

		//set election timer
		n.StartElectionTimer()

		//send back heartbeatResponse
		n.Send_HeartbeatResponse(h.FromAddr, n.CurrentTermId, 0, 0, nil, nil, nil)
	}

	//this node is a CANDIDATE or LEADER, and it has a lower termId compared with heartbeat
	if (n.CurrentState == CANDIDATE || n.CurrentState == LEADER) && n.CurrentTermId < h.TermId {
		n.CurrentState = FOLLOWER
		n.CurrentTermId = h.TermId

		//set election timer
		n.StartElectionTimer()

		//clear SqlInMessage
		n.SqlInMessage = nil

		//clear AddressInMessage
		n.AddressInMessage = nil

		//clear AddressBuffer
		n.AddressBuffer = nil

		//clear SqlBuffer
		n.SqlBuffer = nil

		//send back heartbeatResponse
		n.Send_HeartbeatResponse(h.FromAddr, n.CurrentTermId, n.CurrentSyncId, n.CurrentNodeSyncId, nil, nil, nil)
	}
}

func (n *Node) Receive_HEARTBEATRESPONSE(h HeartbeatResponse) {
	if n.CurrentState == LEADER && n.CurrentTermId == h.TermId {

		//add addresses
		n.AddressBuffer = append(n.AddressBuffer, h.Addresses...)

		//add sqls
		n.SqlBuffer = append(n.SqlBuffer, h.SqlRecords...)

		//update AckedNodeSyncIds
		n.AckedNodeSyncIds[h.FromAddr] = h.AckNodeSyncId

		//update AckedSyncIds
		n.AckedSyncIds[h.FromAddr] = h.AckSyncId

		//update LatestNodeSyncIdAckNum
		if h.AckNodeSyncId == n.CurrentNodeSyncId+1 {
			n.LatestNodeSyncIdAckNum++
		}

		//update LatestSyncIdAckNum
		if h.AckSyncId == n.CurrentSyncId+1 {
			n.LatestSyncIdAckNum++
		}

		//check LatestNodeSyncIdAckNum
		if n.LatestNodeSyncIdAckNum == n.CurrentNodeNum {

			//update CurrentNodeSyncId
			n.CurrentNodeSyncId++

			//update AddressLog
			n.AddressLog = append(n.AddressLog, n.AddressInMessage...)
			n.AddressInMessage = nil
		}

		//check LatestSyncIdAckNum
		if n.LatestSyncIdAckNum > n.CurrentNodeNum/2 {

			//update CurrentSyncId
			n.CurrentSyncId++

			//update SqlLog
			n.SqlLog = append(n.SqlLog, n.SqlInMessage...)
			n.SqlInMessage = nil
		}

	}
}

func (n *Node) Receive_VOTEREQUEST(v VoteRequest) {
	if n.CurrentState == FOLLOWER && n.CurrentTermId < v.TermId {
		n.Set_ElectionTimeout()
		n.CurrentTermId = v.TermId

		//send vote
		n.Send_Vote(v.FromAddr, n.CurrentTermId)

	}
}

func (n *Node) Receive_VOTE(v Vote) {
	if n.CurrentState == CANDIDATE && n.CurrentTermId == v.TermId {
		n.VoteNum++
		if n.VoteNum == n.CurrentNodeNum/2 {

			//state to LEADER
			n.CurrentState = LEADER

			//send heartbeat to every node
			for _, nodeAddr := range n.AddressLog {
				n.Send_Heartbeat(nodeAddr, n.CurrentTermId, n.CurrentSyncId, n.CurrentNodeSyncId, nil, nil, nil)
			}

			//start heartbeat timer
			n.StartHeartbeatTimer()
		}
	}
}

func (n *Node) Receive_CLIENTREQUEST(c ClientRequest) {

	//check if this request is read only
	if n.Check_SqlReadOnly(c.Sql) {
		stateCode, result := n.Db.ExecuteSql(c.UserId, c.Sql)

		//send back clientResponse
		n.Send_ClientResponse(c.FromAddr, stateCode, result)
	} else { // add it into next heartbeatResponse
		newSqlRecord := SqlRecord{
			Sql:           c.Sql,
			ClientAddress: c.FromAddr,
			UserId:        c.UserId,
		}
		n.SqlInMessage = append(n.SqlInMessage, newSqlRecord)
	}
}

func (n *Node) Receive_FAREWELL(f Farewell) {
	//if it is a LEADER, add farewell address into next heartbeat
	if n.CurrentState == LEADER {
		n.DropAddressBuffer = append(n.DropAddressBuffer, f.FromAddr)
	}

	//if it is a FOLLOWER, add farewell address into next heartbeatResponse
	if n.CurrentState == FOLLOWER {
		n.DropAddressInMessage = append(n.DropAddressInMessage, f.FromAddr)
	}
}

func (n *Node) Receive_FAREWELLRESPONSE(f FarewellResponse) {
	//TODO leave from group permitted
}

func (n *Node) Receive_SIMULATORCONTROL(s SimulatorControl) {
	//TODO receive some simulator signals
}
