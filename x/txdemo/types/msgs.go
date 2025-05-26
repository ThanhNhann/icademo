package types

import (
	fmt "fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	proto "github.com/gogo/protobuf/proto"
)

// txdemo message types
const (
	MsgTypeRegisterAccount = "register_account"
	MsgTypeSubmitTx        = "submit_tx"
	MsgTypeAddValidators   = "add_validators"
)

var (
	_ sdk.Msg = &MsgRegisterAccount{}
	_ sdk.Msg = &MsgSubmitTx{}
	_ sdk.Msg = &MsgAddValidators{}
)

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

func (msg *MsgRegisterAccount) Route() string {
	return RouterKey
}

func (msg *MsgRegisterAccount) Type() string {
	return MsgTypeRegisterAccount
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

// ValidateBasic implements sdk.Msg
func (msg MsgSubmitTx) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid owner address")
	}

	return nil
}

func (msg *MsgSubmitTx) Route() string {
	return RouterKey
}

func (msg *MsgSubmitTx) Type() string {
	return MsgTypeSubmitTx
}

// GetSigners implements sdk.Msg
func (msg MsgSubmitTx) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{accAddr}
}

// NewMsgAddValidators creates a new MsgAddValidators instance
func NewMsgAddValidators(owner string, validators []*Validator) *MsgAddValidators {
	return &MsgAddValidators{
		Owner:      owner,
		Validators: validators,
	}
}

func (msg MsgAddValidators) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid owner address")
	}

	if len(msg.Validators) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "validators cannot be empty")
	}

	for _, validator := range msg.Validators {
		if validator.Name == "" {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "validator name cannot be empty")
		}
		if validator.Address == "" {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "validator address cannot be empty")
		}
	}

	return nil
}

func (msg *MsgAddValidators) Route() string {
	return RouterKey
}

func (msg *MsgAddValidators) Type() string {
	return MsgTypeAddValidators
}

func (msg MsgAddValidators) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{accAddr}
}
