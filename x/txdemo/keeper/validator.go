package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ThanhNhann/icademo/x/txdemo/types"
)

// GetValidator returns validator by chainID and address.
func (k Keeper) GetValidator(ctx sdk.Context, address []byte) (types.Validator, bool) {
	val := types.Validator{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixValidator)
	bz := store.Get(address)
	if len(bz) == 0 {
		return val, false
	}

	k.cdc.MustUnmarshal(bz, &val)
	return val, true
}

// SetValidators set validators.
func (k Keeper) SetValidator(ctx sdk.Context, val types.Validator) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixValidator)
	bz := k.cdc.MustMarshal(&val)
	valAddr, err := GetAddressBytes(val.Address)
	if err != nil {
		return err
	}
	store.Set(valAddr, bz)
	return nil
}

// DeleteValidator delete validator by chainID and address.
func (k Keeper) DeleteValidator(ctx sdk.Context, address []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixValidator)
	store.Delete(address)
}

// IterateZones iterates through zones.
func (k Keeper) IterateValidators(ctx sdk.Context, fn func(index int64, validator types.Validator) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixValidator)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		validator := types.Validator{}
		k.cdc.MustUnmarshal(iterator.Value(), &validator)

		stop := fn(i, validator)

		if stop {
			break
		}
		i++
	}
}

func GetAddressBytes(valoperAddress string) ([]byte, error) {
	_, addr, err := bech32.DecodeAndConvert(valoperAddress)
	return addr, err
}