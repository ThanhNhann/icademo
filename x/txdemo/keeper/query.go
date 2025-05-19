package keeper

import (
	"icademo/x/txdemo/types"
)

var _ types.QueryServer = Keeper{}
