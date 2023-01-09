package wars_test

import (
	"github.com/warmage-sports/wars/x/wars"
	"github.com/warmage-sports/wars/x/wars/internal/types"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
)

func TestInvalidMsgFails(t *testing.T) {
	_, ctx := createTestApp(false)
	h := wars.NewHandler(wars.Keeper{})

	msg := sdk.NewTestMsg()
	_, err := h(ctx, msg)

	require.Error(t, err)
}

func TestCreateValidWar(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war
	_, err := h(ctx, newValidMsgCreateWar())

	require.NoError(t, err)
	require.True(t, app.WarsKeeper.WarExists(ctx, token))

	// Check assigned initial state
	war := app.WarsKeeper.MustGetWar(ctx, token)
	require.Equal(t, types.OpenState, war.State)
}

func TestCreateValidAugmentedWarHatchState(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create augmented function war
	_, err := h(ctx, newValidMsgCreateAugmentedWar())

	require.NoError(t, err)
	require.True(t, app.WarsKeeper.WarExists(ctx, token))

	// Check initial state (hatch since augmented)
	war := app.WarsKeeper.MustGetWar(ctx, token)
	require.Equal(t, types.HatchState, war.State)

	// Check function params (R0, S0, V0 added)
	paramsMap := war.FunctionParameters.AsMap()
	d0, _ := paramsMap["d0"]
	p0, _ := paramsMap["p0"]
	theta, _ := paramsMap["theta"]
	kappa, _ := paramsMap["kappa"]

	initialParams := functionParametersAugmented().AsMap()
	require.Equal(t, d0, initialParams["d0"])
	require.Equal(t, p0, initialParams["p0"])
	require.Equal(t, theta, initialParams["theta"])
	require.Equal(t, kappa, initialParams["kappa"])

	R0 := d0.Mul(sdk.OneDec().Sub(theta))
	S0 := d0.Quo(p0)
	V0 := types.Invariant(R0, S0, kappa.TruncateInt64())

	require.Equal(t, R0, paramsMap["R0"])
	require.Equal(t, S0, paramsMap["S0"])
	require.Equal(t, V0, paramsMap["V0"])
	require.Len(t, war.FunctionParameters, 7)
}

func TestCreateWarThatAlreadyExistsFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Set war to simulate creation
	war := types.War{Token: token}
	app.WarsKeeper.SetWar(ctx, token, war)

	// Create war with same token
	_, err := h(ctx, newValidMsgCreateWar())

	require.Error(t, err)
}

func TestCreatingAWarUsingStakingTokenFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war with token set to staking token
	msg := newValidMsgCreateWar()
	msg.Token = app.StakingKeeper.GetParams(ctx).WarDenom
	_, err := h(ctx, msg)

	require.Error(t, err)
	require.False(t, app.WarsKeeper.WarExists(ctx, token))
}

func TestEditingANonExistingWarFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Edit war
	msg := types.NewMsgEditWar(token, initName, initDescription, "",
		"0", "0", initCreator, initSigners)
	_, err := h(ctx, msg)

	require.Error(t, err)
	require.False(t, app.WarsKeeper.WarExists(ctx, token))
}

func TestEditingAWarWithDifferentCreatorAndSignersFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Set war to simulate creation
	app.WarsKeeper.SetWar(ctx, token, newSimpleWar())

	// Edit war
	msg := types.NewMsgEditWar(token, initName, initDescription, "",
		"0", "0", initCreator, []sdk.AccAddress{anotherAddress})
	_, err := h(ctx, msg)

	require.Error(t, err)
}

func TestEditingAWarWithNegativeOrderQuantityLimitsFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Set war to simulate creation
	app.WarsKeeper.SetWar(ctx, token, newSimpleWar())

	// Edit war
	msg := types.NewMsgEditWar(token, initName, initDescription, "-10testtoken",
		"0", "0", initCreator, initSigners)
	_, err := h(ctx, msg)

	require.Error(t, err)
}

func TestEditingAWarWithFloatOrderQuantityLimitsFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Set war to simulate creation
	app.WarsKeeper.SetWar(ctx, token, newSimpleWar())

	// Edit war
	msg := types.NewMsgEditWar(token, initName, initDescription, "10.5testtoken",
		"0", "0", initCreator, initSigners)
	_, err := h(ctx, msg)

	require.Error(t, err)
}

func TestEditingAWarWithSanityRateEmptyStringMakesSanityFieldsZero(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Set war to simulate creation
	war := newSimpleWar()
	war.SanityRate = sdk.OneDec()
	war.SanityMarginPercentage = sdk.OneDec()
	app.WarsKeeper.SetWar(ctx, token, war)

	// Check sanity values before
	war, _ = app.WarsKeeper.GetWar(ctx, token)
	require.NotEqual(t, sdk.ZeroDec(), war.SanityRate)
	require.NotEqual(t, sdk.ZeroDec(), war.SanityMarginPercentage)

	// Edit war
	msg := types.NewMsgEditWar(token, initName, initDescription, "10testtoken",
		"", "", initCreator, initSigners)
	_, err := h(ctx, msg)

	// Check sanity values after
	require.NoError(t, err)
	war, _ = app.WarsKeeper.GetWar(ctx, token)
	require.Equal(t, sdk.ZeroDec(), war.SanityRate)
	require.Equal(t, sdk.ZeroDec(), war.SanityMarginPercentage)
}

func TestEditingAWarWithNegativeSanityRateFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Set war to simulate creation
	app.WarsKeeper.SetWar(ctx, token, newSimpleWar())

	// Edit war
	msg := types.NewMsgEditWar(token, initName, initDescription, "10testtoken",
		"-10", "", initCreator, initSigners)
	_, err := h(ctx, msg)

	require.Error(t, err)
}

func TestEditingAWarWithNonFloatSanityRateFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Set war to simulate creation
	app.WarsKeeper.SetWar(ctx, token, newSimpleWar())

	// Edit war
	msg := types.NewMsgEditWar(token, initName, initDescription, "10testtoken",
		"20t", "", initCreator, initSigners)
	_, err := h(ctx, msg)

	require.Error(t, err)
}

func TestEditingAWarWithNegativeSanityMarginPercentageFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Set war to simulate creation
	app.WarsKeeper.SetWar(ctx, token, newSimpleWar())

	// Edit war
	msg := types.NewMsgEditWar(token, initName, initDescription, "10testtoken",
		"10", "-5", initCreator, initSigners)
	_, err := h(ctx, msg)

	require.Error(t, err)
}

func TestEditingAWarWithNonFloatSanityMarginPercentageFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Set war to simulate creation
	app.WarsKeeper.SetWar(ctx, token, newSimpleWar())

	// Edit war
	msg := types.NewMsgEditWar(token, initName, initDescription, "10testtoken",
		"20", "20t", initCreator, initSigners)
	_, err := h(ctx, msg)

	require.Error(t, err)
}

func TestEditingAWarCorrectlyPasses(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Set war to simulate creation
	app.WarsKeeper.SetWar(ctx, token, newSimpleWar())

	// Edit war
	newName := "a new name"
	newDescription := "a new description"
	msg := types.NewMsgEditWar(token, newName, newDescription, "",
		"0", "0", initCreator, initSigners)
	_, err := h(ctx, msg)

	require.NoError(t, err)
	war, _ := app.WarsKeeper.GetWar(ctx, token)
	require.Equal(t, newName, war.Name)
	require.Equal(t, newDescription, war.Description)
	require.Equal(t, sdk.Coins(nil), war.OrderQuantityLimits)
	require.Equal(t, sdk.ZeroDec(), war.SanityRate)
	require.Equal(t, sdk.ZeroDec(), war.SanityMarginPercentage)
}

func TestBuyingANonExistingWarFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Buy 1 token
	_, err := h(ctx, newValidMsgBuy(1, 10))

	require.Error(t, err)
}

func TestBuyingAWarWithNonExistentToken(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war
	h(ctx, newValidMsgCreateWar())

	// Buy tokens of another war
	msg := newValidMsgBuy(amountLTMaxSupply, 0) // 0 max prices replaced below
	msg.MaxPrices = sdk.Coins{sdk.NewInt64Coin(token2, 10)}
	_, err := h(ctx, msg)

	require.Error(t, err)
	currentSupply := app.WarsKeeper.MustGetWar(ctx, token).CurrentSupply
	require.True(t, currentSupply.Amount.IsZero())
}

func TestBuyingAWarWithMaxPriceBiggerThanBalanceFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war
	h(ctx, newValidMsgCreateWar())

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 4000)})
	require.Nil(t, err)

	// Buy 10 tokens
	_, err = h(ctx, newValidMsgBuy(10, 5000))

	require.Error(t, err)
	currentSupply := app.WarsKeeper.MustGetWar(ctx, token).CurrentSupply
	require.True(t, currentSupply.Amount.IsZero())
}

func TestBuyingWarWithOrderQuantityLimitExceededFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war with order quantity limit
	msg := newValidMsgCreateWar()
	msg.OrderQuantityLimits = sdk.NewCoins(sdk.NewInt64Coin(token, 4))
	h(ctx, msg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 4000)})
	require.Nil(t, err)

	// Buy 10 tokens
	_, err = h(ctx, newValidMsgBuy(10, 4000))

	require.Error(t, err)
	currentSupply := app.WarsKeeper.MustGetWar(ctx, token).CurrentSupply
	require.True(t, currentSupply.Amount.IsZero())
}

func TestBuyingAWarExceedingMaxSupplyFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war
	h(ctx, newValidMsgCreateWar())

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 6000)})
	require.Nil(t, err)

	// Buy an amount greater than max supply
	_, err = h(ctx, newValidMsgBuy(amountGTMaxSupply, 10))

	require.Error(t, err)
	currentSupply := app.WarsKeeper.MustGetWar(ctx, token).CurrentSupply
	require.True(t, currentSupply.Amount.IsZero())
}

func TestBuyingAWarExceedingMaxPriceFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war
	h(ctx, newValidMsgCreateWar())

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 6000)})
	require.Nil(t, err)

	// Buy an amount less than max supply but with low max prices
	_, err = h(ctx, newValidMsgBuy(amountLTMaxSupply, 1))

	require.Error(t, err)
	currentSupply := app.WarsKeeper.MustGetWar(ctx, token).CurrentSupply
	require.True(t, currentSupply.Amount.IsZero())
}

func TestBuyingAWarWithoutSufficientFundsFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war
	h(ctx, newValidMsgCreateWar())

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 4000)})
	require.Nil(t, err)

	// Buy 10 tokens
	_, err = h(ctx, newValidMsgBuy(10, 4000))

	require.Error(t, err)
	currentSupply := app.WarsKeeper.MustGetWar(ctx, token).CurrentSupply
	require.True(t, currentSupply.Amount.IsZero())
}

func TestBuyingAWarWithoutSufficientFundsDueToTxFeeFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war
	h(ctx, newValidMsgCreateWar())

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 5000)})
	require.Nil(t, err)

	// Buy 10 tokens
	_, err = h(ctx, newValidMsgBuy(10, 5000))

	require.Error(t, err)
	currentSupply := app.WarsKeeper.MustGetWar(ctx, token).CurrentSupply
	require.True(t, currentSupply.Amount.IsZero())
}

func TestBuyingAWarCorrectlyPasses(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war
	h(ctx, newValidMsgCreateWar())

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 4000)})
	require.Nil(t, err)

	// Buy 2 tokens
	_, err = h(ctx, newValidMsgBuy(2, 4000))
	wars.EndBlocker(ctx, app.WarsKeeper)

	userBalance := app.WarsKeeper.BankKeeper.GetCoins(ctx, userAddress)
	reserveBalance := app.WarsKeeper.GetReserveBalances(ctx, initToken)
	feeBalance := app.WarsKeeper.BankKeeper.GetCoins(ctx, initFeeAddress)
	currentSupply := app.WarsKeeper.MustGetWar(ctx, token).CurrentSupply
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(3767), userBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(2), userBalance.AmountOf(token))
	require.Equal(t, sdk.NewInt(232), reserveBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.OneInt(), feeBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(2), currentSupply.Amount)
}

func TestSellingANonExistingWarFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Sell 10 tokens
	_, err := h(ctx, newValidMsgSell(10))

	require.Error(t, err)
	require.False(t, app.WarsKeeper.WarExists(ctx, token))
}

func TestSellingAWarWhichCannotBeSoldFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war
	createMsg := newValidMsgCreateWar()
	createMsg.AllowSells = false
	h(ctx, createMsg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 10000)})
	require.Nil(t, err)

	// Buy 10 tokens
	h(ctx, newValidMsgBuy(10, 10000))
	wars.EndBlocker(ctx, app.WarsKeeper)

	// Sell 10 tokens
	warPreSell := app.WarsKeeper.MustGetWar(ctx, token)
	_, err = h(ctx, newValidMsgSell(10))
	warPostSell := app.WarsKeeper.MustGetWar(ctx, token)

	require.Error(t, err)
	require.Equal(t, warPostSell.CurrentSupply.Amount, warPreSell.CurrentSupply.Amount)
}

func TestSellWarExceedingOrderQuantityLimitFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war with order quantity limit
	msg := newValidMsgCreateWar()
	msg.OrderQuantityLimits = sdk.NewCoins(sdk.NewInt64Coin(token, 4))
	h(ctx, msg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 10000)})
	require.Nil(t, err)

	// Buy 10 tokens
	h(ctx, newValidMsgBuy(10, 10000))
	wars.EndBlocker(ctx, app.WarsKeeper)

	// Sell 10 tokens
	warPreSell := app.WarsKeeper.MustGetWar(ctx, token)
	_, err = h(ctx, newValidMsgSell(10))
	warPostSell := app.WarsKeeper.MustGetWar(ctx, token)

	require.Error(t, err)
	require.Equal(t, warPostSell.CurrentSupply.Amount, warPreSell.CurrentSupply.Amount)
}

func TestSellingAWarWithAmountGreaterThanBalanceFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war
	h(ctx, newValidMsgCreateWar())

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 10000)})
	require.Nil(t, err)

	// Buy 10 tokens
	h(ctx, newValidMsgBuy(10, 10000))
	wars.EndBlocker(ctx, app.WarsKeeper)

	// Sell 11 tokens
	warPreSell := app.WarsKeeper.MustGetWar(ctx, token)
	_, err = h(ctx, newValidMsgSell(11))
	warPostSell := app.WarsKeeper.MustGetWar(ctx, token)

	userBalance := app.WarsKeeper.BankKeeper.GetCoins(ctx, userAddress)
	require.Error(t, err)
	require.Equal(t, warPostSell.CurrentSupply.Amount, warPreSell.CurrentSupply.Amount)
	require.Equal(t, sdk.NewInt(10), userBalance.AmountOf(token))
}

func TestSellingAWarWhichSellerDoesNotOwnFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create first war
	h(ctx, newValidMsgCreateWar())

	// Create second war (different token)
	war2Msg := newValidMsgCreateWar()
	war2Msg.Token = token2
	h(ctx, war2Msg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 10000)})
	require.Nil(t, err)

	// Buy 10 tokens
	h(ctx, newValidMsgBuy(10, 10000))
	wars.EndBlocker(ctx, app.WarsKeeper)

	// Sell 11 of a different war
	msg := newValidMsgSell(0) // 0 amount replaced below
	msg.Amount = sdk.NewInt64Coin(token2, 11)
	_, err = h(ctx, msg)

	userBalance := app.WarsKeeper.BankKeeper.GetCoins(ctx, userAddress)
	require.Error(t, err)
	require.Equal(t, sdk.NewInt(10), userBalance.AmountOf(token))
	require.Equal(t, sdk.ZeroInt(), userBalance.AmountOf(token2))
}

func TestSellingMoreTokensThanThereIsSupplyFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war
	h(ctx, newValidMsgCreateWar())

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 10000)})
	require.Nil(t, err)

	// Buy 10 tokens
	h(ctx, newValidMsgBuy(10, 10000))
	wars.EndBlocker(ctx, app.WarsKeeper)

	// Sell an amount greater than the max supply
	warPreSell := app.WarsKeeper.MustGetWar(ctx, token)
	_, err = h(ctx, newValidMsgSell(amountGTMaxSupply))
	warPostSell := app.WarsKeeper.MustGetWar(ctx, token)

	userBalance := app.WarsKeeper.BankKeeper.GetCoins(ctx, userAddress)
	require.Error(t, err)
	require.Equal(t, sdk.NewInt(10), userBalance.AmountOf(token))
	require.Equal(t, warPreSell.CurrentSupply.Amount, warPostSell.CurrentSupply.Amount)
}

func TestSellingAWarCorrectlyPasses(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war
	h(ctx, newValidMsgCreateWar())

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 4000)})
	require.Nil(t, err)

	// Buy 2 tokens
	h(ctx, newValidMsgBuy(2, 4000))
	wars.EndBlocker(ctx, app.WarsKeeper)

	// Sell 2 tokens
	msg := newValidMsgSell(2)
	_, err = h(ctx, msg)
	wars.EndBlocker(ctx, app.WarsKeeper)

	userBalance := app.WarsKeeper.BankKeeper.GetCoins(ctx, userAddress)
	reserveBalance := app.WarsKeeper.GetReserveBalances(ctx, initToken)
	feeBalance := app.WarsKeeper.BankKeeper.GetCoins(ctx, initFeeAddress)
	currentSupply := app.WarsKeeper.MustGetWar(ctx, token).CurrentSupply
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(3997), userBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.ZeroInt(), userBalance.AmountOf(token))
	require.Equal(t, sdk.ZeroInt(), reserveBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(3), feeBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.ZeroInt(), currentSupply.Amount)
}

func TestSwapWarDoesNotExistFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Swap tokens
	_, err := h(ctx, newValidMsgSwap(reserveToken, reserveToken2, 1))

	require.Error(t, err)
	require.False(t, app.WarsKeeper.WarExists(ctx, token))
}

func TestSwapOrderInvalidReserveDenomsFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war
	h(ctx, newValidMsgCreateSwapperWar())

	// Add reserve tokens to user
	coins := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 100000),
		sdk.NewInt64Coin(reserveToken2, 100000),
	)
	err := addCoinsToUser(app, ctx, coins)
	require.Nil(t, err)

	// Buy 2 tokens
	buyMsg := newValidMsgBuy(2, 0) // 0 max prices replaced below
	buyMsg.MaxPrices = sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 10000),
		sdk.NewInt64Coin(reserveToken2, 10000),
	)
	h(ctx, buyMsg)
	wars.EndBlocker(ctx, app.WarsKeeper)

	// Perform swap (invalid instead of reserveToken)
	_, err = h(ctx, newValidMsgSwap("invalid", reserveToken2, 10))
	wars.EndBlocker(ctx, app.WarsKeeper)

	userBalance := app.AccountKeeper.GetAccount(ctx, userAddress).GetCoins()
	require.Error(t, err)
	require.Equal(t, sdk.NewInt(2), userBalance.AmountOf(token))

	// Perform swap (invalid instead of reserveToken2)
	_, err = h(ctx, newValidMsgSwap(reserveToken, "invalid", 10))
	wars.EndBlocker(ctx, app.WarsKeeper)

	userBalance = app.AccountKeeper.GetAccount(ctx, userAddress).GetCoins()
	require.Error(t, err)
	require.Equal(t, sdk.NewInt(2), userBalance.AmountOf(token))
}

func TestSwapOrderQuantityLimitExceededFails(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war with order quantity limit
	createMsg := newValidMsgCreateSwapperWar()
	createMsg.OrderQuantityLimits = sdk.NewCoins(sdk.NewInt64Coin(reserveToken, 4))
	h(ctx, createMsg)

	// Add reserve tokens to user
	coins := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 100000),
		sdk.NewInt64Coin(reserveToken2, 100000),
	)
	err := addCoinsToUser(app, ctx, coins)
	require.Nil(t, err)

	// Buy 2 tokens
	buyMsg := newValidMsgBuy(2, 0) // 0 max prices replaced below
	buyMsg.MaxPrices = sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 10000),
		sdk.NewInt64Coin(reserveToken2, 10000),
	)
	h(ctx, buyMsg)
	wars.EndBlocker(ctx, app.WarsKeeper)

	// Perform swap
	msg := types.NewMsgSwap(userAddress, token, sdk.NewInt64Coin(reserveToken, 5), reserveToken2)
	_, err = h(ctx, msg)

	userBalance := app.AccountKeeper.GetAccount(ctx, userAddress).GetCoins()
	require.Error(t, err)
	require.Equal(t, sdk.NewInt(2), userBalance.AmountOf(token))
}

func TestSwapInvalidAmount(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war
	h(ctx, newValidMsgCreateSwapperWar())

	// Add reserve tokens to user (but not enough)
	nineReserveTokens := sdk.NewInt64Coin(reserveToken, 9)
	tenReserveTokens := sdk.NewInt64Coin(reserveToken, 10)
	err := addCoinsToUser(app, ctx, sdk.Coins{nineReserveTokens})
	require.Nil(t, err)

	// Buy 2 tokens
	buyMsg := newValidMsgBuy(2, 0) // 0 max prices replaced below
	buyMsg.MaxPrices = sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 10000),
		sdk.NewInt64Coin(reserveToken2, 10000),
	)
	h(ctx, buyMsg)
	wars.EndBlocker(ctx, app.WarsKeeper)

	// Perform swap
	msg := types.NewMsgSwap(userAddress, token, tenReserveTokens, reserveToken2)
	_, err = h(ctx, msg)

	require.Error(t, err)
}

func TestSwapValidAmount(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war
	h(ctx, newValidMsgCreateSwapperWar())

	// Add reserve tokens to user
	coins := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 100000),
		sdk.NewInt64Coin(reserveToken2, 100000),
	)
	err := addCoinsToUser(app, ctx, coins)
	require.Nil(t, err)

	// Buy 2 tokens
	buyMsg := newValidMsgBuy(2, 0) // 0 max prices replaced below
	buyMsg.MaxPrices = sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 10000),
		sdk.NewInt64Coin(reserveToken2, 10000),
	)
	h(ctx, buyMsg)
	wars.EndBlocker(ctx, app.WarsKeeper)

	// Perform swap
	_, err = h(ctx, newValidMsgSwap(reserveToken, reserveToken2, 10))
	wars.EndBlocker(ctx, app.WarsKeeper)

	userBalance := app.WarsKeeper.BankKeeper.GetCoins(ctx, userAddress)
	reserveBalance := app.WarsKeeper.GetReserveBalances(ctx, initToken)
	feeBalance := app.WarsKeeper.BankKeeper.GetCoins(ctx, initFeeAddress)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(89990), userBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(90008), userBalance.AmountOf(reserveToken2))
	require.Equal(t, sdk.NewInt(2), userBalance.AmountOf(token))
	require.Equal(t, sdk.NewInt(10009), reserveBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(9992), reserveBalance.AmountOf(reserveToken2))
	require.Equal(t, sdk.OneInt(), feeBalance.AmountOf(reserveToken))
}

func TestSwapValidAmountReversed(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war
	h(ctx, newValidMsgCreateSwapperWar())

	// Add reserve tokens to user
	coins := sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 100000),
		sdk.NewInt64Coin(reserveToken2, 100000),
	)
	err := addCoinsToUser(app, ctx, coins)
	require.Nil(t, err)

	// Buy 2 tokens
	buyMsg := newValidMsgBuy(2, 0) // 0 max prices replaced below
	buyMsg.MaxPrices = sdk.NewCoins(
		sdk.NewInt64Coin(reserveToken, 10000),
		sdk.NewInt64Coin(reserveToken2, 10000),
	)
	h(ctx, buyMsg)
	wars.EndBlocker(ctx, app.WarsKeeper)

	// Perform swap
	_, err = h(ctx, newValidMsgSwap(reserveToken2, reserveToken, 10))
	wars.EndBlocker(ctx, app.WarsKeeper)

	userBalance := app.WarsKeeper.BankKeeper.GetCoins(ctx, userAddress)
	reserveBalance := app.WarsKeeper.GetReserveBalances(ctx, initToken)
	feeBalance := app.WarsKeeper.BankKeeper.GetCoins(ctx, initFeeAddress)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(90008), userBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(89990), userBalance.AmountOf(reserveToken2))
	require.Equal(t, sdk.NewInt(2), userBalance.AmountOf(token))
	require.Equal(t, sdk.NewInt(9992), reserveBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(10009), reserveBalance.AmountOf(reserveToken2))
	require.Equal(t, sdk.OneInt(), feeBalance.AmountOf(reserveToken2))
}

func TestMakeOutcomePayment(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war with 100k outcome payment
	warMsg := newValidMsgCreateWar()
	warMsg.OutcomePayment = sdk.NewCoins(sdk.NewInt64Coin(reserveToken, 100000))
	h(ctx, warMsg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 100000)})
	require.Nil(t, err)

	// Make outcome payment
	_, err = h(ctx, newValidMsgMakeOutcomePayment())
	wars.EndBlocker(ctx, app.WarsKeeper)

	// Check that outcome payment is now in the war reserve
	userBalance := app.WarsKeeper.BankKeeper.GetCoins(ctx, userAddress)
	reserveBalance := app.WarsKeeper.GetReserveBalances(ctx, initToken)
	require.NoError(t, err)
	require.Equal(t, sdk.ZeroInt(), userBalance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(100000), reserveBalance.AmountOf(reserveToken))

	// Check that the war is now in SETTLE state
	require.Equal(t, types.SettleState, app.WarsKeeper.MustGetWar(ctx, token).State)
}

func TestWithdrawShare(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war
	h(ctx, newValidMsgCreateWar())

	// Set war current supply to 3 and state to SETTLE
	war := app.WarsKeeper.MustGetWar(ctx, token)
	war.CurrentSupply = sdk.NewCoin(war.Token, sdk.NewInt(3))
	war.State = types.SettleState
	app.WarsKeeper.SetWar(ctx, token, war)

	// Mint 3 war tokens and send [2 to user 1] and [1 to user 2]
	err := app.SupplyKeeper.MintCoins(ctx, types.WarsMintBurnAccount,
		sdk.NewCoins(sdk.NewInt64Coin(token, 3)))
	require.Nil(t, err)
	err = app.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.WarsMintBurnAccount,
		userAddress, sdk.NewCoins(sdk.NewInt64Coin(token, 2)))
	require.Nil(t, err)
	err = app.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.WarsMintBurnAccount,
		anotherAddress, sdk.NewCoins(sdk.NewInt64Coin(token, 1)))
	require.Nil(t, err)

	// Simulate outcome payment by depositing (freshly minted) 100k into reserve
	hundredK := sdk.NewCoins(sdk.NewCoin(reserveToken, sdk.NewInt(100000)))
	err = app.SupplyKeeper.MintCoins(ctx, types.WarsMintBurnAccount, hundredK)
	require.Nil(t, err)
	err = app.WarsKeeper.DepositReserveFromModule(
		ctx, war.Token, types.WarsMintBurnAccount, hundredK)
	require.Nil(t, err)

	// User 1 withdraws share
	_, err = h(ctx, newValidMsgWithdrawShareFrom(userAddress))
	require.NoError(t, err)
	wars.EndBlocker(ctx, app.WarsKeeper)

	// User 1 had 2 tokens out of the supply of 3 tokens, so user 1 gets 2/3
	user1Balance := app.WarsKeeper.BankKeeper.GetCoins(ctx, userAddress)
	reserveBalance := app.WarsKeeper.GetReserveBalances(ctx, initToken)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(66666), user1Balance.AmountOf(reserveToken))
	require.Equal(t, sdk.NewInt(33334), reserveBalance.AmountOf(reserveToken))

	// Note: rounding is rounded to floor, so despite user 1 being owed 66666.67
	// tokens, user 1 gets 66666 and not 66667 tokens. Then, since user 2 now owns
	// the entire share of the war tokens, they will get 100% of the remaining
	// 33334 tokens, which is more than what was initially owed (33333.33).

	// User 2 withdraws share
	_, err = h(ctx, newValidMsgWithdrawShareFrom(anotherAddress))
	require.NoError(t, err)
	wars.EndBlocker(ctx, app.WarsKeeper)

	// User 2 had 1 token out of the remaining supply of 1 token, so user 2 gets all remaining
	user2Balance := app.WarsKeeper.BankKeeper.GetCoins(ctx, anotherAddress)
	reserveBalance = app.WarsKeeper.GetReserveBalances(ctx, initToken)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(33334), user2Balance.AmountOf(reserveToken))
	require.Equal(t, sdk.ZeroInt(), reserveBalance.AmountOf(reserveToken))
}

func TestDecrementRemainingBlocksCountAfterEndBlock(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	two := sdk.NewUint(2)
	one := sdk.NewUint(1)

	// Create war
	createMsg := newValidMsgCreateWar()
	createMsg.BatchBlocks = two
	h(ctx, createMsg)

	require.Equal(t, two, app.WarsKeeper.MustGetBatch(ctx, token).BlocksRemaining)
	wars.EndBlocker(ctx, app.WarsKeeper)
	require.Equal(t, one, app.WarsKeeper.MustGetBatch(ctx, token).BlocksRemaining)
}

func TestEndBlockerDoesNotPerformOrdersBeforeASpecifiedNumberOfBlocks(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war with batch blocks set to 2
	createMsg := newValidMsgCreateWar()
	createMsg.BatchBlocks = sdk.NewUint(2)
	h(ctx, createMsg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 1000000)})
	require.Nil(t, err)

	// Buy 4 tokens
	h(ctx, newValidMsgBuy(2, 10000))
	h(ctx, newValidMsgBuy(2, 10000))
	wars.EndBlocker(ctx, app.WarsKeeper)

	require.Equal(t, len(app.WarsKeeper.MustGetBatch(ctx, token).Buys), 2)
}

func TestEndBlockerPerformsOrdersAfterASpecifiedNumberOfBlocks(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war
	createMsg := newValidMsgCreateWar()
	createMsg.BatchBlocks = sdk.NewUint(3)
	h(ctx, createMsg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 1000000)})
	require.Nil(t, err)

	// Buy 4 tokens
	h(ctx, newValidMsgBuy(2, 10000))
	h(ctx, newValidMsgBuy(2, 10000))

	// Run EndBlocker for N times, where N = BatchBlocks
	batchBlocksInt := int(createMsg.BatchBlocks.Uint64())
	for i := 0; i <= batchBlocksInt; i++ {
		wars.EndBlocker(ctx, app.WarsKeeper)
	}

	// Buys have been performed
	require.Equal(t, 0, len(app.WarsKeeper.MustGetBatch(ctx, token).Buys))
}

func TestEndBlockerAugmentedFunction(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war with augmented function type
	createMsg := newValidMsgCreateAugmentedWar()
	h(ctx, createMsg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 1000000)})
	require.Nil(t, err)

	// Get war to confirm allowSells==false, S0==50000, state==hatch
	war := app.WarsKeeper.MustGetWar(ctx, token)
	require.False(t, war.AllowSells)
	require.Equal(t, sdk.NewDec(50000), war.FunctionParameters.AsMap()["S0"])
	require.Equal(t, types.HatchState, war.State)

	// - Buy 49999 tokens; just below S0
	// - Cannot buy 2 tokens in the meantime, since this exceeds S0
	// - Cannot sell tokens (not even 1) in hatch state
	_, err = h(ctx, newValidMsgBuy(49999, 100000))
	require.NoError(t, err)
	_, err = h(ctx, newValidMsgBuy(2, 100000))
	require.Error(t, err)
	_, err = h(ctx, newValidMsgSell(1))
	require.Error(t, err)
	wars.EndBlocker(ctx, app.WarsKeeper)

	// Confirm allowSells and state still the same
	war = app.WarsKeeper.MustGetWar(ctx, token)
	require.False(t, war.AllowSells)
	require.Equal(t, types.HatchState, war.State)

	// Buy 1 more token, to reach S0 => state is now open
	h(ctx, newValidMsgBuy(1, 100000))
	wars.EndBlocker(ctx, app.WarsKeeper)

	// Confirm allowSells==true, state==open
	war = app.WarsKeeper.MustGetWar(ctx, token)
	require.True(t, war.AllowSells)
	require.Equal(t, types.OpenState, war.State)

	// Check user balance of tokens
	balance := app.BankKeeper.GetCoins(ctx, userAddress).AmountOf(token).Int64()
	require.Equal(t, int64(50000), balance)

	// Can now sell tokens (all 50,000 of them)
	_, err = h(ctx, newValidMsgSell(50000))
	require.NoError(t, err)
	wars.EndBlocker(ctx, app.WarsKeeper)
	balance = app.BankKeeper.GetCoins(ctx, userAddress).AmountOf(token).Int64()
	require.Equal(t, int64(0), balance)
}

func TestEndBlockerAugmentedFunctionDecimalS0(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war with augmented function type
	createMsg := newValidMsgCreateAugmentedWar()
	h(ctx, createMsg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 1000000)})
	require.Nil(t, err)

	// Change war's S0 parameter to 49999.5
	decimalS0 := sdk.MustNewDecFromStr("49999.5")
	war := app.WarsKeeper.MustGetWar(ctx, token)
	for i, p := range war.FunctionParameters {
		if p.Param == "S0" {
			war.FunctionParameters[i].Value = decimalS0
			break
		}
	}
	app.WarsKeeper.SetWar(ctx, war.Token, war)

	// Get war to confirm S0==49999.5, allowSells==false, state==hatch
	war = app.WarsKeeper.MustGetWar(ctx, token)
	require.Equal(t, decimalS0, war.FunctionParameters.AsMap()["S0"])
	require.False(t, war.AllowSells)
	require.Equal(t, types.HatchState, war.State)

	// Buy 49999 tokens; just below S0
	h(ctx, newValidMsgBuy(49999, 100000))
	wars.EndBlocker(ctx, app.WarsKeeper)

	// Confirm allowSells and state still the same
	war = app.WarsKeeper.MustGetWar(ctx, token)
	require.False(t, war.AllowSells)
	require.Equal(t, types.HatchState, war.State)

	// Buy 1 more token, to exceed S0
	h(ctx, newValidMsgBuy(1, 100000))
	wars.EndBlocker(ctx, app.WarsKeeper)

	// Confirm allowSells==true, state==open
	war = app.WarsKeeper.MustGetWar(ctx, token)
	require.True(t, war.AllowSells)
	require.Equal(t, types.OpenState, war.State)
}

func TestEndBlockerAugmentedFunctionSmallBuys(t *testing.T) {
	app, ctx := createTestApp(false)
	h := wars.NewHandler(app.WarsKeeper)

	// Create war with augmented function type, small params, and zero fees
	createMsg := newValidMsgCreateAugmentedWar()
	createMsg.FunctionParameters = types.FunctionParams{
		types.NewFunctionParam("d0", sdk.MustNewDecFromStr("10.0")),
		types.NewFunctionParam("p0", sdk.MustNewDecFromStr("1.0")),
		types.NewFunctionParam("theta", sdk.MustNewDecFromStr("0.9")),
		types.NewFunctionParam("kappa", sdk.MustNewDecFromStr("3.0"))}
	createMsg.TxFeePercentage = sdk.ZeroDec()
	createMsg.ExitFeePercentage = sdk.ZeroDec()
	h(ctx, createMsg)

	// Add reserve tokens to user
	err := addCoinsToUser(app, ctx, sdk.Coins{sdk.NewInt64Coin(reserveToken, 1000000)})
	require.Nil(t, err)

	// Get war to confirm allowSells==false, S0==10, R0==1 state==hatch
	war := app.WarsKeeper.MustGetWar(ctx, token)
	require.False(t, war.AllowSells)
	require.Equal(t, sdk.NewDec(10), war.FunctionParameters.AsMap()["S0"])
	require.Equal(t, sdk.NewDec(1), war.FunctionParameters.AsMap()["R0"])
	require.Equal(t, types.HatchState, war.State)

	// Perform 10 buys of 1 token each
	for i := 0; i < 10; i++ {
		_, err := h(ctx, newValidMsgBuy(1, 1))
		require.NoError(t, err)
	}
	wars.EndBlocker(ctx, app.WarsKeeper)

	// Confirm allowSells==true, state==open
	war = app.WarsKeeper.MustGetWar(ctx, token)
	require.True(t, war.AllowSells)
	require.Equal(t, types.OpenState, war.State)

	// Confirm reserve balance is R0 [i.e. d0*(1-theta)] = 1
	require.Equal(t, int64(1), war.CurrentReserve[0].Amount.Int64())

	// Confirm fee address balance is d0*theta = 9
	feeAddressBalance := app.BankKeeper.GetCoins(
		ctx, war.FeeAddress).AmountOf(reserveToken).Int64()
	require.Equal(t, int64(9), feeAddressBalance)
}
