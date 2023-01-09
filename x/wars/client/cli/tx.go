package cli

import (
	"bufio"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	client2 "github.com/warmage-sports/wars/x/wars/client"
	"github.com/warmage-sports/wars/x/wars/internal/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	warsTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Wars transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	warsTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreateWar(cdc),
		GetCmdEditWar(cdc),
		GetCmdBuy(cdc),
		GetCmdSell(cdc),
		GetCmdSwap(cdc),
		GetCmdMakeOutcomePayment(cdc),
		GetCmdWithdrawShare(cdc),
	)...)

	return warsTxCmd
}

func GetCmdCreateWar(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-war",
		Short: "Create war",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			_token := viper.GetString(FlagToken)
			_name := viper.GetString(FlagName)
			_description := viper.GetString(FlagDescription)
			_functionType := viper.GetString(FlagFunctionType)
			_functionParameters := viper.GetString(FlagFunctionParameters)
			_reserveTokens := viper.GetString(FlagReserveTokens)
			_txFeePercentage := viper.GetString(FlagTxFeePercentage)
			_exitFeePercentage := viper.GetString(FlagExitFeePercentage)
			_feeAddress := viper.GetString(FlagFeeAddress)
			_maxSupply := viper.GetString(FlagMaxSupply)
			_orderQuantityLimits := viper.GetString(FlagOrderQuantityLimits)
			_sanityRate := viper.GetString(FlagSanityRate)
			_sanityMarginPercentage := viper.GetString(FlagSanityMarginPercentage)
			_allowSells := viper.GetBool(FlagAllowSells)
			_signers := viper.GetString(FlagSigners)
			_batchBlocks := viper.GetString(FlagBatchBlocks)
			_outcomePayment := viper.GetString(FlagOutcomePayment)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Parse function parameters
			functionParams, err := client2.ParseFunctionParams(_functionParameters)
			if err != nil {
				return fmt.Errorf(err.Error())
			}

			// Parse reserve tokens
			reserveTokens := strings.Split(_reserveTokens, ",")

			// Parse tx fee percentage
			txFeePercentage, err := sdk.NewDecFromStr(_txFeePercentage)
			if err != nil {
				return sdkerrors.Wrap(types.ErrArgumentMissingOrNonFloat, "tx fee percentage")
			}

			// Parse exit fee percentage
			exitFeePercentage, err := sdk.NewDecFromStr(_exitFeePercentage)
			if err != nil {
				return sdkerrors.Wrap(types.ErrArgumentMissingOrNonFloat, "exit fee percentage")
			}

			// Parse fee address
			feeAddress, err := sdk.AccAddressFromBech32(_feeAddress)
			if err != nil {
				return err
			}

			// Parse max supply
			maxSupply, err := sdk.ParseCoin(_maxSupply)
			if err != nil {
				return err
			}

			// Parse order quantity limits
			orderQuantityLimits, err := sdk.ParseCoins(_orderQuantityLimits)
			if err != nil {
				return err
			}

			// Parse sanity rate
			sanityRate, err := sdk.NewDecFromStr(_sanityRate)
			if err != nil {
				return fmt.Errorf(err.Error())
			}

			// Parse sanity margin percentage
			sanityMarginPercentage, err := sdk.NewDecFromStr(_sanityMarginPercentage)
			if err != nil {
				return fmt.Errorf(err.Error())
			}

			// Parse signers
			signers, err := client2.ParseSigners(_signers)
			if err != nil {
				return err
			}

			// Parse batch blocks
			batchBlocks, err := sdk.ParseUint(_batchBlocks)
			if err != nil {
				return sdkerrors.Wrap(types.ErrArgumentMissingOrNonUInteger, "max batch blocks")
			}

			// Parse order quantity limits
			outcomePayment, err := sdk.ParseCoins(_outcomePayment)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateWar(_token, _name, _description,
				cliCtx.GetFromAddress(), _functionType, functionParams,
				reserveTokens, txFeePercentage, exitFeePercentage, feeAddress,
				maxSupply, orderQuantityLimits, sanityRate, sanityMarginPercentage,
				_allowSells, signers, batchBlocks, outcomePayment)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsWarGeneral)
	cmd.Flags().AddFlagSet(fsWarCreate)

	_ = cmd.MarkFlagRequired(FlagToken)
	_ = cmd.MarkFlagRequired(FlagName)
	_ = cmd.MarkFlagRequired(FlagDescription)
	_ = cmd.MarkFlagRequired(FlagFunctionType)
	_ = cmd.MarkFlagRequired(FlagFunctionParameters)
	_ = cmd.MarkFlagRequired(FlagReserveTokens)
	_ = cmd.MarkFlagRequired(FlagTxFeePercentage)
	_ = cmd.MarkFlagRequired(FlagExitFeePercentage)
	_ = cmd.MarkFlagRequired(FlagFeeAddress)
	_ = cmd.MarkFlagRequired(FlagMaxSupply)
	_ = cmd.MarkFlagRequired(FlagOrderQuantityLimits)
	_ = cmd.MarkFlagRequired(FlagSanityRate)
	_ = cmd.MarkFlagRequired(FlagSanityMarginPercentage)
	_ = cmd.MarkFlagRequired(FlagSigners)
	_ = cmd.MarkFlagRequired(FlagBatchBlocks)
	// _ = cmd.MarkFlagRequired(FlagOutcomePayment) // Optional

	return cmd
}

func GetCmdEditWar(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-war",
		Short: "Edit war",
		RunE: func(cmd *cobra.Command, args []string) error {
			_token := viper.GetString(FlagToken)
			_name := viper.GetString(FlagName)
			_description := viper.GetString(FlagDescription)
			_orderQuantityLimits := viper.GetString(FlagOrderQuantityLimits)
			_sanityRate := viper.GetString(FlagSanityRate)
			_sanityMarginPercentage := viper.GetString(FlagSanityMarginPercentage)
			_signers := viper.GetString(FlagSigners)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Parse signers
			signers, err := client2.ParseSigners(_signers)
			if err != nil {
				return err
			}

			msg := types.NewMsgEditWar(
				_token, _name, _description, _orderQuantityLimits, _sanityRate,
				_sanityMarginPercentage, cliCtx.GetFromAddress(), signers)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsWarGeneral)
	cmd.Flags().AddFlagSet(fsWarEdit)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagToken)
	_ = cmd.MarkFlagRequired(FlagSigners)

	return cmd
}

func GetCmdBuy(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "buy [war-token-with-amount] [max-prices]",
		Example: "" +
			"buy 10abc 1000res1\n" +
			"buy 10abc 1000res1,1000res2",
		Short: "Buy from a war",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			warCoinWithAmount, err := sdk.ParseCoin(args[0])
			if err != nil {
				return err
			}

			maxPrices, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgBuy(cliCtx.GetFromAddress(),
				warCoinWithAmount, maxPrices)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	return cmd
}

func GetCmdSell(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sell [war-token-with-amount]",
		Example: "sell 10abc",
		Short:   "Sell from a war",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			warCoinWithAmount, err := sdk.ParseCoin(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgSell(cliCtx.GetFromAddress(), warCoinWithAmount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	return cmd
}

func GetCmdSwap(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "swap [war-token] [from-amount] [from-token] [to-token]",
		Example: "" +
			"swap abc 100 res1 res2\n" +
			"swap abc 100 res2 res1",
		Short: "Perform a swap between two tokens",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Check that from amount and token can be parsed to a coin
			from, err := client2.ParseTwoPartCoin(args[1], args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgSwap(cliCtx.GetFromAddress(), args[0], from, args[3])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	return cmd
}

func GetCmdMakeOutcomePayment(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "make-outcome-payment [war-token]",
		Example: "make-outcome-payment abc",
		Short:   "Make an outcome payment to a war",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			msg := types.NewMsgMakeOutcomePayment(cliCtx.GetFromAddress(), args[0])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	return cmd
}

func GetCmdWithdrawShare(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw-share [war-token]",
		Example: "withdraw-share abc",
		Short:   "Withdraw share from a war that is in settlement state",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			msg := types.NewMsgWithdrawShare(cliCtx.GetFromAddress(), args[0])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	return cmd
}
