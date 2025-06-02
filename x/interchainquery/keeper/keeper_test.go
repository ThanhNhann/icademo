package keeper_test

import (
	"encoding/json"
	"testing"
	// "time"

	// "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	// "github.com/cometbft/cometbft/proto/tendermint/crypto"
	// sdk "github.com/cosmos/cosmos-sdk/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"

	icaapp "github.com/ThanhNhann/icademo/app"
	"github.com/ThanhNhann/icademo/x/interchainquery/keeper"
	"github.com/ThanhNhann/icademo/x/interchainquery/types"
)

const TestOwnerAddress = "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs"

func SetupICATestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	db := dbm.NewMemDB()
	encCdc := icaapp.MakeEncodingConfig()
	app := icaapp.New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, icaapp.DefaultNodeHome, 0, encCdc, icaapp.EmptyAppOptions{})
	return app, icaapp.NewDefaultGenesisState(encCdc.Marshaler)
}

func init() {
	ibctesting.DefaultTestingAppInit = SetupICATestingApp
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
	path   *ibctesting.Path
}

func (suite *KeeperTestSuite) GetSimApp(chain *ibctesting.TestChain) *icaapp.App {
	icademo, ok := chain.App.(*icaapp.App)
	if !ok {
		panic("not icademo app")
	}

	return icademo
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))

	suite.path = newSimAppPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(suite.path)
}

func (s *KeeperTestSuite) GetMsgServer() types.MsgServer {
	return keeper.NewMsgServerImpl(s.GetSimApp(s.chainA).InterchainqueryKeeper)
}

func newSimAppPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}

// func (suite *KeeperTestSuite) TestSubmitQueryResponse() {
// 	testCases := []struct {
// 		name          string
// 		query         types.Query
// 		response      *types.MsgSubmitQueryResponse
// 		expectedError bool
// 	}{
// 		{
// 			name: "non-existent query",
// 			response: &types.MsgSubmitQueryResponse{
// 				QueryId: "non-existent",
// 				ChainId: "chain-1",
// 			},
// 			expectedError: false, // Should return empty response without error
// 		},
// 		{
// 			name: "contentless response",
// 			query: types.Query{
// 				Id:               "test-query-1",
// 				ConnectionId:     suite.path.EndpointA.ConnectionID,
// 				ChainId:          "chain-1",
// 				QueryType:        "test/query",
// 				CallbackModule:   "test",
// 				CallbackId:       "test-callback",
// 				TimeoutTimestamp: uint64(time.Now().Add(time.Hour).UnixNano()),
// 			},
// 			response: &types.MsgSubmitQueryResponse{
// 				QueryId: "test-query-1",
// 				ChainId: "chain-1",
// 				Result:  []byte{},
// 			},
// 			expectedError: false,
// 		},
// 		{
// 			name: "timeout case",
// 			query: types.Query{
// 				Id:               "test-query-2",
// 				ConnectionId:     suite.path.EndpointA.ConnectionID,
// 				ChainId:          "chain-1",
// 				QueryType:        "test/query",
// 				CallbackModule:   "test",
// 				CallbackId:       "test-callback",
// 				TimeoutTimestamp: uint64(time.Now().Add(-time.Hour).UnixNano()), // Expired
// 			},
// 			response: &types.MsgSubmitQueryResponse{
// 				QueryId: "test-query-2",
// 				ChainId: "chain-1",
// 				Result:  []byte("test result"),
// 			},
// 			expectedError: false, // Should retry without error
// 		},
// 		{
// 			name: "key proof query",
// 			query: types.Query{
// 				Id:               "test-query-3",
// 				ConnectionId:     suite.path.EndpointA.ConnectionID,
// 				ChainId:          "chain-1",
// 				QueryType:        "test/key",
// 				CallbackModule:   "test",
// 				CallbackId:       "test-callback",
// 				TimeoutTimestamp: uint64(time.Now().Add(time.Hour).UnixNano()),
// 			},
// 			response: &types.MsgSubmitQueryResponse{
// 				QueryId: "test-query-3",
// 				ChainId: "chain-1",
// 				Result:  []byte("test result"),
// 				ProofOps: &crypto.ProofOps{
// 					Ops: []crypto.ProofOp{
// 						{
// 							Type: "test",
// 							Key:  []byte("test"),
// 							Data: []byte("test"),
// 						},
// 					},
// 				},
// 				Height: 1,
// 			},
// 			expectedError: true, // Should fail due to invalid proof
// 		},
// 	}

// 	for _, tc := range testCases {
// 		suite.Run(tc.name, func() {
// 			// Set up the query if provided
// 			icaapp := suite.GetSimApp(suite.chainA)
// 			if tc.query.Id != "" {
// 				icaapp.InterchainqueryKeeper.SetQuery(suite.chainA.GetContext(), tc.query)
// 			}

// 			// Submit the response with a proper Cosmos SDK context
// 			sdkCtx := suite.chainA.GetContext()
// 			goCtx := sdk.WrapSDKContext(sdkCtx)
// 			_, err := suite.GetMsgServer().SubmitQueryResponse(goCtx, tc.response)

// 			if tc.expectedError {
// 				require.Error(suite.T(), err)
// 			} else {
// 				require.NoError(suite.T(), err)
// 			}

// 			// Verify query is deleted after processing
// 			if tc.query.Id != "" {
// 				_, found := icaapp.InterchainqueryKeeper.GetQuery(suite.chainA.GetContext(), tc.query.Id)
// 				require.False(suite.T(), found)
// 			}
// 		})
// 	}
// }
