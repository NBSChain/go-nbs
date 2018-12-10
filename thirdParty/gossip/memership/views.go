package memership

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

//TODO:: make sure this random is ok
func (node *MemManager) choseRandomInPartialView() *viewNode {
	count := len(node.partialView)
	j := 0
	random, _ := rand.Int(rand.Reader, big.NewInt(int64(count)))
	logger.Debug("chose random in partialView :->", random)
	for _, item := range node.partialView {
		if j == int(random.Int64()) {
			return item
		} else {
			j++
		}
	}
	return nil
}

func (node *MemManager) removeFromView(item *viewNode, views map[string]*viewNode) {

	delete(views, item.nodeId)

	if item.outConn != nil {
		if err := item.outConn.Close(); err != nil {
			logger.Warning(err)
		}
	}

	logger.Warning("remove node from view:->", item.nodeId)

	updateProbability(views)
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
	return nil
}
