package keeper_test

import (
	"icademo/x/txdemo/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestParamsQuery() {
	k := GetICAApp(suite.chainA).TxdemoKeeper
	ctx := suite.chainA.GetContext()
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	k.SetParams(ctx, params)

	response, err := k.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), &types.QueryParamsResponse{Params: params}, response)
}
