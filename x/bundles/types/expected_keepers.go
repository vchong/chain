package types

import (
	"cosmossdk.io/math"
	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"
	"github.com/KYVENetwork/chain/x/funders/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

type PoolKeeper interface {
	AssertPoolExists(ctx sdk.Context, poolId uint64) error
	GetPoolWithError(ctx sdk.Context, poolId uint64) (pooltypes.Pool, error)
	GetPool(ctx sdk.Context, id uint64) (val pooltypes.Pool, found bool)

	IncrementBundleInformation(ctx sdk.Context, poolId uint64, currentHeight uint64, currentKey string, currentValue string)

	GetAllPools(ctx sdk.Context) (list []pooltypes.Pool)
	ChargeInflationPool(ctx sdk.Context, poolId uint64) (payout uint64, err error)
	GetProtocolInflationShare(ctx sdk.Context) (res math.LegacyDec)
	GetMaxVotingPowerPerPool(ctx sdk.Context) (res math.LegacyDec)
}

type StakerKeeper interface {
	GetAllStakerAddressesOfPool(ctx sdk.Context, poolId uint64) (stakers []string)
	GetCommission(ctx sdk.Context, stakerAddress string) math.LegacyDec
	IncreaseStakerCommissionRewards(ctx sdk.Context, address string, payerModuleName string, amount sdk.Coins) error
	AssertValaccountAuthorized(ctx sdk.Context, poolId uint64, stakerAddress string, valaddress string) error

	DoesStakerExist(ctx sdk.Context, staker string) bool
	DoesValaccountExist(ctx sdk.Context, poolId uint64, stakerAddress string) bool

	LeavePool(ctx sdk.Context, staker string, poolId uint64)

	IncrementPoints(ctx sdk.Context, poolId uint64, stakerAddress string) (newPoints uint64)
	ResetPoints(ctx sdk.Context, poolId uint64, stakerAddress string) (previousPoints uint64)
}

type DelegationKeeper interface {
	GetDelegationAmount(ctx sdk.Context, staker string) uint64
	GetDelegationOfPool(ctx sdk.Context, poolId uint64) uint64
	GetTotalAndHighestDelegationOfPool(ctx sdk.Context, poolId uint64) (uint64, uint64)
	PayoutRewards(ctx sdk.Context, staker string, amount sdk.Coins, payerModuleName string) error
	SlashDelegators(ctx sdk.Context, poolId uint64, staker string, slashType delegationTypes.SlashType)
}

type FundersKeeper interface {
	GetCoinWhitelistMap(ctx sdk.Context) (whitelist map[string]types.WhitelistCoinEntry)
	ChargeFundersOfPool(ctx sdk.Context, poolId uint64, recipient string) (payout sdk.Coins, err error)
}

type TeamKeeper interface {
	GetTeamBlockProvision(ctx sdk.Context) int64
}
