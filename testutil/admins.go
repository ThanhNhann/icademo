package testutil

import (
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)
var Admins = map[string]bool{
	"icademo1k8c2m5cn322akk5wy8lpt87dd2f4yh9azg7jlh": true, // F5
	"icademo10d07y265gmmuvt4z0w9aw880jnsr700japtjqr": true, // gov module
}

func ValidateAdminAddress(address string) error {
	if !Admins[address] {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "address (%s) is not an admin", address)
	}
	return nil
}