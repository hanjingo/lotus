package storage

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/dline"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"

	"github.com/ipfs/go-cid"
)

// SchedulerState defines the possible states in which the scheduler could be,
// for the purposes of journalling.
type SchedulerState string

const (
	// SchedulerStateStarted gets recorded when a WdPoSt cycle for an
	// epoch begins.
	// 时空证明循环过程开始
	SchedulerStateStarted = SchedulerState("started")
	// SchedulerStateAborted gets recorded when a WdPoSt cycle for an
	// epoch is aborted, normally because of a chain reorg or advancement.
	// 时空证明循环过程中断
	SchedulerStateAborted = SchedulerState("aborted")
	// SchedulerStateFaulted gets recorded when a WdPoSt cycle for an
	// epoch terminates abnormally, in which case the error is also recorded.
	// 时空证明循环过程失败
	SchedulerStateFaulted = SchedulerState("faulted")
	// SchedulerStateSucceeded gets recorded when a WdPoSt cycle for an
	// epoch ends successfully.
	// 时空证明循环过程完成
	SchedulerStateSucceeded = SchedulerState("succeeded")
)

// Journal event types.
const (
	evtTypeWdPoStScheduler = iota
	evtTypeWdPoStProofs
	evtTypeWdPoStRecoveries
	evtTypeWdPoStFaults
)

// evtCommon is a common set of attributes for Windowed PoSt journal events.
// 事件公有类
type evtCommon struct {
	Deadline *dline.Info    // 到期信息
	Height   abi.ChainEpoch // 链高度
	TipSet   []cid.Cid      // 标签集合
	Error    error          `json:",omitempty"` // 错误码
}

// WdPoStSchedulerEvt is the journal event that gets recorded on scheduler
// actions.
type WdPoStSchedulerEvt struct {
	evtCommon
	State SchedulerState // 时空证明循环过程状态
}

// WdPoStProofsProcessedEvt is the journal event that gets recorded when
// Windowed PoSt proofs have been processed.
// 证明过程通过事件
type WdPoStProofsProcessedEvt struct {
	evtCommon
	Partitions []miner.PoStPartition
	MessageCID cid.Cid `json:",omitempty"`
}

// WdPoStRecoveriesProcessedEvt is the journal event that gets recorded when
// Windowed PoSt recoveries have been processed.
// 证明过程恢复事件
type WdPoStRecoveriesProcessedEvt struct {
	evtCommon
	Declarations []miner.RecoveryDeclaration
	MessageCID   cid.Cid `json:",omitempty"`
}

// WdPoStFaultsProcessedEvt is the journal event that gets recorded when
// Windowed PoSt faults have been processed.
// 证明过程失败
type WdPoStFaultsProcessedEvt struct {
	evtCommon
	Declarations []miner.FaultDeclaration
	MessageCID   cid.Cid `json:",omitempty"`
}
