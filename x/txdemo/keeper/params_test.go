package keeper_test

import (
	"icademo/x/txdemo/types"
)

func (suite *KeeperTestSuite) TestGetParams() {
	k := GetICAApp(suite.chainA).TxdemoKeeper
	ctx := suite.chainA.GetContext()
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	suite.Require().EqualValues(params, k.GetParams(ctx))
}
