package cli

import (
	"github.com/warmage-sports/wars/x/wars/internal/types"
	flag "github.com/spf13/pflag"
)

const (
	FlagToken                  = "token"
	FlagName                   = "name"
	FlagDescription            = "description"
	FlagFunctionType           = "function-type"
	FlagFunctionParameters     = "function-parameters"
	FlagReserveTokens          = "reserve-tokens"
	FlagTxFeePercentage        = "tx-fee-percentage"
	FlagExitFeePercentage      = "exit-fee-percentage"
	FlagFeeAddress             = "fee-address"
	FlagMaxSupply              = "max-supply"
	FlagOrderQuantityLimits    = "order-quantity-limits"
	FlagSanityRate             = "sanity-rate"
	FlagSanityMarginPercentage = "sanity-margin-percentage"
	FlagAllowSells             = "allow-sells"
	FlagSigners                = "signers"
	FlagBatchBlocks            = "batch-blocks"
	FlagOutcomePayment         = "outcome-payment"
)

var (
	fsWarGeneral = flag.NewFlagSet("", flag.ContinueOnError)
	fsWarCreate  = flag.NewFlagSet("", flag.ContinueOnError)
	fsWarEdit    = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {

	fsWarGeneral.String(FlagToken, "", "The war's token")
	fsWarGeneral.String(FlagSigners, "", "The list of signers required to create/edit the war")

	fsWarCreate.String(FlagName, "", "The war's name")
	fsWarCreate.String(FlagDescription, "", "The war's description")
	fsWarCreate.String(FlagFunctionType, "", "The type of function that the war will be")
	fsWarCreate.String(FlagFunctionParameters, "", "The parameters that will define the function")
	fsWarCreate.String(FlagReserveTokens, "", "The token(s) that will serve as the reserve token(s)")
	fsWarCreate.String(FlagTxFeePercentage, "", "The percentage fee charged on buys and sells")
	fsWarCreate.String(FlagExitFeePercentage, "", "The percentage fee charged on sells")
	fsWarCreate.String(FlagFeeAddress, "", "The address that will hold any charged fees")
	fsWarCreate.String(FlagMaxSupply, "", "The maximum supply that can be achieved")
	fsWarCreate.String(FlagOrderQuantityLimits, "", "The max number of tokens bought/sold/swapped per order")
	fsWarCreate.String(FlagSanityRate, "", "For swappers, this is the typical t1 per t2 rate")
	fsWarCreate.String(FlagSanityMarginPercentage, "", "For swappers, this is the acceptable deviation from the sanity rate")
	fsWarCreate.Bool(FlagAllowSells, false, "Whether or not sells will be allowed")
	fsWarCreate.String(FlagBatchBlocks, "", "The duration in terms of blocks of each orders batch")
	fsWarCreate.String(FlagOutcomePayment, "", "The payment that would be required to transition the war to settlement")

	fsWarEdit.String(FlagName, types.DoNotModifyField, "The war's name")
	fsWarEdit.String(FlagDescription, types.DoNotModifyField, "The war's description")
	fsWarEdit.String(FlagOrderQuantityLimits, types.DoNotModifyField, "The max number of tokens bought/sold/swapped per order")
	fsWarEdit.String(FlagSanityRate, types.DoNotModifyField, "For swappers, this is the typical t1 per t2 rate")
	fsWarEdit.String(FlagSanityMarginPercentage, types.DoNotModifyField, "For swappers, this is the acceptable deviation from the sanity rate")
}
