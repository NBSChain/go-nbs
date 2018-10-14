package engine

import (
	"github.com/libp2p/go-libp2p-peer"
	"math"
	"sync"
)

type SwapLedger interface {

	Score() float64

	Threshold() float64
}


//TODO::change the mathematic function.
const LedgerThreshold = 0.1

type swapLedger struct {
	received	float64
	sent		float64
	score		float64
}

type ledgerEngine struct {
	sync.Mutex
	ledger	map[peer.ID]*swapLedger
}

func NewLedgerEngine()  *ledgerEngine{
	return &ledgerEngine{
		ledger:make(map[peer.ID]*swapLedger),
	}
}

/*****************************************************************
*
*		interface SwapLedger implements.
*
*****************************************************************/
func (l swapLedger) Score() float64{
	return l.score
}

func (l swapLedger) Threshold() float64{
	return LedgerThreshold
}

func (l swapLedger) send(dataLen float64){
	l.sent += dataLen
	l.calculateScore()
}

func (l swapLedger) gotData(dataLen float64){
	l.received += dataLen
	l.calculateScore()
}

func (l swapLedger) calculateScore(){
	r := l.sent / (l.received + 1)
	l.score = 1 - 1 / (1 - (math.Exp(6 -3*r)))
}

/*****************************************************************
*
*		interface LedgerEngine implements.
*
*****************************************************************/

func (engine *ledgerEngine) ReceiveData(fromNode peer.ID, data []byte) SwapLedger{

	engine.Lock()
	defer engine.Unlock()

	ledger := engine.getOrCreateLedger(fromNode)

	ledger.gotData(float64(len(data)))

	return ledger
}

func (engine *ledgerEngine) SupportData(toNode peer.ID, data []byte) (SwapLedger){

	engine.Lock()
	defer engine.Unlock()

	ledger := engine.getOrCreateLedger(toNode)

	ledger.send(float64(len(data)))

	return ledger
}

func (engine *ledgerEngine) GetLedger(nodeId peer.ID) SwapLedger{

	if l, ok := engine.ledger[nodeId]; ok{
		return l
	}

	return nil
}

func (engine *ledgerEngine) getOrCreateLedger(nodeId peer.ID) *swapLedger{

	ledger, ok := engine.ledger[nodeId]
	if !ok {
		ledger = &swapLedger{
			received: 0,
			score: 1.0,
			sent:0,
		}
		engine.ledger[nodeId] = ledger
	}

	return ledger
}