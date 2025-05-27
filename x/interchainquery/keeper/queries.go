package keeper

import (
	"fmt"
	"strings"
	"time"

	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/crypto"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"

	"github.com/ThanhNhann/icademo/x/interchainquery/types"
)

func GenerateQueryHash(query types.Query) string {
	return fmt.Sprintf("%x", crypto.Sha256(append([]byte(query.ConnectionId+query.ChainId+query.CallbackModule+query.QueryType+query.CallbackId), query.Request...)))
}

// GetQuery returns query
func (k Keeper) GetQuery(ctx sdk.Context, id string) (types.Query, bool) {
	query := types.Query{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixQuery)
	bz := store.Get([]byte(id))
	if len(bz) == 0 {
		return query, false
	}
	k.cdc.MustUnmarshal(bz, &query)
	return query, true
}

// SetQuery set query info
func (k Keeper) SetQuery(ctx sdk.Context, query types.Query) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixQuery)
	bz := k.cdc.MustMarshal(&query)
	k.Logger(ctx).Info("Created/updated query", "ID", query.Id)
	store.Set([]byte(query.Id), bz)
}

// DeleteQuery delete query info
func (k Keeper) DeleteQuery(ctx sdk.Context, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixQuery)
	store.Delete([]byte(id))
}

// IterateQueries iterate through querys
func (k Keeper) IterateQueries(ctx sdk.Context, fn func(index int64, queryInfo types.Query) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixQuery)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		query := types.Query{}
		k.cdc.MustUnmarshal(iterator.Value(), &query)
		stop := fn(i, query)

		if stop {
			break
		}
		i++
	}
}

// AllQueries returns every queryInfo in the store
func (k Keeper) AllQueries(ctx sdk.Context) []types.Query {
	querys := []types.Query{}
	k.IterateQueries(ctx, func(_ int64, queryInfo types.Query) (stop bool) {
		querys = append(querys, queryInfo)
		return false
	})
	return querys
}

// ValidateQuery validates that all the required attributes of a query are supplied when submitting an ICQ
func (k Keeper) ValidateQuery(ctx sdk.Context, query types.Query) error {
	if query.ChainId == "" {
		return errorsmod.Wrapf(types.ErrInvalidICQRequest, "chain-id cannot be empty")
	}
	if query.ConnectionId == "" {
		return errorsmod.Wrapf(types.ErrInvalidICQRequest, "connection-id cannot be empty")
	}
	if !strings.HasPrefix(query.ConnectionId, connectiontypes.ConnectionPrefix) {
		return errorsmod.Wrapf(types.ErrInvalidICQRequest, "invalid connection-id (%s)", query.ConnectionId)
	}
	if query.QueryType == "" {
		return errorsmod.Wrapf(types.ErrInvalidICQRequest, "query type cannot be empty")
	}
	if query.CallbackModule == "" {
		return errorsmod.Wrapf(types.ErrInvalidICQRequest, "callback module must be specified")
	}
	if query.CallbackId == "" {
		return errorsmod.Wrapf(types.ErrInvalidICQRequest, "callback-id cannot be empty")
	}
	if query.TimeoutDuration == time.Duration(0) {
		return errorsmod.Wrapf(types.ErrInvalidICQRequest, "timeout duration must be set")
	}
	if _, exists := k.callbacks[query.CallbackModule]; !exists {
		return errorsmod.Wrapf(types.ErrInvalidICQRequest, "no callback handler registered for module (%s)", query.CallbackModule)
	}
	if exists := k.callbacks[query.CallbackModule].HasICQCallback(query.CallbackId); !exists {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "callback-id (%s) is not registered for module (%s)", query.CallbackId, query.CallbackModule)
	}

	return nil
}
