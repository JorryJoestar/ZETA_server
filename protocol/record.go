package protocol

type SqlRecord struct {
	SyncId        uint32
	SqlId         uint32
	Sql           string
	ClientAddress string
	UserId        int32
}
