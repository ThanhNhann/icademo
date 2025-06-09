package types

import (
	fmt "fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	proto "github.com/gogo/protobuf/proto"

	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
)

var (
	_ sdk.Msg = &MsgRegisterAccount{}
	_ sdk.Msg = &MsgSubmitTx{}
	_ sdk.Msg = &MsgRegisterHostZone{}
)

const TypeMsgRegisterHostZone = "register_host_zone"

// NewMsgRegisterAccount creates a new MsgRegisterAccount instance
func NewMsgRegisterAccount(owner, connectionID, version string) *MsgRegisterAccount {
	return &MsgRegisterAccount{
		Owner:        owner,
		ConnectionId: connectionID,
		Version:      version,
	}
}

// ValidateBasic implements sdk.Msg
func (msg MsgRegisterAccount) ValidateBasic() error {
	if strings.TrimSpace(msg.Owner) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "failed to parse address: %s", msg.Owner)
	}

	return nil
}

// GetSigners implements sdk.Msg
func (msg MsgRegisterAccount) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{accAddr}
}

// NewMsgSubmitTx creates and returns a new MsgSubmitTx instance
func NewMsgSubmitTx(sdkMsg sdk.Msg, connectionID, owner string) (*MsgSubmitTx, error) {
	msg, ok := sdkMsg.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("cannot proto marshal %T", msg)
	}

	protoAny, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	return &MsgSubmitTx{
		ConnectionId: connectionID,
		Owner:        owner,
		Msg:          protoAny,
	}, nil
}

// UnpackInterfaces implements codectypes.UnpackInterfacesMessage
func (msg MsgSubmitTx) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var sdkMsg sdk.Msg

	return unpacker.UnpackAny(msg.Msg, &sdkMsg)
}

// GetTxMsg fetches the cached any message
func (msg *MsgSubmitTx) GetTxMsg() sdk.Msg {
	sdkMsg, ok := msg.Msg.GetCachedValue().(sdk.Msg)
	if !ok {
		return nil
	}

	return sdkMsg
}

// GetSigners implements sdk.Msg
func (msg MsgSubmitTx) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{accAddr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgSubmitTx) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid owner address")
	}

	return nil
}

func NewMsgRegisterHostZone(
	creator string,
	connectionId string,
	bech32prefix string,
	hostDenom string,
	ibcDenom string,
	transferChannelId string,
	unbondingPeriod uint64,
	maxMessagePerIcaTx uint64,
) *MsgRegisterHostZone {
	return &MsgRegisterHostZone{
		Creator:             creator,
		ConnectionId:        connectionId,
		Bech32Prefix:        bech32prefix,
		HostDenom:           hostDenom,
		IbcDenom:            ibcDenom,
		TransferChannelId:   transferChannelId,
		UnbondingPeriod:     unbondingPeriod,
		MaxMessagesPerIcaTx: maxMessagePerIcaTx,
	}
}

func (msg *MsgRegisterHostZone) Route() string {
	return RouterKey
}

func (msg *MsgRegisterHostZone) Type() string {
	return TypeMsgRegisterHostZone
}

func (msg *MsgRegisterHostZone) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRegisterHostZone) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterHostZone) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	// TODO: add a list admin for registerHostZone
	// if err := utils.ValidateAdminAddress(msg.Creator); err != nil {
	// 	return err
	// }
	// VALIDATE DENOMS
	// host denom cannot be empty
	if msg.HostDenom == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "host denom cannot be empty")
	}
	// host denom must be a valid asset denom
	if err := sdk.ValidateDenom(msg.HostDenom); err != nil {
		return err
	}

	// ibc denom cannot be empty and must begin with "ibc"
	if msg.IbcDenom == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "ibc denom cannot be empty")
	}
	if !strings.HasPrefix(msg.IbcDenom, "ibc") {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "ibc denom must begin with 'ibc'")
	}
	// ibc denom must be valid
	err = ibctransfertypes.ValidateIBCDenom(msg.IbcDenom)
	if err != nil {
		return err
	}
	// bech32 prefix must be non-empty (we validate it fully in msg_server)
	if strings.TrimSpace(msg.Bech32Prefix) == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "bech32 prefix must be non-empty")
	}
	// connection id cannot be empty and must begin with "connection"
	if msg.ConnectionId == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "connection id cannot be empty")
	}
	if !strings.HasPrefix(msg.ConnectionId, "connection") {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "connection id must begin with 'connection'")
	}
	// transfer channel id cannot be empty
	if msg.TransferChannelId == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "transfer channel id cannot be empty")
	}
	// transfer channel id must begin with "channel"
	if !strings.HasPrefix(msg.TransferChannelId, "channel") {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "transfer channel id must begin with 'channel'")
	}
	// unbonding frequency must be positive nonzero
	if msg.UnbondingPeriod < 1 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "unbonding frequency must be greater than zero")
	}

	return nil
}
