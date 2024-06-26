package keeper

import (
	"context"
	"fmt"

	"github.com/KYVENetwork/chain/x/funders/types"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// DefundPool handles the logic to defund a pool.
// If the user is a funder, it will subtract the provided amount
// and send the tokens back. If there are no more funds left, the funding will get inactive.
func (k msgServer) DefundPool(goCtx context.Context, msg *types.MsgDefundPool) (*types.MsgDefundPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Funding has to exist
	funding, found := k.GetFunding(ctx, msg.Creator, msg.PoolId)
	if !found {
		return nil, errors.Wrapf(errorsTypes.ErrNotFound, types.ErrFundingDoesNotExist.Error(), msg.PoolId, msg.Creator)
	}

	// FundingState has to exist
	fundingState, found := k.GetFundingState(ctx, msg.PoolId)
	if !found {
		util.PanicHalt(k.upgradeKeeper, ctx, fmt.Sprintf("FundingState for pool %d does not exist", msg.PoolId))
	}

	// If funder defunds more than he has we defund the entire amount of that coin
	defundAmounts := funding.Amounts.Min(msg.Amounts)
	if defundAmounts.IsZero() {
		return nil, errors.Wrapf(errorsTypes.ErrInvalidRequest, types.ErrFundsTooLow.Error())
	}

	// Subtract amount from funding
	funding.Amounts = funding.Amounts.Sub(defundAmounts...)

	if err := k.ensureParamsCompatibility(ctx, &funding); err != nil {
		return nil, err
	}

	if funding.Amounts.IsZero() {
		fundingState.SetInactive(&funding)
	}

	// Transfer tokens from this module to sender.
	recipient := sdk.MustAccAddressFromBech32(msg.Creator)
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, defundAmounts); err != nil {
		return nil, err
	}

	// Save funding and funding state
	k.SetFunding(ctx, &funding)
	k.SetFundingState(ctx, &fundingState)

	// Emit a defund event.
	_ = ctx.EventManager().EmitTypedEvent(&types.EventDefundPool{
		PoolId:  msg.PoolId,
		Address: msg.Creator,
		Amounts: defundAmounts.String(),
	})

	return &types.MsgDefundPoolResponse{}, nil
}
