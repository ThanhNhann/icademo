package keeper_test

import (
	"icademo/x/txdemo/keeper"
)

func (suite *KeeperTestSuite) TestMsgServer() {
	msgServer := keeper.NewMsgServerImpl(GetICAApp(suite.chainA).TxdemoKeeper)
	suite.Require().NotNil(msgServer)
}
