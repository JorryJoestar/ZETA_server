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
	//create sendSyncId, sendNodeSyncId, senSqlRecords, sendAddresses, sendDropAddresses
	var sendSyncId uint32
	var sendNodeSyncId uint32
	var sendFromSyncId uint32
	var sendFromNodeSyncId uint32
	var senSqlRecords []SqlRecord
	var sendAddresses []string
	var sendDropAddresses []string

	//update AddressInMessage and DropAddressInMessage
	//if AddressInMessage is empty and AddressBuffer is not empty, push addresses in AddressBuffer into AddressInMessage
	if len(n.AddressInMessage) == 0 && len(n.AddressBuffer) != 0 {
		n.AddressInMessage = n.AddressBuffer
		n.AddressBuffer = nil
	}
	//if DropAddressInMessage is empty and DropAddressBuffer is not empty, push addresses in DropAddressBuffer into DropAddressInMessage
	if len(n.DropAddressInMessage) == 0 && len(n.DropAddressBuffer) != 0 {
		n.DropAddressInMessage = n.DropAddressBuffer
		n.DropAddressBuffer = nil
	}

	//when heartbeat timeout is reached, as a LEADER it will send heartbeat to all other nodes
	for _, nodeAddr := range n.AddressLog {

		//set sendSyncId, sendFromSyncId and senSqlRecords according to AckedSyncIds[nodeAddr] and len of
		if n.AckedSyncIds[nodeAddr] == n.CurrentSyncId { //this node is up to date on sqls
			//if SqlInMessage not empty, resend these sqls
			if len(n.SqlInMessage) != 0 {
				sendFromSyncId = n.CurrentSyncId + 1
				sendSyncId = n.CurrentSyncId + 1
				senSqlRecords = n.SqlInMessage
			} else { //if SqlInMessage is empty
				//check if SqlBuffer is empty
				if len(n.SqlBuffer) == 0 { //no new sqls to sync
					sendFromSyncId = n.CurrentSyncId
					sendSyncId = n.CurrentSyncId
					senSqlRecords = nil
				} else { //load sqls from buffer to SqlInMessage, plus 1 to sendSyncId
					sendFromSyncId = n.CurrentSyncId + 1
					sendSyncId = n.CurrentSyncId + 1

					//update every records
					beginIndex := uint32(len(n.SqlLog) + 1)
					loopIndex := beginIndex
					for _, record := range n.SqlBuffer {
						record.SqlId = loopIndex
						loopIndex++
						record.SyncId = sendSyncId
						n.SqlInMessage = append(n.SqlInMessage, record)
					}

					n.SqlBuffer = nil
					senSqlRecords = n.SqlInMessage

				}
			}
		} else { //this node is not up to date on sqls
			//fetch all sqlRecords that this node misses into heartbeat
			sendFromSyncId = n.AckedSyncIds[nodeAddr] + 1

			//find out all missing records from log
			var unSyncSqlRecords []SqlRecord
			for _, record := range n.SqlLog {
				if record.SyncId >= sendFromSyncId {
					unSyncSqlRecords = append(unSyncSqlRecords, record)
				}
			}
			senSqlRecords = unSyncSqlRecords

			//if SqlInMessage not empty
			if len(n.SqlInMessage) != 0 {
				sendSyncId = n.CurrentSyncId + 1

				//append records in SqlInMessage
				senSqlRecords = append(senSqlRecords, n.SqlInMessage...)
			} else { //if SqlInMessage is empty
				//check if SqlBuffer is empty
				if len(n.SqlBuffer) == 0 { //no new sqls to sync
					sendSyncId = n.CurrentSyncId
				} else { //load sqls from buffer to SqlInMessage, plus 1 to sendSyncId
					sendSyncId = n.CurrentSyncId + 1

					//update every records
					beginIndex := uint32(len(n.SqlLog) + 1)
					loopIndex := beginIndex
					for _, record := range n.SqlBuffer {
						record.SqlId = loopIndex
						loopIndex++
						record.SyncId = sendSyncId
						n.SqlInMessage = append(n.SqlInMessage, record)
					}

					n.SqlBuffer = nil
					senSqlRecords = n.SqlInMessage

				}
			}
		}

		//set sendNodeSyncId and sendAddresses according to AckedNodeSyncIds[nodeAddr]
		if n.AckedNodeSyncIds[nodeAddr] == n.CurrentNodeSyncId { //this node is up to date on addresses

			//check if AddressInMessage and DropAddressInMessage is empty
			if len(n.AddressInMessage) == 0 && len(n.DropAddressInMessage) == 0 {
				sendFromNodeSyncId = n.CurrentNodeSyncId
				sendNodeSyncId = n.CurrentNodeSyncId
				sendAddresses = nil
				sendDropAddresses = nil
			} else { //some new addresses alter should be pushed into next heartbeat
				sendFromNodeSyncId = n.CurrentNodeSyncId + 1
				sendNodeSyncId = n.CurrentNodeSyncId + 1
				sendAddresses = n.AddressInMessage
				sendDropAddresses = n.DropAddressInMessage
			}

		} else { //this node is not up to date on addresses
			//set sendFromNodeSyncId
			sendFromNodeSyncId = n.AckedNodeSyncIds[nodeAddr] + 1

			//set sendNodeSyncId
			if len(n.AddressInMessage) == 0 && len(n.DropAddressInMessage) == 0 {
				sendNodeSyncId = n.CurrentNodeSyncId
			} else {
				sendNodeSyncId = n.CurrentNodeSyncId + 1
			}

			//set sendAddresses and sendDropAddresses
			sendAddresses = n.AddressLog
			sendAddresses = append(sendAddresses, n.AddressInMessage...)
			sendDropAddresses = n.DropAddressInMessage
		}

		n.Send_Heartbeat(nodeAddr, n.CurrentTermId, sendSyncId, sendFromSyncId, sendNodeSyncId, sendFromNodeSyncId, senSqlRecords, sendAddresses, sendDropAddresses)
	}
}

func (n *Node) Set_ElectionTimeout()   {}
func (n *Node) Set_HeartbeatTimeout()  {}
func (n *Node) Stop_ElectionTimeout()  {}
func (n *Node) Stop_HeartbeatTimeout() {}
