package idPoolTools

import (
	"github.com/cilium/cilium/pkg/idpool"
)

type IdPool struct {
	idpool.IDPool
	members map[string]string
}

func IsValid(id int64) bool {
	return ((id > 0) && (id <= 9223372036854775807))
}

func (*IdPool) IsValid() {

}
