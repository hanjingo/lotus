package sealing

type SectorState string

// 扇区状态机
var ExistSectorStateList = map[SectorState]struct{}{
	Empty:                {}, // 空,
	WaitDeals:            {}, // 待交易
	Packing:              {}, // 打包中
	GetTicket:            {}, // ？？？
	PreCommit1:           {}, // 预提交1
	PreCommit2:           {}, // 预提交2
	PreCommitting:        {}, // 正在预提交
	PreCommitWait:        {}, // 预提交等待
	WaitSeed:             {}, // ？？？
	Committing:           {}, // 提交中
	SubmitCommit:         {}, // 接受提交
	CommitWait:           {}, // ???
	FinalizeSector:       {}, // 固定扇区
	Proving:              {}, // 证明
	FailedUnrecoverable:  {}, // 无法恢复
	SealPreCommit1Failed: {}, // 密封预提交1失败
	SealPreCommit2Failed: {}, // 密封预提交2失败
	PreCommitFailed:      {}, // 预提交失败
	ComputeProofFailed:   {}, // 计算证明失败
	CommitFailed:         {}, // 提交失败
	PackingFailed:        {}, // 打包失败
	FinalizeFailed:       {}, // 固定扇区失败
	DealsExpired:         {}, // 交易超时
	RecoverDealIDs:       {}, // 恢复交易id
	Faulty:               {}, // 有bug
	FaultReported:        {}, // bug已上报
	FaultedFinal:         {}, // bug修复
	Removing:             {}, // 删除中
	RemoveFailed:         {}, // 删除失败
	Removed:              {}, // 已删除
}

const (
	UndefinedSectorState SectorState = ""

	// happy path
	Empty          SectorState = "Empty"
	WaitDeals      SectorState = "WaitDeals"     // waiting for more pieces (deals) to be added to the sector
	Packing        SectorState = "Packing"       // sector not in sealStore, and not on chain
	GetTicket      SectorState = "GetTicket"     // generate ticket
	PreCommit1     SectorState = "PreCommit1"    // do PreCommit1
	PreCommit2     SectorState = "PreCommit2"    // do PreCommit2
	PreCommitting  SectorState = "PreCommitting" // on chain pre-commit
	PreCommitWait  SectorState = "PreCommitWait" // waiting for precommit to land on chain
	WaitSeed       SectorState = "WaitSeed"      // waiting for seed
	Committing     SectorState = "Committing"    // compute PoRep
	SubmitCommit   SectorState = "SubmitCommit"  // send commit message to the chain
	CommitWait     SectorState = "CommitWait"    // wait for the commit message to land on chain
	FinalizeSector SectorState = "FinalizeSector"
	Proving        SectorState = "Proving"
	// error modes
	FailedUnrecoverable  SectorState = "FailedUnrecoverable"
	SealPreCommit1Failed SectorState = "SealPreCommit1Failed"
	SealPreCommit2Failed SectorState = "SealPreCommit2Failed"
	PreCommitFailed      SectorState = "PreCommitFailed"
	ComputeProofFailed   SectorState = "ComputeProofFailed"
	CommitFailed         SectorState = "CommitFailed"
	PackingFailed        SectorState = "PackingFailed" // TODO: deprecated, remove
	FinalizeFailed       SectorState = "FinalizeFailed"
	DealsExpired         SectorState = "DealsExpired"
	RecoverDealIDs       SectorState = "RecoverDealIDs"

	Faulty        SectorState = "Faulty"        // sector is corrupted or gone for some reason
	FaultReported SectorState = "FaultReported" // sector has been declared as a fault on chain
	FaultedFinal  SectorState = "FaultedFinal"  // fault declared on chain

	Removing     SectorState = "Removing"
	RemoveFailed SectorState = "RemoveFailed"
	Removed      SectorState = "Removed"
)

func toStatState(st SectorState) statSectorState {
	switch st {
	case Empty, WaitDeals, Packing, GetTicket, PreCommit1, PreCommit2, PreCommitting, PreCommitWait, WaitSeed, Committing, SubmitCommit, CommitWait, FinalizeSector:
		return sstSealing
	case Proving, Removed, Removing:
		return sstProving
	}

	return sstFailed
}
