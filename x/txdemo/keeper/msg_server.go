package keeper

import (
	"context"
	"time"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"github.com/ThanhNhann/icademo/x/txdemo/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var (
	DefaultMaxMessagesPerIcaTx = uint64(32)
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) RegisterAccount(goCtx context.Context, msg *types.MsgRegisterAccount) (*types.MsgRegisterAccountResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	ckMsgServer := icacontrollerkeeper.NewMsgServerImpl(&k.icaControllerKeeper)

	_, err := ckMsgServer.RegisterInterchainAccount(ctx, icacontrollertypes.NewMsgRegisterInterchainAccount(
		msg.ConnectionId,
		msg.Owner,
		msg.Version,
	))
	if err != nil {
		return nil, err
	}

	return &types.MsgRegisterAccountResponse{}, nil
}

func (k msgServer) SubmitTx(goCtx context.Context, msg *types.MsgSubmitTx) (*types.MsgSubmitTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	data, err := icatypes.SerializeCosmosTx(k.cdc, []proto.Message{msg.GetTxMsg()})
	if err != nil {
		return nil, err
	}

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}

	timeoutTimestamp := ctx.BlockTime().Add(time.Minute).UnixNano()
	ckMsgServer := icacontrollerkeeper.NewMsgServerImpl(&k.icaControllerKeeper)
	msgSendTx := icacontrollertypes.NewMsgSendTx(msg.Owner, msg.ConnectionId, uint64(timeoutTimestamp), packetData)
	_, err = ckMsgServer.SendTx(ctx, msgSendTx)
	if err != nil {
		return nil, err
	}

	return &types.MsgSubmitTxResponse{}, nil
}

func (k msgServer) AddValidators(goCtx context.Context, msg *types.MsgAddValidators) (*types.MsgAddValidatorsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	for _, validator := range msg.Validators {
		if err := k.AddValidatorToHostZone(ctx, msg.HostZone, *validator, false); err != nil {
			return nil, err
		}
		// TODO: add icq to tracking status of validator here (slash, tombs toned ?)
	}

	return &types.MsgAddValidatorsResponse{}, nil
}

// createModuleAccount creates a module account for the given address
func (k Keeper) createModuleAccount(ctx sdk.Context, addr sdk.AccAddress) error {
	// Check if account already exists
	acc := k.accountKeeper.GetAccount(ctx, addr)
	if acc != nil {
		return nil // Account already exists
	}

	// Create the module account
	moduleAccount := authtypes.NewEmptyModuleAccount(addr.String())
	k.accountKeeper.SetAccount(ctx, moduleAccount)
	return nil
}

func (k msgServer) RegisterHostZone(goCtx context.Context, msg *types.MsgRegisterHostZone) (*types.MsgRegisterHostZoneResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Get ConnectionEnd (for counterparty connection)
	connectionEnd, found := k.ibcKeeper.ConnectionKeeper.GetConnection(ctx, msg.ConnectionId)
	if !found {
		return nil, errorsmod.Wrapf(connectiontypes.ErrConnectionNotFound, "connection-id %s does not exist", msg.ConnectionId)
	}
	counterpartyConnection := connectionEnd.Counterparty

	// Get chain id from connection
	chainId, err := k.GetChainIdFromConnectionId(ctx, msg.ConnectionId)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "unable to obtain chain id from connection %s", msg.ConnectionId)
	}

	// get zone
	_, found = k.GetHostZone(ctx, chainId)
	if found {
		return nil, errorsmod.Wrapf(types.ErrFailedToRegisterHostZone, "host zone already registered for chain-id %s", chainId)
	}

	// check the denom is not already registered
	hostZones := k.GetAllHostZone(ctx)
	for _, hostZone := range hostZones {
		if hostZone.HostDenom == msg.HostDenom {
			return nil, errorsmod.Wrapf(types.ErrFailedToRegisterHostZone, "host denom %s already registered", msg.HostDenom)
		}
		if hostZone.ConnectionId == msg.ConnectionId {
			return nil, errorsmod.Wrapf(types.ErrFailedToRegisterHostZone, "connection-id %s already registered", msg.ConnectionId)
		}
		if hostZone.TransferChannelId == msg.TransferChannelId {
			return nil, errorsmod.Wrapf(types.ErrFailedToRegisterHostZone, "transfer channel %s already registered", msg.TransferChannelId)
		}
		if hostZone.Bech32Prefix == msg.Bech32Prefix {
			return nil, errorsmod.Wrapf(types.ErrFailedToRegisterHostZone, "bech32 prefix %s already registered", msg.Bech32Prefix)
		}
	}

	// Create module account for the host zone
	moduleName := types.ModuleName + "_" + chainId
	depositAddress := authtypes.NewModuleAddress(moduleName)
	if err := k.createModuleAccount(ctx, depositAddress); err != nil {
		return nil, errorsmod.Wrapf(err, "unable to create deposit account for host zone %s", chainId)
	}

	// Set the max messages per ICA tx to the default value if it's not specified
	maxMessagesPerIcaTx := msg.MaxMessagesPerIcaTx
	if maxMessagesPerIcaTx == 0 {
		maxMessagesPerIcaTx = DefaultMaxMessagesPerIcaTx
	}

	// set the zone
	zone := types.HostZone{
		ChainId:             chainId,
		ConnectionId:        msg.ConnectionId,
		Bech32Prefix:        msg.Bech32Prefix,
		IbcDenom:            msg.IbcDenom,
		HostDenom:           msg.HostDenom,
		TransferChannelId:   msg.TransferChannelId,
		UnbondingPeriod:     msg.UnbondingPeriod,
		DepositAddress:      depositAddress.String(),
		MaxMessagesPerIcaTx: maxMessagesPerIcaTx,
	}
	// write the zone back to the store
	k.SetHostZone(ctx, zone)

	appVersion := string(icatypes.ModuleCdc.MustMarshalJSON(&icatypes.Metadata{
		Version:                icatypes.Version,
		ControllerConnectionId: zone.ConnectionId,
		HostConnectionId:       counterpartyConnection.ConnectionId,
		Encoding:               icatypes.EncodingProtobuf,
		TxType:                 icatypes.TxTypeSDKMultiMsg,
	}))

	// generate delegate account
	delegatePortOwner := chainId + ".delegate"
	if err := k.icaControllerKeeper.RegisterInterchainAccountWithOrdering(ctx, zone.ConnectionId, delegatePortOwner, appVersion, channeltypes.ORDERED); err != nil {
		return nil, errorsmod.Wrap(err, "failed to register delegation ICA")
	}

	// generate withdrawal account
	withdrawalPortOwner := chainId + ".withdrawal"
	if err := k.icaControllerKeeper.RegisterInterchainAccountWithOrdering(ctx, zone.ConnectionId, withdrawalPortOwner, appVersion, channeltypes.ORDERED); err != nil {
		return nil, errorsmod.Wrap(err, "failed to register withdrawal ICA")
	}

	// generate redemption account
	redemptionPortOwner := chainId + ".redemption"
	if err := k.icaControllerKeeper.RegisterInterchainAccountWithOrdering(ctx, zone.ConnectionId, redemptionPortOwner, appVersion, channeltypes.ORDERED); err != nil {
		return nil, errorsmod.Wrap(err, "failed to register redemption ICA")
	}
	
	// TODO: Add a init record for deposite here

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRegisterZone,
			sdk.NewAttribute(types.AttributeKeyConnectionId, msg.ConnectionId),
			sdk.NewAttribute(types.AttributeKeyRecipientChain, chainId),
		),
	)

	return &types.MsgRegisterHostZoneResponse{}, nil
}

