package memership

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

//TODO:: make sure this random is ok
func (node *MemManager) choseRandomInPartialView(nodeId string) *viewNode {
	count := len(node.partialView)
	j := 0
	random, _ := rand.Int(rand.Reader, big.NewInt(int64(count)))

	for _, item := range node.partialView {
		if j == int(random.Int64()) && nodeId != item.nodeId {
			return item
		} else {
			j++
		}
	}
	return nil
}

func (node *MemManager) removeFromView(item *viewNode, views map[string]*viewNode) {

	delete(views, item.nodeId)

	if err := item.outConn.Close(); err != nil {
		logger.Warning(err)
	}

	logger.Warning("remove node from partial view:->", item.nodeId)

	updateProbability(views)
}

func (node *MemManager) unsubItem(item *viewNode) {

}

func (node *MemManager) getHeartBeat(task *msgTask) error {
	beat := task.msg.HeartBeat
	item, ok := node.inputView[beat.FromID]
	if !ok {
		err := fmt.Errorf("no such input view item:->%s", beat.FromID)
		logger.Warning(err)
		return err
	}
	item.updateTime = time.Now()
	node.updateTime = time.Now()
	return nil
}
