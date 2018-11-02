package idService

import (
	"github.com/mr-tron/base58/base58"
)

func IDB58Encode(id ID) string {
	return base58.Encode([]byte(id))
}
