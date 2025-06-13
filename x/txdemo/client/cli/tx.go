package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ThanhNhann/icademo/x/txdemo/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
	FlagMaxMessagesPerIcaTx               = "max-messages-per-ica-tx"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(RegisterAccountCmd())
	cmd.AddCommand(SubmitTxCmd())
	cmd.AddCommand(RegisterHostZone())
	cmd.AddCommand(AddValidators())
	// this line is used by starport scaffolding # 1

	return cmd
}

func RegisterAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-ica-account [connection-id] [version]",
		Short: "Register a new ICA account",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRegisterAccount(clientCtx.GetFromAddress().String(), args[0], args[1])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func SubmitTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "submit-tx [path/to/sdk_msg.json] [connection-id]",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			cdc := codec.NewProtoCodec(clientCtx.InterfaceRegistry)

			var txMsg sdk.Msg
			if err := cdc.UnmarshalInterfaceJSON([]byte(args[0]), &txMsg); err != nil {

				// check for file path if JSON input is not provided
				contents, err := os.ReadFile(args[0])
				if err != nil {
					return errors.Wrap(err, "neither JSON input nor path to .json file for sdk msg were provided")
				}

				if err := cdc.UnmarshalInterfaceJSON(contents, &txMsg); err != nil {
					return errors.Wrap(err, "error unmarshalling sdk msg file")
				}
			}

			msg, err := types.NewMsgSubmitTx(txMsg, args[1], clientCtx.GetFromAddress().String())
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func RegisterHostZone() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-host-zone [connection-id] [host-denom] [bech32prefix] [ibc-denom] [channel-id] [unbonding-period]",
		Short: "Broadcast message register-host-zone",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			connectionId := args[0]
			hostDenom := args[1]
			bech32prefix := args[2]
			ibcDenom := args[3]
			channelId := args[4]
			unbondingPeriod, err := strconv.ParseUint(args[5], 10, 64)
			if err != nil {
				return err
			}

			maxMessagesPerIcaTx, err := cmd.Flags().GetUint64(FlagMaxMessagesPerIcaTx)
			if err != nil {
				return err
			}

			msg := types.NewMsgRegisterHostZone(
				clientCtx.GetFromAddress().String(),
				connectionId,
				bech32prefix,
				hostDenom,
				ibcDenom,
				channelId,
				unbondingPeriod,
				maxMessagesPerIcaTx,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().Uint64(FlagMaxMessagesPerIcaTx, 0, "maximum number of ICA txs in a given tx")

	return cmd
}

func AddValidators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-validators [host-zone] [validator-list-file]",
		Short: "Broadcast message add-validators",
		Long: strings.TrimSpace(
			`Add validators and weights using a JSON file in the following format
	{
		"validator_weights": [
			{"address": "cosmosXXX", "weight": 1},
			{"address": "cosmosXXX", "weight": 2}
		]
	}	
`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			hostZone := args[0]
			validatorListProposalFile := args[1]

			validators, err := parseAddValidatorsFile(validatorListProposalFile)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddValidators(
				clientCtx.GetFromAddress().String(),
				hostZone,
				validators.Validators,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

type ValidatorsList struct {
	Validators []*types.Validator `json:"validators,omitempty"`
}

// Parse a JSON with a list of validators in the format
//
//	{
//		  "validators": [
//		     {"name": "val1", "address": "cosmosXXX", "weight": 1},
//			 {"name": "val2", "address": "cosmosXXX", "weight": 2}
//	   ]
//	}
func parseAddValidatorsFile(validatorsFile string) (validators ValidatorsList, err error) {
	fileContents, err := os.ReadFile(validatorsFile)
	if err != nil {
		return validators, err
	}

	if err = json.Unmarshal(fileContents, &validators); err != nil {
		return validators, err
	}

	return validators, nil
}
