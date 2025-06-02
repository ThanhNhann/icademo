package keeper_test

// import (
// 	"context"
// 	"time"

// 	"github.com/stretchr/testify/require"

// 	"github.com/cometbft/cometbft/proto/tendermint/crypto"

// 	"github.com/ThanhNhann/icademo/x/interchainquery/types"
// )

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
// 				ConnectionId:     "connection-0",
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
// 				ConnectionId:     "connection-0",
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
// 				ConnectionId:     "connection-0",
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

// 			// Submit the response
// 			_, err := suite.GetMsgServer().SubmitQueryResponse(context.Background(), tc.response)

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
