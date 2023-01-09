package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

// MsgCreateWar: Missing arguments

func TestValidateBasicMsgCreateTokenArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.Token = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateNameArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.Name = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateDescriptionArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.Description = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateCreatorMissingGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.Creator = sdk.AccAddress{}

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateReserveTokenArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.ReserveTokens = nil

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgFeeAddressArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.FeeAddress = sdk.AccAddress{}

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgFunctionTypeArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.FunctionType = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgCreateWar: War token denomination

func TestValidateBasicMsgCreateInvalidTokenArgumentGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.Token = "123abc" // starts with number
	err := message.ValidateBasic()
	require.NotNil(t, err)

	message.Token = "a" // too short
	err = message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgCreateWar: Function parameters and function type

func TestValidateBasicMsgCreateMissingFunctionParamGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.FunctionParameters = []FunctionParam{
		message.FunctionParameters[0],
		message.FunctionParameters[1],
	}

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateTypoFunctionParamGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.FunctionParameters = []FunctionParam{
		NewFunctionParam("invalidParam", message.FunctionParameters[0].Value),
		message.FunctionParameters[1],
		message.FunctionParameters[2],
	}

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateNegativeFunctionParamGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.FunctionParameters[0].Value = sdk.NewDec(-1)

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgFunctionTypeArgumentInvalidGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.FunctionType = "invalid_function_type"

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgCreateWar: Reserve tokens

func TestValidateBasicMsgCreateReserveTokenArgumentInvalidGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.ReserveTokens[0] = "123abc" // invalid denomination

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateNoReserveTokensInvalidGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.ReserveTokens = nil

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateReserveTokensWrongAmountInvalidGivesError(t *testing.T) {
	message := newValidMsgCreateSwapperWar()
	message.ReserveTokens = append(message.ReserveTokens, "extra")

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgCreateWar: Max supply validity

func TestValidateBasicMsgCreateInvalidMaxSupplyGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.MaxSupply.Amount = message.MaxSupply.Amount.Neg() // negate

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgCreateWar: Order quantity limits validity

func TestValidateBasicMsgCreateInvalidOrderQuantityLimitGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.OrderQuantityLimits = sdk.NewCoins(sdk.NewCoin("abc", sdk.OneInt()))
	message.OrderQuantityLimits[0].Amount = message.OrderQuantityLimits[0].Amount.Neg()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgCreateWar: Max supply denom matches war token denom

func TestValidateBasicMsgCreateMaxSupplyDenomTokenDenomMismatchGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.Token = message.MaxSupply.Denom + "a" // to ensure different

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgCreateWar: Sanity values must be positive

func TestValidateBasicMsgCreateNegativeSanityRateGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.SanityRate = sdk.OneDec().Neg()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateNegativeSanityPercentageGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.SanityMarginPercentage = sdk.OneDec().Neg()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgCreateWar: Fee percentages must be positive and not add up to 100

func TestValidateBasicMsgCreateTxFeeIsNegativeGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.TxFeePercentage = sdk.NewDec(-1)

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateTxFeeIsZeroGivesNoError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.TxFeePercentage = sdk.ZeroDec()

	err := message.ValidateBasic()
	require.Nil(t, err)
}

func TestValidateBasicMsgCreateExitFeeIsNegativeGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.ExitFeePercentage = sdk.NewDec(-1)

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateExitFeeIsZeroGivesNoError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.ExitFeePercentage = sdk.ZeroDec()

	err := message.ValidateBasic()
	require.Nil(t, err)
}

func TestValidateBasicMsgCreate100PercentFeeGivesError(t *testing.T) {
	message := newValidMsgCreateWar()

	message.TxFeePercentage = sdk.NewDec(100)
	message.ExitFeePercentage = sdk.ZeroDec()
	err := message.ValidateBasic()
	require.NotNil(t, err)

	message.TxFeePercentage = sdk.NewDec(50)
	message.ExitFeePercentage = sdk.NewDec(50)
	err = message.ValidateBasic()
	require.NotNil(t, err)

	message.TxFeePercentage = sdk.ZeroDec()
	message.ExitFeePercentage = sdk.NewDec(100)
	err = message.ValidateBasic()
	require.NotNil(t, err)

	message.TxFeePercentage = sdk.MustNewDecFromStr("49.999999")
	message.ExitFeePercentage = sdk.NewDec(50)
	require.Nil(t, message.ValidateBasic())

	message.TxFeePercentage = sdk.NewDec(50)
	message.ExitFeePercentage = sdk.MustNewDecFromStr("49.999999")
	require.Nil(t, message.ValidateBasic())
}

// MsgCreateWar: Batch blocks and max supply cannot be zero

func TestValidateBasicMsgCreateZeroBatchBlocksGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.BatchBlocks = sdk.ZeroUint()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgCreateZeroMaxSupplyGivesError(t *testing.T) {
	message := newValidMsgCreateWar()
	message.MaxSupply = sdk.NewCoin(token, sdk.ZeroInt())

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgCreateWar: Valid war creation

func TestValidateBasicMsgCreateWarCorrectlyGivesNoError(t *testing.T) {
	message := newValidMsgCreateWar()

	err := message.ValidateBasic()
	require.Nil(t, err)
}

// MsgEditWar: missing arguments

func TestValidateBasicMsgEditWarTokenArgumentMissingGivesError(t *testing.T) {
	message := newEmptyStringsMsgEditWar()
	message.Token = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgEditWarNameArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgEditWar()
	message.Name = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgEditWarDescriptionArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgEditWar()
	message.Description = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgEditWarOrderQuantityLimitsArgumentMissingGivesNoError(t *testing.T) {
	message := newValidMsgEditWar()
	message.OrderQuantityLimits = ""

	err := message.ValidateBasic()
	require.Nil(t, err)
}

func TestValidateBasicMsgEditWarSanityRateArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgEditWar()
	message.SanityRate = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgEditWarSanityMarginPercentageArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgEditWar()
	message.SanityMarginPercentage = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgEditWarEditorArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgEditWar()
	message.Editor = nil

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgEditWar: no edits

func TestValidateBasicMsgEditWarNoEditsGivesError(t *testing.T) {
	message := NewMsgEditWar(DoNotModifyField, DoNotModifyField,
		DoNotModifyField, DoNotModifyField, DoNotModifyField,
		DoNotModifyField, initCreator, initSigners)

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgEditWar: correct edit

func TestValidateBasicMsgEditWarCorrectlyGivesNoError(t *testing.T) {
	message := newValidMsgEditWar()

	err := message.ValidateBasic()
	require.Nil(t, err)
}

// MsgBuy: missing arguments

func TestValidateBasicMsgBuyBuyerArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgBuy()
	message.Buyer = sdk.AccAddress{}

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgBuy: invalid arguments

func TestValidateBasicMsgBuyInvalidAmountGivesError(t *testing.T) {
	message := newValidMsgBuy()
	message.Amount.Amount = message.Amount.Amount.Neg()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgBuyZeroAmountGivesError(t *testing.T) {
	message := newValidMsgBuy()
	message.Amount.Amount = sdk.ZeroInt()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgBuyMaxPricesInvalidGivesError(t *testing.T) {
	message := newValidMsgBuy()
	message.MaxPrices[0].Amount = sdk.ZeroInt()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgBuy: correct buy

func TestValidateBasicMsgBuyCorrectlyGivesNoError(t *testing.T) {
	message := newValidMsgBuy()

	err := message.ValidateBasic()
	require.Nil(t, err)
}

// MsgSell: missing arguments

func TestValidateBasicMsgSellSellerArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgSell()
	message.Seller = sdk.AccAddress{}

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgSell: invalid arguments

func TestValidateBasicMsgSellInvalidAmountGivesError(t *testing.T) {
	message := newValidMsgSell()
	message.Amount.Amount = message.Amount.Amount.Neg()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgSellZeroAmountGivesError(t *testing.T) {
	message := newValidMsgSell()
	message.Amount.Amount = sdk.ZeroInt()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgSell: correct sell

func TestValidateBasicMsgSellCorrectlyGivesNoError(t *testing.T) {
	message := newValidMsgSell()

	err := message.ValidateBasic()
	require.Nil(t, err)
}

// MsgSwap: missing arguments

func TestValidateBasicMsgSwapSwapperArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.Swapper = sdk.AccAddress{}

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgSwapWarTokenArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.WarToken = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgSwapToTokenArgumentMissingGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.ToToken = ""

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgSwap: invalid arguments

func TestValidateBasicMsgSwapInvalidFromAmountGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.From.Amount = message.From.Amount.Neg()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgSwapInvalidToTokenGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.ToToken = "123abc"

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

func TestValidateBasicMsgSwapZeroFromAmountGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.From.Amount = sdk.ZeroInt()

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgSwap: fromToken==toToken

func TestValidateBasicMsgSwapFromAndToSameTokenGivesError(t *testing.T) {
	message := newValidMsgSwap()
	message.From = sdk.NewInt64Coin(reserveToken, 10)
	message.ToToken = message.From.Denom

	err := message.ValidateBasic()
	require.NotNil(t, err)
}

// MsgSwap: correct swap

func TestValidateBasicMsgSwapCorrectlyGivesNoError(t *testing.T) {
	message := newValidMsgSwap()

	err := message.ValidateBasic()
	require.Nil(t, err)
}
