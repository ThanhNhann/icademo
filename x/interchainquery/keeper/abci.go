package keeper

import (
	"encoding/hex"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ThanhNhann/icademo/x/interchainquery/types"
)

// EndBlocker of interchainquery module
func (k Keeper) EndBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	_ = k.Logger(ctx)
	events := sdk.Events{}

	for _, query := range k.AllQueries(ctx) {
		if query.IsSent {
			continue
		}

		k.Logger(ctx).Info("Interchainquery event emitted", "id", query.Id)

		event := sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueQuery),
			sdk.NewAttribute(types.AttributeKeyQueryId, query.Id),
			sdk.NewAttribute(types.AttributeKeyChainId, query.ChainId),
			sdk.NewAttribute(types.AttributeKeyConnectionId, query.ConnectionId),
			sdk.NewAttribute(types.AttributeKeyType, query.QueryType),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatUint(query.LastHeight, 10)),
			sdk.NewAttribute(types.AttributeKeyRequest, hex.EncodeToString(query.Request)),
		)
		event.Type = "query_request"
		events = append(events, event)

		query.IsSent = true
		k.SetQuery(ctx, query)
	}

	if len(events) > 0 {
		ctx.EventManager().EmitEvents(events)
	}
}
