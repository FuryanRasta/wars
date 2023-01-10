package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/mage-war/wars/x/wars/internal/keeper"
	"github.com/mage-war/wars/x/wars/internal/types"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	warsQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Wars querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	warsQueryCmd.AddCommand(flags.GetCommands(
		GetCmdWars(storeKey, cdc),
		GetCmdWar(storeKey, cdc),
		GetCmdBatch(storeKey, cdc),
		GetCmdLastBatch(storeKey, cdc),
		GetCmdCurrentPrice(storeKey, cdc),
		GetCmdCurrentReserve(storeKey, cdc),
		GetCmdCustomPrice(storeKey, cdc),
		GetCmdBuyPrice(storeKey, cdc),
		GetCmdSellReturn(storeKey, cdc),
		GetCmdSwapReturn(storeKey, cdc),
		GetCmdQueryParams(cdc),
	)...)

	return warsQueryCmd
}

func GetCmdWars(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "wars-list",
		Short: "List of all wars",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/wars",
					queryRoute), nil)
			if err != nil {
				fmt.Printf("%s", err.Error())
				return nil
			}

			var out types.QueryWars
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdWar(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "war [war-token]",
		Short: "Query info of a war",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			warToken := args[0]

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/war/%s",
					queryRoute, warToken), nil)
			if err != nil {
				fmt.Printf("%s", err.Error())
				return nil
			}

			var out types.War
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdBatch(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "batch [war-token]",
		Short: "Query info of a war's current batch",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			warToken := args[0]

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/batch/%s",
					queryRoute, warToken), nil)
			if err != nil {
				fmt.Printf("%s", err.Error())
				return nil
			}

			var out types.Batch
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdLastBatch(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "last-batch [war-token]",
		Short: "Query info of a war's last batch",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			warToken := args[0]

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/last_batch/%s",
					queryRoute, warToken), nil)
			if err != nil {
				fmt.Printf("%s", err.Error())
				return nil
			}

			var out types.Batch
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdCurrentPrice(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "current-price [war-token]",
		Short: "Query current price(s) of the war",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			warToken := args[0]

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/current_price/%s",
					queryRoute, warToken), nil)
			if err != nil {
				fmt.Printf("%s", err.Error())
				return nil
			}

			var out sdk.DecCoins
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdCurrentReserve(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "current-reserve [war-token]",
		Example: "current-reserve abc",
		Short:   "Query current balance(s) of the reserve pool",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			warToken := args[0]

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/current_reserve/%s",
					queryRoute, warToken), nil)
			if err != nil {
				fmt.Printf("%s", err.Error())
				return nil
			}

			var out sdk.Coins
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdCustomPrice(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "price [war-token-with-amount]",
		Example: "price 10abc",
		Short:   "Query price(s) of the war at a specific supply",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			warTokenWithAmount := args[0]

			warCoinWithAmount, err := sdk.ParseCoin(warTokenWithAmount)
			if err != nil {
				fmt.Printf("%s", err.Error())
				return nil
			}

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/custom_price/%s/%s",
					queryRoute, warCoinWithAmount.Denom,
					warCoinWithAmount.Amount.String()), nil)
			if err != nil {
				fmt.Printf("%s", err.Error())
				return nil
			}

			var out sdk.DecCoins
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdBuyPrice(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "buy-price [war-token-with-amount]",
		Example: "buy-price 10abc",
		Short:   "Query price(s) of buying an amount of tokens of the war",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			warTokenWithAmount := args[0]

			warCoinWithAmount, err := sdk.ParseCoin(warTokenWithAmount)
			if err != nil {
				fmt.Printf("%s", err.Error())
				return nil
			}

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/buy_price/%s/%s",
					queryRoute, warCoinWithAmount.Denom,
					warCoinWithAmount.Amount.String()), nil)
			if err != nil {
				fmt.Printf("%s", err.Error())
				return nil
			}

			var out types.QueryBuyPrice
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdSellReturn(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "sell-return [war-token-with-amount]",
		Example: "sell-return 10abc",
		Short:   "Query return(s) on selling an amount of tokens of the war",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			warTokenWithAmount := args[0]

			warCoinWithAmount, err := sdk.ParseCoin(warTokenWithAmount)
			if err != nil {
				fmt.Printf("%s", err.Error())
				return nil
			}

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/sell_return/%s/%s",
					queryRoute, warCoinWithAmount.Denom,
					warCoinWithAmount.Amount.String()), nil)
			if err != nil {
				fmt.Printf("%s", err.Error())
				return nil
			}

			var out types.QuerySellReturn
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdSwapReturn(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "swap-return [war-token] [from-token-with-amount] [to-token]",
		Example: "swap-return abc 10res1 res2",
		Short:   "Query return(s) on swapping an amount of tokens to another token",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			warToken := args[0]
			fromTokenWithAmount := args[1]
			toToken := args[2]

			fromCoinWithAmount, err := sdk.ParseCoin(fromTokenWithAmount)
			if err != nil {
				fmt.Printf("%s", err.Error())
				return nil
			}

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/swap_return/%s/%s/%s/%s",
					queryRoute, warToken, fromCoinWithAmount.Denom,
					fromCoinWithAmount.Amount.String(), toToken), nil)
			if err != nil {
				fmt.Printf("%s", err.Error())
				return nil
			}

			var out types.QuerySwapReturn
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdQueryParams implements a command to fetch wars parameters.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current wars parameters",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(`Query genesis parameters for the wars module:

$ <appcli> query wars params
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s",
				types.QuerierRoute, keeper.QueryParams)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(res, &params)
			return cliCtx.PrintOutput(params)
		},
	}
}
