package cli

import (
	"encoding/json"
	"fmt"
	"os"
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
	cmd.AddCommand(AddValidatorsCmd())
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

func AddValidatorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-validators [validator-list-file]",
		Short: "Add validators of host chain",
		Long: strings.TrimSpace(
			`Add validators and weights using a JSON file in the following format
	{
		"validator_weights": [
			{"address": "cosmosXXX", "weight": 1},
			{"address": "cosmosXXX", "weight": 2}
		]
	}	
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
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
//		     {"name": "val1", "address": "cosmosXXX", "commission": 0.1},
//			 {"name": "val2", "address": "cosmosXXX", "commission": 0.1}
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
