package bundles

import (
	"cosmossdk.io/math"
	"github.com/KYVENetwork/chain/util"
	bundlesKeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	"github.com/KYVENetwork/chain/x/bundles/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// Auth
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	// Mint
	mintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
)

func SplitInflation(ctx sdk.Context, k bundlesKeeper.Keeper, bk util.BankKeeper, mk mintKeeper.Keeper, pk types.PoolKeeper, tk types.TeamKeeper, uk util.UpgradeKeeper) {
	minter, err := mk.Minter.Get(ctx)
	if err != nil {
		util.PanicHalt(uk, ctx, "failed to get minter")
	}
	params, err := mk.Params.Get(ctx)
	if err != nil {
		util.PanicHalt(uk, ctx, "failed to get params")
	}

	// get total inflation rewards for current block
	blockProvision := minter.BlockProvision(params).Amount.Int64()

	// calculate the remaining block provision for chain and protocol after x/team took its share
	remainingBlockProvision := blockProvision - tk.GetTeamBlockProvision(ctx)

	// calculate block provision for protocol based on protocol inflation share
	protocolBlockProvision := math.LegacyNewDec(remainingBlockProvision).Mul(pk.GetProtocolInflationShare(ctx)).TruncateInt64()

	if protocolBlockProvision == 0 {
		return
	}

	// track actual distributed block provision for protocol
	distributed := uint64(0)

	// calculate total inflation share weight of pools to get each pool's reward share
	totalInflationShareWeight := math.LegacyZeroDec()

	for _, pool := range pk.GetAllPools(ctx) {
		// only include active pools
		if err := k.AssertPoolCanRun(ctx, pool.Id); err == nil {
			totalInflationShareWeight = totalInflationShareWeight.Add(pool.InflationShareWeight)
		}
	}

	// if the total inflation share weight is zero all rewards go the chain
	if totalInflationShareWeight.IsZero() {
		return
	}

	for _, pool := range pk.GetAllPools(ctx) {
		// only include active pools
		if err := k.AssertPoolCanRun(ctx, pool.Id); err == nil {
			// calculate pool share based of inflation share weight
			amount := uint64(pool.InflationShareWeight.
				Quo(totalInflationShareWeight).
				Mul(math.LegacyNewDec(protocolBlockProvision)).
				TruncateInt64())

			// transfer funds to pool account
			if err := util.TransferFromModuleToAddress(bk, ctx, authTypes.FeeCollectorName, pool.GetPoolAccount().String(), amount); err != nil {
				util.PanicHalt(uk, ctx, err.Error())
			}

			// track transferred $KYVE to protocol
			distributed += amount
		}
	}

	// calculate if a remainder is left
	remainder := uint64(protocolBlockProvision) - distributed

	if remainder > 0 {
		// find an active pool
		for _, pool := range pk.GetAllPools(ctx) {
			if err := k.AssertPoolCanRun(ctx, pool.Id); err != nil {
				// add remainder to first active pool we find
				if err := util.TransferFromModuleToAddress(bk, ctx, authTypes.FeeCollectorName, pool.GetPoolAccount().String(), remainder); err != nil {
					util.PanicHalt(uk, ctx, err.Error())
				}
			}
		}
	}

	// rest gets transferred to chain
	k.Logger().Info("split portion of minted coins to protocol", "amount", protocolBlockProvision)
}
