package storage

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"

	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	"github.com/filecoin-project/lotus/chain/types"
)

// 地址选择器api接口
type addrSelectApi interface {
	// 钱包负载均衡
	WalletBalance(context.Context, address.Address) (types.BigInt, error)
	// 钱包是否存在
	WalletHas(context.Context, address.Address) (bool, error)
	//
	StateAccountKey(context.Context, address.Address, types.TipSetKey) (address.Address, error)
	// 查找id
	StateLookupID(context.Context, address.Address, types.TipSetKey) (address.Address, error)
}

// 地址选择器
type AddressSelector struct {
	api.AddressConfig
}

// 根据地址类型选择post地址 a:地址选择api, mi:矿机信息, use:地址类型, goodFunds:, minFunds:
func (as *AddressSelector) AddressFor(ctx context.Context, a addrSelectApi, mi miner.MinerInfo, use api.AddrUse, goodFunds, minFunds abi.TokenAmount) (address.Address, abi.TokenAmount, error) {
	var addrs []address.Address
	switch use {
	case api.PreCommitAddr: // 0, 预提交地址
		addrs = append(addrs, as.PreCommitControl...)
	case api.CommitAddr: // 1, 提交地址
		addrs = append(addrs, as.CommitControl...)
	default: // 2... post地址
		defaultCtl := map[address.Address]struct{}{}
		for _, a := range mi.ControlAddresses { // post地址
			defaultCtl[a] = struct{}{}
		}
		delete(defaultCtl, mi.Owner)  // 排除 owner地址
		delete(defaultCtl, mi.Worker) // 排除 worker地址
		// 排除预提交控制地址
		for _, addr := range append(append([]address.Address{}, as.PreCommitControl...), as.CommitControl...) {
			if addr.Protocol() != address.ID {
				var err error
				addr, err = a.StateLookupID(ctx, addr, types.EmptyTSK)
				if err != nil {
					log.Warnw("looking up control address", "address", addr, "error", err)
					continue
				}
			}

			delete(defaultCtl, addr)
		}

		for a := range defaultCtl {
			addrs = append(addrs, a)
		}
	}
	addrs = append(addrs, mi.Owner, mi.Worker)

	return pickAddress(ctx, a, mi, goodFunds, minFunds, addrs)
}

// 在地址集合中选择最优地址,返回选中的地址和其token数量; a:地址选择api, mi:矿机信息, goodFunds:最优值, minFunds:最差值, addrs:地址集合
func pickAddress(ctx context.Context, a addrSelectApi, mi miner.MinerInfo, goodFunds, minFunds abi.TokenAmount, addrs []address.Address) (address.Address, abi.TokenAmount, error) {
	leastBad := mi.Worker
	bestAvail := minFunds

	ctl := map[address.Address]struct{}{}
	for _, a := range append(mi.ControlAddresses, mi.Owner, mi.Worker) {
		ctl[a] = struct{}{}
	}

	for _, addr := range addrs {
		if addr.Protocol() != address.ID {
			var err error
			addr, err = a.StateLookupID(ctx, addr, types.EmptyTSK)
			if err != nil {
				log.Warnw("looking up control address", "address", addr, "error", err)
				continue
			}
		}

		if _, ok := ctl[addr]; !ok {
			log.Warnw("non-control address configured for sending messages", "address", addr)
			continue
		}

		if maybeUseAddress(ctx, a, addr, goodFunds, &leastBad, &bestAvail) {
			return leastBad, bestAvail, nil
		}
	}

	log.Warnw("No address had enough funds to for full PoSt message Fee, selecting least bad address", "address", leastBad, "balance", types.FIL(bestAvail), "optimalFunds", types.FIL(goodFunds), "minFunds", types.FIL(minFunds))

	return leastBad, bestAvail, nil
}

// 判断地址能不能用; a:地址选择api, addr:地址, goodFunds:良数, leastBad:, bestAvail:最有用的值
func maybeUseAddress(ctx context.Context, a addrSelectApi, addr address.Address, goodFunds abi.TokenAmount, leastBad *address.Address, bestAvail *abi.TokenAmount) bool {
	b, err := a.WalletBalance(ctx, addr)
	if err != nil {
		log.Errorw("checking control address balance", "addr", addr, "error", err)
		return false
	}

	if b.GreaterThanEqual(goodFunds) {
		k, err := a.StateAccountKey(ctx, addr, types.EmptyTSK)
		if err != nil {
			log.Errorw("getting account key", "error", err)
			return false
		}

		have, err := a.WalletHas(ctx, k)
		if err != nil {
			log.Errorw("failed to check control address", "addr", addr, "error", err)
			return false
		}

		if !have {
			log.Errorw("don't have key", "key", k, "address", addr)
			return false
		}

		*leastBad = addr
		*bestAvail = b
		return true
	}

	if b.GreaterThan(*bestAvail) {
		*leastBad = addr
		*bestAvail = b
	}

	log.Warnw("address didn't have enough funds for window post message", "address", addr, "required", types.FIL(goodFunds), "balance", types.FIL(b))
	return false
}
