package storage

import (
	"context"
	"io"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"

	sealing "github.com/filecoin-project/lotus/extern/storage-sealing"
)

// TODO: refactor this to be direct somehow
// 返回矿机地址
func (m *Miner) Address() address.Address {
	return m.sealing.Address()
}

// 添加块到合适的扇区
func (m *Miner) AddPieceToAnySector(ctx context.Context, size abi.UnpaddedPieceSize, r io.Reader, d sealing.DealInfo) (abi.SectorNumber, abi.PaddedPieceSize, error) {
	return m.sealing.AddPieceToAnySector(ctx, size, r, d)
}

// 开始打包扇区
func (m *Miner) StartPackingSector(sectorNum abi.SectorNumber) error {
	return m.sealing.StartPacking(sectorNum)
}

// 列出扇区
func (m *Miner) ListSectors() ([]sealing.SectorInfo, error) {
	return m.sealing.ListSectors()
}

// 获得扇区信息
func (m *Miner) GetSectorInfo(sid abi.SectorNumber) (sealing.SectorInfo, error) {
	return m.sealing.GetSectorInfo(sid)
}

// ???
func (m *Miner) PledgeSector() error {
	return m.sealing.PledgeSector()
}

// 扇区状态更新
func (m *Miner) ForceSectorState(ctx context.Context, id abi.SectorNumber, state sealing.SectorState) error {
	return m.sealing.ForceSectorState(ctx, id, state)
}

// 根据扇区号删除扇区
func (m *Miner) RemoveSector(ctx context.Context, id abi.SectorNumber) error {
	return m.sealing.Remove(ctx, id)
}

// 标记扇区号用来更新
func (m *Miner) MarkForUpgrade(id abi.SectorNumber) error {
	return m.sealing.MarkForUpgrade(id)
}

// 判断扇区号是否标记了更新
func (m *Miner) IsMarkedForUpgrade(id abi.SectorNumber) bool {
	return m.sealing.IsMarkedForUpgrade(id)
}
