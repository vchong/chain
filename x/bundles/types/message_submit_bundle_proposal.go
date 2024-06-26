package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ legacytx.LegacyMsg = &MsgSubmitBundleProposal{}
	_ sdk.Msg            = &MsgSubmitBundleProposal{}
)

func NewMsgSubmitBundleProposal(creator string, staker string, poolId uint64, storageId string, dataSize uint64, dataHash string, fromIndex uint64, bundleSize uint64, fromKey string, toKey string, bundleSummary string) *MsgSubmitBundleProposal {
	return &MsgSubmitBundleProposal{
		Creator:       creator,
		Staker:        staker,
		PoolId:        poolId,
		StorageId:     storageId,
		DataSize:      dataSize,
		FromIndex:     fromIndex,
		BundleSize:    bundleSize,
		FromKey:       fromKey,
		ToKey:         toKey,
		BundleSummary: bundleSummary,
		DataHash:      dataHash,
	}
}

func (msg *MsgSubmitBundleProposal) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSubmitBundleProposal) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgSubmitBundleProposal) Route() string {
	return RouterKey
}

func (msg *MsgSubmitBundleProposal) Type() string {
	return "kyve/bundles/MsgSubmitBundleProposal"
}

func (msg *MsgSubmitBundleProposal) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
