package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgRegisterAccount{}, "icademo/txdemo/MsgRegisterAccount", nil)
	cdc.RegisterConcrete(&MsgSubmitTx{}, "icademo/txdemo/MsgSubmitTx", nil)
	cdc.RegisterConcrete(&MsgRegisterHostZone{}, "icademo/txdemo/MsgRegisterHostZone", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgRegisterAccount{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSubmitTx{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgRegisterHostZone{})
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
