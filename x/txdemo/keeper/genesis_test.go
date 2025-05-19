package keeper_test

import (
	"icademo/testutil/nullify"
	"icademo/x/txdemo/types"
)

func (suite *KeeperTestSuite) TestGenesis() {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k := GetICAApp(suite.chainA).TxdemoKeeper
	ctx := suite.chainA.GetContext()
	k.InitGenesis(ctx, genesisState)
	got := k.ExportGenesis(ctx)
	suite.Require().NotNil(got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	suite.Require().Equal(genesisState.PortId, got.PortId)

	// this line is used by starport scaffolding # genesis/test/assert
}
