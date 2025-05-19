package txdemo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "icademo/testutil/keeper"
	"icademo/testutil/nullify"
	"icademo/x/txdemo"
	"icademo/x/txdemo/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.TxdemoKeeper(t)
	txdemo.InitGenesis(ctx, *k, genesisState)
	got := txdemo.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	// this line is used by starport scaffolding # genesis/test/assert
}
