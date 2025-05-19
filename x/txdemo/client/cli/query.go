package cli

import (
	"context"
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ThanhNhann/icademo/x/txdemo/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group txdemo queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdQueryInterchainAccount())
	// this line is used by starport scaffolding # 1

	return cmd
}

func CmdQueryInterchainAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "interchain-account [connection-id] [owner]",
		Short: "Query the interchain account address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			connectionID := args[0]
			owner := args[1]

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.InterchainAccount(context.Background(), &types.QueryInterchainAccountRequest{
				ConnectionId: connectionID,
				Owner:        owner,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}
