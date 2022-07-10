package protocol

func (n *Node) Reach_ElectionTimeout() {
	//node ready to compete for being LEADER

	//transfer state to CANDIDATE
	n.CurrentState = CANDIDATE

	//vote to itself
	n.VoteNum = 1

	//update termId
	n.CurrentTermId++

	//set voteRequest to all other nodes
	for _, nodeAddr := range n.AddressLog {
		n.Send_VoteRequest(nodeAddr, n.CurrentTermId)
	}

	//set election timer
	n.StartElectionTimer()
}

func (n *Node) Reach_HeartbeatTimeout() {
	//when heartbeat timeout is reached, as a LEADER it will send heartbeat to all other nodes
	for _, nodeAddr := range n.AddressLog {

		//create sendSyncId, sendNodeSyncId, senSqlRecords, sendAddresses, sndDropAddresses
		var sendSyncId uint32
		var sendNodeSyncId uint32
		var senSqlRecords []SqlRecord
		var sendAddresses []string
		var sndDropAddresses []string

		//set sendSyncId and senSqlRecords according to AckedSyncIds[nodeAddr]
		if n.AckedSyncIds[nodeAddr] == n.CurrentSyncId { //this node is up to date on sqls

		} else {

		}

		//set sendNodeSyncId and sendAddresses according to AckedNodeSyncIds[nodeAddr]
		if n.AckedNodeSyncIds[nodeAddr] == n.CurrentNodeSyncId { //this node is up to date on addresses

		} else {

		}

		n.Send_Heartbeat(nodeAddr, n.CurrentTermId, sendSyncId, sendNodeSyncId, senSqlRecords, sendAddresses, sndDropAddresses)
	}
}

func (n *Node) Set_ElectionTimeout()   {}
func (n *Node) Set_HeartbeatTimeout()  {}
func (n *Node) Stop_ElectionTimeout()  {}
func (n *Node) Stop_HeartbeatTimeout() {}
