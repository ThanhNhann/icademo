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

	if err := k.icaControllerKeeper.RegisterInterchainAccount(ctx, msg.ConnectionId, msg.Owner, msg.Version); err != nil {
		return nil, err
	}

	return &types.MsgRegisterAccountResponse{}, nil
}

func (k msgServer) SubmitTx(goCtx context.Context, msg *types.MsgSubmitTx) (*types.MsgSubmitTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	portOwner, err := icatypes.NewControllerPortID(msg.Owner)
	if err != nil {
		return nil, err
	}

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
	msgSendTx := icacontrollertypes.NewMsgSendTx(portOwner, msg.ConnectionId, uint64(timeoutTimestamp), packetData)
	_, err = ckMsgServer.SendTx(ctx, msgSendTx)
	if err != nil {
		return nil, err
	}

	return &types.MsgSubmitTxResponse{}, nil
}
