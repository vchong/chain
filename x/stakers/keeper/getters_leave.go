package keeper

import (
	"cosmossdk.io/store/prefix"
	storeTypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// #####################
// === QUEUE ENTRIES ===
// #####################

// SetLeavePoolEntry ...
func (k Keeper) SetLeavePoolEntry(ctx sdk.Context, leavePoolEntry types.LeavePoolEntry) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.LeavePoolEntryKeyPrefix)
	b := k.cdc.MustMarshal(&leavePoolEntry)
	store.Set(types.LeavePoolEntryKey(
		leavePoolEntry.Index,
	), b)

	// Insert the same entry with a different key prefix for query lookup
	indexStore := prefix.NewStore(storeAdapter, types.LeavePoolEntryKeyPrefixIndex2)
	indexStore.Set(types.LeavePoolEntryKeyIndex2(
		leavePoolEntry.Staker,
		leavePoolEntry.PoolId,
	), []byte{1})
}

// GetLeavePoolEntry ...
func (k Keeper) GetLeavePoolEntry(ctx sdk.Context, index uint64) (val types.LeavePoolEntry, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.LeavePoolEntryKeyPrefix)

	b := store.Get(types.LeavePoolEntryKey(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetLeavePoolEntryByIndex2 ...
func (k Keeper) GetLeavePoolEntryByIndex2(ctx sdk.Context, staker string, poolId uint64) (val types.LeavePoolEntry, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.LeavePoolEntryKeyPrefixIndex2)

	b := store.Get(types.LeavePoolEntryKeyIndex2(staker, poolId))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// DoesLeavePoolEntryExistByIndex2 ...
func (k Keeper) DoesLeavePoolEntryExistByIndex2(ctx sdk.Context, staker string, poolId uint64) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.LeavePoolEntryKeyPrefixIndex2)

	return store.Has(types.LeavePoolEntryKeyIndex2(staker, poolId))
}

// RemoveLeavePoolEntry ...
func (k Keeper) RemoveLeavePoolEntry(ctx sdk.Context, leavePoolEntry *types.LeavePoolEntry) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.LeavePoolEntryKeyPrefix)
	store.Delete(types.LeavePoolEntryKey(leavePoolEntry.Index))

	indexStore := prefix.NewStore(storeAdapter, types.LeavePoolEntryKeyPrefixIndex2)
	indexStore.Delete(types.LeavePoolEntryKeyIndex2(
		leavePoolEntry.Staker,
		leavePoolEntry.PoolId,
	))
}

// GetAllLeavePoolEntries ...
func (k Keeper) GetAllLeavePoolEntries(ctx sdk.Context) (list []types.LeavePoolEntry) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.LeavePoolEntryKeyPrefix)
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.LeavePoolEntry
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
