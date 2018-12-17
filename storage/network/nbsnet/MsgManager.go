package nbsnet

import (
	"fmt"
	"github.com/NBSChain/go-nbs/storage/network/pb"
	"sync"
)

const (
	NatMsgBase         = net_pb.MsgType_NatBase
	NatBootReg         = net_pb.MsgType_NatBootReg
	NatKeepAlive       = net_pb.MsgType_NatKeepAlive
	NatDigApply        = net_pb.MsgType_NatDigApply
	NatPingPong        = net_pb.MsgType_NatPingPong
	NatDigOut          = net_pb.MsgType_NatDigOut
	NatReversInvite    = net_pb.MsgType_NatReversInvite
	NatReversInviteAck = net_pb.MsgType_NatReversInviteAck
	NatPriDigSyn       = net_pb.MsgType_NatPriDigSyn
	NatPriDigAck       = net_pb.MsgType_NatPriDigAck
	NatDigConfirm      = net_pb.MsgType_NatDigConfirm
	NatBootAnswer      = net_pb.MsgType_NatBootAnswer
	DrainOutOldKa      = net_pb.MsgType_NatInnerBase + 1
	NatEnd             = net_pb.MsgType_NatEnd

	GspBase        = net_pb.MsgType_GspBase
	GspSub         = net_pb.MsgType_GspSub
	GspSubACK      = net_pb.MsgType_GspSubACK
	GspVoteContact = net_pb.MsgType_GspVoteContact
	GspVoteResult  = net_pb.MsgType_GspVoteResult
	GspVoteResAck  = net_pb.MsgType_GspVoteResAck
	GspIntroduce   = net_pb.MsgType_GspIntroduce
	GspWelcome     = net_pb.MsgType_GspWelcome
	GspHeartBeat   = net_pb.MsgType_GspHeartBeat
	GspReplaceArc  = net_pb.MsgType_GspReplaceArc
	GspReplaceAck  = net_pb.MsgType_GspReplaceAck
	GspRemoveIVArc = net_pb.MsgType_GspRemoveIVArc
	GspRemoveOVAcr = net_pb.MsgType_GspRemoveOVAcr
	GspUpdateOVWei = net_pb.MsgType_GspUpdateOVWei
	GspUpdateIVWei = net_pb.MsgType_GspUpdateIVWei
	GspInnerBase   = net_pb.MsgType_GspInnerBase
	GspEnd         = net_pb.MsgType_GspEnd
	MsgPoolSize    = 1 << 12
)

var (
	once          sync.Once
	instance      *dispatcher
	MsgConvertErr = fmt.Errorf("convert to message failed")
)

type Task struct {
	TypId  net_pb.MsgType
	Param  interface{}
	Result chan interface{}
}

type dispatcher struct {
	msgQueue chan *Task
	router   map[net_pb.MsgType]MsgProcess
}

type MsgProcess func(param *Task) error

func GetInstance() *dispatcher {

	once.Do(func() {
		d := &dispatcher{
			msgQueue: make(chan *Task, MsgPoolSize),
			router:   make(map[net_pb.MsgType]MsgProcess),
		}
		instance = d
		go d.runLoop()
	})

	return instance
}

func (d *dispatcher) Register(msgType net_pb.MsgType, handler MsgProcess) {
	d.router[msgType] = handler
}

func (d *dispatcher) Enqueue(param *Task) {
	d.msgQueue <- param
}

func (d *dispatcher) runLoop() {

	for {
		select {

		case task := <-d.msgQueue:

			p, ok := d.router[task.TypId]
			if !ok {
				logger.Warning("unknown task type:->", task.TypId)
				continue
			}

			if err := p(task); err != nil {
				logger.Warning("process task err:->", err)
				continue
			}
		}
	}
}
