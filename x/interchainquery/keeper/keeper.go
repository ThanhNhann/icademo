package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	"github.com/ThanhNhann/icademo/x/interchainquery/types"
)

// Keeper of this module maintains collections of registered zones.
type Keeper struct {
	cdc       codec.Codec
	storeKey  storetypes.StoreKey
	callbacks map[string]types.QueryCallbacks
	IBCKeeper *ibckeeper.Keeper
}

// NewKeeper returns a new instance of zones Keeper
func NewKeeper(cdc codec.Codec, storeKey storetypes.StoreKey, ibcKeeper *ibckeeper.Keeper) Keeper {
	if ibcKeeper == nil {
		panic("ibcKeeper is nil")
	}
	return Keeper{
		cdc:       cdc,
		storeKey:  storeKey,
		callbacks: make(map[string]types.QueryCallbacks),
		IBCKeeper: ibcKeeper,
	}
}

func (k *Keeper) SetCallbackHandler(module string, handler types.QueryCallbacks) error {
	_, found := k.callbacks[module]
	if found {
		return fmt.Errorf("callback handler already set for %s", module)
	}
	k.callbacks[module] = handler
	return nil
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) MakeICQRequest(ctx sdk.Context, query types.Query) error {
	k.Logger(ctx).Info("Making ICQ request- module=%s, callbackId=%s, connectionId=%s, queryType=%s, timeout_duration=%d",
		query.CallbackModule, query.CallbackId, query.ConnectionId, query.QueryType, query.TimeoutDuration)
	
	if err := k.ValidateQuery(ctx, query); err != nil {
		return err
	}

	// Set the timeout using the block time and timeout duration
	// This is int64 + int64, but there is no case for the result < 0
	timeoutTimestamp := (ctx.BlockTime().UnixNano() + query.TimeoutDuration.Nanoseconds())
	query.TimeoutTimestamp = uint64(timeoutTimestamp)
	
	query.Id = GenerateQueryHash(query)
	query.IsSent = false
	query.LastHeight = uint64(ctx.BlockHeight())

	connection, found := k.IBCKeeper.ConnectionKeeper.GetConnection(ctx, query.ConnectionId)
	if !found {
		return fmt.Errorf("connection not found: %s", query.ConnectionId)
	}
	clientState, found := k.IBCKeeper.ClientKeeper.GetClientState(ctx, connection.ClientId)
	if !found {
		return fmt.Errorf("client state not found: %s", connection.ClientId)
	}
	query.LastHeight = clientState.GetLatestHeight().GetRevisionHeight()

	k.SetQuery(ctx, query)

	return nil
}

// Re-submit an ICQ, generally used after a timeout
func (k *Keeper) RetryICQRequest(ctx sdk.Context, query types.Query) error {
	// Delete old query
	k.DeleteQuery(ctx, query.Id)

	// Submit a new query (with a new ID)
	if err := k.MakeICQRequest(ctx, query); err != nil {
		return errorsmod.Wrap(err, types.ErrFailedToRetryQuery.Error())
	}

	return nil
}