package keeper

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"

	"github.com/ThanhNhann/icademo/x/txdemo/types"
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
		_, err := sdk.ValAddressFromBech32(validator.Address)
		if err != nil {
			return nil, err
		}

		err = k.SetValidator(ctx, *validator)
		if err != nil {
			return nil, err
		}
		// add icq to check validator information
		// icq := icacontrollertypes.NewICQ(validatorAddr, msg.ConnectionId, msg.Version)
		// k.icaControllerKeeper.SetICQ(ctx, icq)
	}

	return &types.MsgAddValidatorsResponse{}, nil
}
