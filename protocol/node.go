package protocol

import (
	"ZETA_server/ZetaDB"
	"ZETA_server/network"
	"crypto/rand"
	"math/big"
	"sync"
	"time"
)

var nodeInstance *Node
var nodeOnce sync.Once

//call this function to get node
func GetNode() *Node {

	nodeOnce.Do(func() {
		nodeInstance = &Node{}
	})

	return nodeInstance
}

type State uint8

const (
	FOLLOWER  State = 1
	LEADER    State = 2
	CANDIDATE State = 3
	NEWCOMER  State = 4
)

type Node struct {
	//node state
	CurrentTermId     uint32
	CurrentSyncId     uint32 //records with syncId <= CurrentSyncId is acknowledged
	CurrentNodeSyncId uint32
	CurrentNodeNum    uint32 //equals to len(AddressLog)
	CurrentState      State  //LEADER/FOLLOWER/CANDIDATE/NEWCOMER

	//for LEADER to push syncId and nodeSyncId
	//if more than half of nodes acknowledge the syncId, it is acknowledged
	//only if all current nodes acknowledge the nodeSyncId, it is acknowledged
	AckedSyncIds           map[string]uint32 //map a node address to its latest acked syncId
	AckedNodeSyncIds       map[string]uint32 //map a node address to its latest acked nodeSyncId
	LatestSyncIdAckNum     uint32
	LatestNodeSyncIdAckNum uint32

	//record received vote number
	VoteNum uint32

	//node log
	SqlLog     []SqlRecord
	AddressLog []string //AddressLog does not contain SelfAddr

	//for LEADER, store addresses, dropAddresses and sqls in sended but unacked heartbeat
	//for FOLLOWER, store addresses, dropAddresses and sqls for next heartbeatResponse
	AddressInMessage     []string
	DropAddressInMessage []string
	SqlInMessage         []SqlRecord

	//for FOLLOWER, store addresses, dropAddresses and sqls from heartbeat with CurrentSyncId+1 and CurrentNodeSyncId+1
	//for LEADER, store addresses, dropAddresses and sqls in next heartbeat to send
	AddressBuffer     []string
	DropAddressBuffer []string
	SqlBuffer         []SqlRecord

	// ---------- below are some unchanged values when system runs ----------

	//main function argument
	ElectionLowBound  uint32 //unit in millisecond
	ElectionHighBound uint32 //unit in millisecond
	HeartbeatTimeout  uint32 //unit in millisecond
	SelfAddr          string
	Mode              string
	SimulatorAddr     string
	Prot              string

	//randomized when electionTimer is started
	//between ElectionLowBound and ElectionHighBound
	ElectionTimeout uint32 //unit in millisecond

	//timer
	ElectionTimer  time.Timer
	HeartbeatTimer time.Timer

	//system structure
	Db              ZetaDB.Database
	ResponseChannel chan network.Session
}

func (n *Node) StartElectionTimer() time.Timer {

	randomN, _ := rand.Int(rand.Reader, big.NewInt(100))
	randomF64 := float64(randomN.Int64()) / 100 * (float64(n.ElectionHighBound) - float64(n.ElectionLowBound))
	randomInt := int(randomF64) + int(n.ElectionLowBound)

	randomTime := time.Duration(randomInt)

	return *time.NewTimer(randomTime * time.Millisecond)
}

func (n *Node) StopElectionTimer() {}

func (n *Node) StartHeartbeatTimer() {}

func (n *Node) StopHeartbeatTimer() {}
