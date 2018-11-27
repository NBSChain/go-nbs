package memership

import (
	"crypto/rand"
	"math/big"
)

func (node *MemManager) choseRandomInPartialView() *peerNodeItem {
	count := len(node.partialView)
	j := 0
	random, _ := rand.Int(rand.Reader, big.NewInt(int64(count)))

	for _, item := range node.partialView {
		if j == int(random.Int64()) {
			return item
		} else {
			j++
		}
	}
	return nil
}
