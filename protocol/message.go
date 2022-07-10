package protocol

// Greeting 00
type Greeting struct {
	FromAddr string
	ToAddr   string

	NewAddr string
}

//Heartbeat 01
type Heartbeat struct {
	FromAddr string
	ToAddr   string

	TermId     uint32
	SyncId     uint32
	NodeSyncId uint32

	SqlRecords    []SqlRecord
	Addresses     []string
	DropAddresses []string
}

//HeartbeatResponse 02
type HeartbeatResponse struct {
	FromAddr string
	ToAddr   string

	TermId        uint32
	AckSyncId     uint32
	AckNodeSyncId uint32

	SqlRecords    []SqlRecord
	Addresses     []string
	DropAddresses []string
}

//VoteRequest 03
type VoteRequest struct {
	FromAddr string
	ToAddr   string

	TermId uint32
}

//Vote 04
type Vote struct {
	FromAddr string
	ToAddr   string

	TermId uint32
}

//ClientRequest 05
type ClientRequest struct {
	FromAddr string
	ToAddr   string

	UserId int32
	Sql    string
}

//ClientResponse 06
type ClientResponse struct {
	FromAddr string
	ToAddr   string

	StateCode int32
	Result    string
}

//Farewell 07
type Farewell struct {
	FromAddr string
	ToAddr   string
}

//FarewellResponse 08
type FarewellResponse struct {
	FromAddr string
	ToAddr   string
}

//SimulatorControl 09
type SimulatorControl struct {
	FromAddr string
	ToAddr   string

	Command string
}
