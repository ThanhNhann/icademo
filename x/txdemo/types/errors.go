package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/txdemo module sentinel errors
var (
	ErrSample                   = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrInvalidPacketTimeout     = sdkerrors.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion           = sdkerrors.Register(ModuleName, 1501, "invalid version")
	ErrHostZoneNotFound         = sdkerrors.Register(ModuleName, 1502, "host zone not found")
	ErrHaltedHostZone           = sdkerrors.Register(ModuleName, 1503, "host zone halted")
	ErrClientStateNotTendermint = sdkerrors.Register(ModuleName, 1504, "client state not tendermint")
	ErrFailedToRegisterHostZone = sdkerrors.Register(ModuleName, 1505, "failed to register host zone")
)
