package engine

import (
	"github.com/NBSChain/go-nbs/storage/application/dataStore"
	"github.com/NBSChain/go-nbs/storage/bitswap/pb"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-peer"
	"math"
	"sync"
	"time"
)

type SwapLedger interface {

	Score() float64

	Threshold() float64
}

//TODO::change the mathematic function.

const LedgerDataKeyPrefix		= "local_ledger_data_key"
const LedgerThreshold 			= 0.1
const MaxCounterToLose  		= 1 << 10
const MaxTimerIntervalToSyncLedger 	= 30
var logger 				= utils.GetLogInstance()
var timeFormat				= utils.GetConfig().SysTimeFormat

type swapLedger struct {
	received	float64
	sent		float64
	score		float64
	updateTime	time.Time
}

type ledgerEngine struct {
	sync.Mutex
	counter		int
	workSignal	chan struct{}
	ledger      	map[peer.ID]*swapLedger
	ledgerStore 	dataStore.DataStore
}

func NewLedgerEngine()  *ledgerEngine{

	storeService  := dataStore.GetServiceDispatcher().ServiceByType(dataStore.ServiceTypeLocalParam)

	engine := &ledgerEngine{
		ledgerStore: 	storeService,
		workSignal:	make(chan struct{}),
		ledger:     	make(map[peer.ID]*swapLedger),
	}

	engine.loadSavedLedger()

	go engine.runLoop()

	return engine
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
	l.updateTime = time.Now()
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

	engine.metrics()

	return ledger
}

func (engine *ledgerEngine) SupportData(toNode peer.ID, data []byte) (SwapLedger){

	engine.Lock()
	defer engine.Unlock()

	ledger := engine.getOrCreateLedger(toNode)

	ledger.send(float64(len(data)))

	engine.metrics()

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
			received: 	0,
			score: 		1.0,
			sent:		0,
			updateTime:	time.Now(),
		}
		engine.ledger[nodeId] = ledger
	}

	return ledger
}
func (engine *ledgerEngine) metrics()  {

	if engine.counter += 1; engine.counter >= MaxCounterToLose{
		engine.counter = 0
		engine.workSignal <- struct{}{}
	}
}
func (engine *ledgerEngine) loadSavedLedger(){
	ledgerSet := &bitswap_pb.LedgerSet{
	}

	ledgerData, err := engine.ledgerStore.Get(LedgerDataKeyPrefix)
	if err != nil{
		logger.Warning("failed to Unmarshal saved ledger data.", err)
		return
	}

	if err := proto.Unmarshal(ledgerData, ledgerSet); err != nil{
		logger.Warning("failed to Unmarshal saved ledger data.", err)
		return
	}

	engine.Lock()
	defer engine.Unlock()

	for _, l := range ledgerSet.Ledgers{

		ledger := &swapLedger{
			received: 	l.Received,
			score: 		l.Score,
			sent: 		l.Sent,
		}

		ledger.updateTime, _ = time.Parse(timeFormat, l.UpdateTime)

		engine.ledger[peer.ID(l.Id)] = ledger
	}
}

//TODO:: if the ledger grow too big, we need control this and remove the oldest ones.
func (engine *ledgerEngine) syncLedgerToSave(){

	engine.Lock()

	ledgerSet := &bitswap_pb.LedgerSet{
		Ledgers:make([]*bitswap_pb.Ledger, len(engine.ledger)),
	}

	for key, value := range engine.ledger{

		l := &bitswap_pb.Ledger{
			Id:		string(key),
			Received:	value.received,
			Score:		value.score,
			Sent:		value.sent,
			UpdateTime:	value.updateTime.Format(timeFormat),
		}

		ledgerSet.Ledgers = append(ledgerSet.Ledgers, l)
	}

	engine.Unlock()

	ledgerSetBytes, err := proto.Marshal(ledgerSet)
	if err != nil{
		logger.Warning("failed to marshal ledger data.", err)
		return
	}

	if err := engine.ledgerStore.Put(LedgerDataKeyPrefix, ledgerSetBytes); err != nil{
		logger.Warning("failed to save ledger data to local data store.", err)
	}
}

func (engine *ledgerEngine) runLoop(){

	select {
	case <-engine.workSignal:
		go engine.syncLedgerToSave()
	case <-time.After(time.Minute * MaxTimerIntervalToSyncLedger):
		go engine.syncLedgerToSave()
	}
}