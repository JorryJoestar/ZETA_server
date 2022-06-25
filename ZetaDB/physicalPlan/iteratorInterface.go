package execution

import (
	"ZETA_server/ZetaDB/container"
)

type Iterator interface {
	Open(iterator1 Iterator, iterator2 Iterator) error
	GetNext() (*container.Tuple, error)
	HasNext() bool
	Close()
	GetSchema() *container.Schema
}
