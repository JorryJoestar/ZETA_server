package node

type Node struct {
	CurrentTermId  uint32
	CurrentSyncId  uint32
	CurrentNodeNum int32

	AddressBuffer []AddressRecord
	AddressLog    []AddressRecord

	SqlBuffer []SqlRecord
	SqlLog    []SqlRecord

	//unit in millisecond
	ElectionLowBound  int32
	ElectionHighBound int32
	ElectionTimeout   int32
	HeartbeatTimeout  int32

	VoteNum uint32
}
