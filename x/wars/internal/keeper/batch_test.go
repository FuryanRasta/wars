package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/warmage-sports/wars/x/wars/internal/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBatchExistsSetGet(t *testing.T) {
	app, ctx := createTestApp(false)

	// Batch doesn't exist yet
	require.False(t, app.WarsKeeper.BatchExists(ctx, token))

	// Add batch
	batchAdded := getValidBatch()
	app.WarsKeeper.SetBatch(ctx, token, batchAdded)

	// Batch now exists (but this has nothing to do with last batch)
	require.True(t, app.WarsKeeper.BatchExists(ctx, token))
	require.False(t, app.WarsKeeper.LastBatchExists(ctx, token))

	// Must get batch
	batchFetched := app.WarsKeeper.MustGetBatch(ctx, token)

	// Batch fetched is equal to added batch
	require.Equal(t, batchAdded, batchFetched)
}

func TestLastBatchExistsSetGet(t *testing.T) {
	app, ctx := createTestApp(false)

	// Batch doesn't exist yet
	require.False(t, app.WarsKeeper.BatchExists(ctx, token))

	// Add last batch
	batchAdded := getValidBatch()
	app.WarsKeeper.SetLastBatch(ctx, token, batchAdded)

	// Last batch now exists (but this has nothing to do with current batch)
	require.True(t, app.WarsKeeper.LastBatchExists(ctx, token))
	require.False(t, app.WarsKeeper.BatchExists(ctx, token))

	// Must get last batch
	batchFetched := app.WarsKeeper.MustGetLastBatch(ctx, token)

	// Last batch fetched is equal to added batch
	require.Equal(t, batchAdded, batchFetched)
}

func TestBatchAddBuyOrder(t *testing.T) {
	app, ctx := createTestApp(false)

	// Add batch
	batchAdded := getValidBatch()
	app.WarsKeeper.SetBatch(ctx, token, batchAdded)
	require.True(t, app.WarsKeeper.BatchExists(ctx, token))

	// Add buy order
	bo := getValidBuyOrder()
	app.WarsKeeper.AddBuyOrder(ctx, token, bo, buyPrices, sellPrices)

	// Get and check batch
	batchFetched := app.WarsKeeper.MustGetBatch(ctx, token)
	require.Equal(t, len(batchFetched.Buys), 1)
	require.Equal(t, len(batchFetched.Sells), 0)
	require.Equal(t, len(batchFetched.Swaps), 0)
	require.Equal(t, batchFetched.TotalBuyAmount, bo.Amount)
	require.Equal(t, batchFetched.TotalSellAmount, sdk.NewCoin(token, sdk.ZeroInt()))
	require.Equal(t, batchFetched.Buys[0], bo)
}

func TestBatchAddSellOrder(t *testing.T) {
	app, ctx := createTestApp(false)

	// Add batch
	batchAdded := getValidBatch()
	app.WarsKeeper.SetBatch(ctx, token, batchAdded)
	require.True(t, app.WarsKeeper.BatchExists(ctx, token))

	// Add sell order
	so := getValidSellOrder()
	app.WarsKeeper.AddSellOrder(ctx, token, so, buyPrices, sellPrices)

	// Get and check batch
	batchFetched := app.WarsKeeper.MustGetBatch(ctx, token)
	require.Equal(t, len(batchFetched.Buys), 0)
	require.Equal(t, len(batchFetched.Sells), 1)
	require.Equal(t, len(batchFetched.Swaps), 0)
	require.Equal(t, batchFetched.TotalBuyAmount, sdk.NewCoin(token, sdk.ZeroInt()))
	require.Equal(t, batchFetched.TotalSellAmount, so.Amount)
	require.Equal(t, batchFetched.Sells[0], so)
}

func TestBatchAddSwapOrder(t *testing.T) {
	app, ctx := createTestApp(false)

	// Add batch
	batchAdded := getValidBatch()
	app.WarsKeeper.SetBatch(ctx, token, batchAdded)
	require.True(t, app.WarsKeeper.BatchExists(ctx, token))

	// Add swap order
	swapOrder := getValidSwapOrder()
	app.WarsKeeper.AddSwapOrder(ctx, token, swapOrder)

	// Get and check batch
	batchFetched := app.WarsKeeper.MustGetBatch(ctx, token)
	require.Equal(t, len(batchFetched.Buys), 0)
	require.Equal(t, len(batchFetched.Sells), 0)
	require.Equal(t, len(batchFetched.Swaps), 1)
	require.Equal(t, batchFetched.TotalBuyAmount, sdk.NewCoin(token, sdk.ZeroInt()))
	require.Equal(t, batchFetched.TotalSellAmount, sdk.NewCoin(token, sdk.ZeroInt()))
	require.Equal(t, batchFetched.Swaps[0], swapOrder)
}

func TestGetBatchBuySellPrices(t *testing.T) {
	app, ctx := createTestApp(false)

	// Create war and get current batch-independent prices at supply=10
	war := getValidWar()
	war.CurrentSupply = sdk.NewInt64Coin(war.Token, 10)
	app.WarsKeeper.SetWar(ctx, token, war)
	currentPrices, _ := war.GetCurrentPricesPT(nil)
	expectedCurrentPrices, _ := war.GetPricesAtSupply(war.CurrentSupply.Amount)
	require.Equal(t, expectedCurrentPrices, currentPrices)

	// Set fixed buy/sell amount
	fiveTokens := sdk.NewInt64Coin(war.Token, 5)
	fiveDec := sdk.NewDec(5)

	// Add appropriate amount of reserve tokens (freshly minted) to reserve
	expectedReserve := war.ReserveAtSupply(war.CurrentSupply.Amount)
	expectedRounded := expectedReserve.Ceil().TruncateInt()
	reserveBalance := sdk.NewCoins(sdk.NewCoin(war.ReserveTokens[0], expectedRounded))
	err := app.SupplyKeeper.MintCoins(ctx, types.WarsMintBurnAccount, reserveBalance)
	require.Nil(t, err)
	err = app.WarsKeeper.DepositReserveFromModule(
		ctx, war.Token, types.WarsMintBurnAccount, reserveBalance)
	require.Nil(t, err)

	// Create empty batch
	batch := getValidBatch()

	// Initially equal to current prices due to no orders
	buyPrices, sellPrices, err := app.WarsKeeper.GetBatchBuySellPrices(ctx, war.Token, batch)
	require.Equal(t, currentPrices, buyPrices)
	require.Equal(t, currentPrices, sellPrices)
	require.Nil(t, err)

	// ------------------------------------------------

	// (Re)Create batch with buy order (nil max prices)
	batch = getValidBatch()
	bo := types.NewBuyOrder(buyerAddress, fiveTokens, reserveBalance)
	batch.Buys = append(batch.Buys, bo)
	batch.TotalBuyAmount = batch.TotalBuyAmount.Add(bo.Amount)

	// Calculate expected buy price
	expectedPrices1, err := war.GetPricesToMint(bo.Amount.Amount, reserveBalance)
	require.NotNil(t, expectedPrices1)
	require.Nil(t, err)
	expectedBuyPricesPerToken := types.DivideDecCoinsByDec(expectedPrices1, fiveDec)

	// Since sells=0, buy prices are based on just the buy and sell prices are current prices
	buyPrices, sellPrices, err = app.WarsKeeper.GetBatchBuySellPrices(ctx, war.Token, batch)
	require.Equal(t, expectedBuyPricesPerToken, buyPrices)
	require.Equal(t, currentPrices, sellPrices)
	require.Nil(t, err)

	// ------------------------------------------------

	// (Re)Create batch with sell order
	batch = getValidBatch()
	so := types.NewSellOrder(sellerAddress, fiveTokens)
	batch.Sells = append(batch.Sells, so)
	batch.TotalSellAmount = batch.TotalSellAmount.Add(so.Amount)

	// Calculate expected sell price
	expectedReturns := war.GetReturnsForBurn(so.Amount.Amount, reserveBalance)
	require.NotNil(t, expectedReturns)
	expectedSellPricesPerToken := types.DivideDecCoinsByDec(expectedReturns, fiveDec)

	// Since sells=0, buy prices are based on just the buy and sell prices are current prices
	buyPrices, sellPrices, err = app.WarsKeeper.GetBatchBuySellPrices(ctx, war.Token, batch)
	require.Equal(t, currentPrices, buyPrices)
	require.Equal(t, expectedSellPricesPerToken, sellPrices)
	require.Nil(t, err)

	// ------------------------------------------------

	// (Re)Create batch with buy amount > sell amount
	batch = getValidBatch()
	bo1 := types.NewBuyOrder(buyerAddress, fiveTokens, nil)
	bo2 := types.NewBuyOrder(buyerAddress, fiveTokens, nil) // 5 more
	so = types.NewSellOrder(sellerAddress, fiveTokens)
	batch.Buys = append(batch.Buys, bo1, bo2)
	batch.Sells = append(batch.Sells, so)
	batch.TotalBuyAmount = batch.TotalBuyAmount.Add(bo1.Amount).Add(bo2.Amount)
	batch.TotalSellAmount = batch.TotalSellAmount.Add(so.Amount)

	// Calculate expected buy price (for 5 [mint-price] + 5 [current-price] tokens)
	expectedPrices1, err = war.GetPricesToMint(fiveTokens.Amount, reserveBalance)
	require.Nil(t, err)
	require.NotNil(t, expectedPrices1)
	expectedPrices2 := currentPrices.MulDec(fiveDec)
	totalExpectedPrices := expectedPrices1.Add(expectedPrices2...)
	expectedBuyPricesPerToken = types.DivideDecCoinsByDec(totalExpectedPrices, fiveDec.Add(fiveDec))

	// Since buys>sells, buy prices are affected by extra buys and sell prices are current prices
	buyPrices, sellPrices, err = app.WarsKeeper.GetBatchBuySellPrices(ctx, war.Token, batch)
	require.Equal(t, expectedBuyPricesPerToken, buyPrices)
	require.Equal(t, currentPrices, sellPrices)
	require.Nil(t, err)

	// ------------------------------------------------

	// (Re)Create batch with sell amount > buy amount
	batch = getValidBatch()
	bo = types.NewBuyOrder(buyerAddress, fiveTokens, nil)
	so1 := types.NewSellOrder(sellerAddress, fiveTokens)
	so2 := types.NewSellOrder(sellerAddress, fiveTokens)
	batch.Buys = append(batch.Buys, bo)
	batch.Sells = append(batch.Sells, so1, so2)
	batch.TotalBuyAmount = batch.TotalBuyAmount.Add(bo1.Amount)
	batch.TotalSellAmount = batch.TotalSellAmount.Add(so1.Amount).Add(so2.Amount)

	// Calculate expected sell price (for 5 [burn-price] + 5 [current-price] tokens)
	expectedReturns1 := war.GetReturnsForBurn(fiveTokens.Amount, reserveBalance)
	require.Nil(t, err)
	require.NotNil(t, expectedReturns1)
	expectedReturns2 := currentPrices.MulDec(fiveDec)
	totalExpectedReturns := expectedReturns1.Add(expectedReturns2...)
	expectedSellPricesPerToken = types.DivideDecCoinsByDec(totalExpectedReturns, fiveDec.Add(fiveDec))

	// Since sells>buys, sell prices are affected by extra sells and buy prices are current prices
	buyPrices, sellPrices, err = app.WarsKeeper.GetBatchBuySellPrices(ctx, war.Token, batch)
	require.Equal(t, currentPrices, buyPrices)
	require.Equal(t, expectedSellPricesPerToken, sellPrices)
	require.Nil(t, err)
}

func TestGetUpdatedBatchPricesAfterBuy(t *testing.T) {
	app, ctx := createTestApp(false)

	// Create war and batch
	war := getValidWar()
	batch := getValidBatch()
	app.WarsKeeper.SetWar(ctx, war.Token, war)
	app.WarsKeeper.SetBatch(ctx, war.Token, batch)

	// Fixed buy amount
	buyAmount := sdk.NewCoin(war.Token, sdk.OneInt())

	// Buy order with buy amount greater than max supply not fulfillable
	// MaxPrices is set to nil since it is not relevant in this scenario
	maxSupplyPlus1 := war.MaxSupply.Add(sdk.NewCoin(war.Token, sdk.OneInt()))
	bo := types.NewBuyOrder(buyerAddress, maxSupplyPlus1, nil)
	_, _, err := app.WarsKeeper.GetUpdatedBatchPricesAfterBuy(ctx, war.Token, bo)
	require.Error(t, err)

	// Buy order with max prices lower than prices not fulfillable
	maxPrices := sdk.NewCoins(sdk.NewCoin(war.ReserveTokens[0], sdk.OneInt()))
	bo = types.NewBuyOrder(buyerAddress, buyAmount, maxPrices)
	_, _, err = app.WarsKeeper.GetUpdatedBatchPricesAfterBuy(ctx, war.Token, bo)
	require.Error(t, err)

	// Check buy prices for fulfillable buy order
	maxPrices = sdk.NewCoins(sdk.NewInt64Coin(war.ReserveTokens[0], 10000000))
	bo = types.NewBuyOrder(buyerAddress, buyAmount, maxPrices)
	buyPrices, sellPrices, err = app.WarsKeeper.GetUpdatedBatchPricesAfterBuy(ctx, war.Token, bo)
	expectedBuyPrices, _ := war.GetPricesToMint(buyAmount.Amount, nil)
	expectedSellPrices, _ := war.GetCurrentPricesPT(nil)
	require.Nil(t, err)
	require.Equal(t, expectedBuyPrices, buyPrices)
	require.Equal(t, expectedSellPrices, sellPrices)
}

func TestGetUpdatedBatchPricesAfterSell(t *testing.T) {
	app, ctx := createTestApp(false)

	// Create war and batch
	war := getValidWar()
	batch := getValidBatch()
	app.WarsKeeper.SetWar(ctx, war.Token, war)
	app.WarsKeeper.SetBatch(ctx, war.Token, batch)

	// Fixed sell amount
	sellAmount := sdk.NewCoin(war.Token, sdk.OneInt())

	// Sell order when current supply is zero is not fulfillable
	so := types.NewSellOrder(sellerAddress, sellAmount)
	_, _, err := app.WarsKeeper.GetUpdatedBatchPricesAfterSell(ctx, war.Token, so)
	require.Error(t, err)

	// Increase current supply to amount to be sold and set an appropriate reserve
	war.CurrentSupply = sellAmount
	app.WarsKeeper.SetWar(ctx, war.Token, war)
	reserveBalance := sdk.NewCoins(sdk.NewInt64Coin(war.ReserveTokens[0], 10000000))
	err = app.SupplyKeeper.MintCoins(ctx, types.WarsMintBurnAccount, reserveBalance)
	require.Nil(t, err)
	err = app.WarsKeeper.DepositReserveFromModule(
		ctx, war.Token, types.WarsMintBurnAccount, reserveBalance)

	// Check sell prices for fulfillable sell order
	so = types.NewSellOrder(sellerAddress, sellAmount)
	buyPrices, sellPrices, err = app.WarsKeeper.GetUpdatedBatchPricesAfterSell(ctx, war.Token, so)
	expectedBuyPrices, _ := war.GetCurrentPricesPT(nil)
	expectedSellPrices := war.GetReturnsForBurn(sellAmount.Amount, reserveBalance)
	require.Nil(t, err)
	require.Equal(t, expectedBuyPrices, buyPrices)
	require.Equal(t, expectedSellPrices, sellPrices)
}

func TestPerformBuyAtPrice(t *testing.T) {
	app, ctx := createTestApp(false)
	war := getValidWar()

	buyPrices := sdk.DecCoins{sdk.NewInt64DecCoin(reserveToken, 100)}
	maxPrices := sdk.Coins{sdk.NewInt64Coin(reserveToken, 1100)}

	testCases := []struct {
		amount         sdk.Int
		maxPrices      sdk.Coins
		txFee          sdk.Dec
		expectedPrices sdk.Int
		fulfillable    bool
	}{
		{
			sdk.NewInt(10), maxPrices, sdk.ZeroDec(), sdk.NewInt(1000), true,
		}, // (10 * 100) + (10 * FEE) = 1000 <= 1100, where FEE=0
		{
			sdk.NewInt(11), maxPrices, sdk.ZeroDec(), sdk.NewInt(1100), true,
		}, // (11 * 100) + (11 * FEE) = 1100 <= 1100, where FEE=0
		{
			sdk.NewInt(12), maxPrices, sdk.ZeroDec(), sdk.NewInt(1200), false,
		}, // (12 * 100) + (12 * FEE) = 1200 > 1100, where FEE=0 [not fulfillable]
		{
			sdk.NewInt(10), maxPrices, sdk.NewDec(10), sdk.NewInt(1100), true,
		}, // (10 * 100) + (10 * FEE) = 1100 <= 1100, where FEE=10
		{
			sdk.NewInt(10), maxPrices, sdk.NewDec(20), sdk.NewInt(1200), false,
		}, // (10 * 100) + (10 * FEE) = 1200 > 1100, where FEE=20 [not fulfillable]
	}

	for _, tc := range testCases {
		// Create buy order
		amount := sdk.NewCoin(war.Token, tc.amount)
		bo := types.NewBuyOrder(buyerAddress, amount, tc.maxPrices)

		// Set transaction fee
		war.TxFeePercentage = tc.txFee
		app.WarsKeeper.SetWar(ctx, war.Token, war)

		// Calculate total prices
		reservePrice := buyPrices[0].Amount.MulInt(bo.Amount.Amount)
		reservePrices := sdk.DecCoins{sdk.NewDecCoinFromDec(buyPrices[0].Denom, reservePrice)}
		reservePricesRounded := types.RoundReserveReturns(reservePrices)
		txFees := war.GetTxFees(reservePrices)
		totalPrices := reservePricesRounded.Add(txFees...)

		// Check expected prices
		require.Equal(t, totalPrices.AmountOf(reserveToken), tc.expectedPrices)

		// Add reserve tokens paid by buyer to module account address
		moduleAcc := app.SupplyKeeper.GetModuleAccount(ctx, types.BatchesIntermediaryAccount)
		err := app.BankKeeper.SetCoins(ctx, moduleAcc.GetAddress(), tc.maxPrices)
		require.NoError(t, err)

		// Previous values
		prevSupplySDK := app.SupplyKeeper.GetSupply(ctx).GetTotal().AmountOf(war.Token)
		prevSupplyWars := app.WarsKeeper.MustGetWar(ctx, war.Token).CurrentSupply
		prevModuleAccBal := app.BankKeeper.GetCoins(ctx, moduleAcc.GetAddress())
		prevFeeAddrBal := app.BankKeeper.GetCoins(ctx, war.FeeAddress)
		prevBuyerBal := app.BankKeeper.GetCoins(ctx, buyerAddress)
		prevReserveBal := app.WarsKeeper.GetReserveBalances(ctx, war.Token)

		// Perform buy
		err = app.WarsKeeper.PerformBuyAtPrice(ctx, war.Token, bo, buyPrices)

		// Check if buy is fulfillable (i.e. if maxPrices >= totalPrices)
		if tc.fulfillable {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
			continue // app would panic at this stage
		}

		// Calculate increase in buyer balance
		remainderForBuyer := tc.maxPrices.Sub(totalPrices)
		tokensBought := sdk.NewCoin(war.Token, tc.amount)
		increaseInBuyerBal := sdk.Coins{tokensBought}.Add(remainderForBuyer...)

		// New values
		newSupplySDK := app.SupplyKeeper.GetSupply(ctx).GetTotal().AmountOf(war.Token)
		newSupplyWars := app.WarsKeeper.MustGetWar(ctx, war.Token).CurrentSupply
		newModuleAccBal := app.BankKeeper.GetCoins(ctx, moduleAcc.GetAddress())
		newFeeAddrBal := app.BankKeeper.GetCoins(ctx, war.FeeAddress)
		newBuyerBal := app.BankKeeper.GetCoins(ctx, buyerAddress)
		newReserveBal := app.WarsKeeper.GetReserveBalances(ctx, war.Token)

		require.Equal(t, prevSupplySDK.Add(tc.amount), newSupplySDK)
		require.Equal(t, prevSupplyWars.Add(tokensBought), newSupplyWars)
		require.Equal(t, prevModuleAccBal.Sub(tc.maxPrices), newModuleAccBal)
		require.Equal(t, txFees.IsZero(), tc.txFee.IsZero())
		require.Equal(t, prevFeeAddrBal.Add(txFees...), newFeeAddrBal.Add(nil...))
		require.Equal(t, prevBuyerBal.Add(increaseInBuyerBal...), newBuyerBal)
		require.Equal(t, prevReserveBal.Add(reservePricesRounded...), newReserveBal)
	}
}

func TestPerformBuyAtPriceAugmentedFunction(t *testing.T) {
	app, ctx := createTestApp(false)
	war := getValidAugmentedFunctionWar()
	war.FunctionParameters = types.FunctionParams{
		types.NewFunctionParam("d0", sdk.MustNewDecFromStr("500.0")),
		types.NewFunctionParam("p0", sdk.MustNewDecFromStr("100.0")),
		types.NewFunctionParam("theta", sdk.MustNewDecFromStr("0.4")),
		types.NewFunctionParam("kappa", sdk.MustNewDecFromStr("3.0"))}
	args := war.FunctionParameters.AsMap()
	// price p0 set to 100 so that this test matches other TestPerformBuyAtPrice

	buyPrices := sdk.DecCoins{sdk.NewDecCoinFromDec(reserveToken, args["p0"])}
	maxPrices := sdk.Coins{sdk.NewInt64Coin(reserveToken, 1100)}

	testCases := []struct {
		amount         sdk.Int
		maxPrices      sdk.Coins
		txFee          sdk.Dec
		state          string
		expectedPrices sdk.Int
		fulfillable    bool
	}{
		{
			sdk.NewInt(10), maxPrices, sdk.ZeroDec(), types.HatchState, sdk.NewInt(1000), true,
		}, // (10 * 100) + (10 * FEE) = 1000 <= 1100, where FEE=0
		{
			sdk.NewInt(10), maxPrices, sdk.NewDec(10), types.HatchState, sdk.NewInt(1100), true,
		}, // (10 * 100) + (10 * FEE) = 1100 <= 1100, where FEE=10
		{
			sdk.NewInt(10), maxPrices, sdk.NewDec(20), types.HatchState, sdk.NewInt(1200), false,
		}, // (10 * 100) + (10 * FEE) = 1200 > 1100, where FEE=20 [not fulfillable]
		{
			sdk.NewInt(10), maxPrices, sdk.ZeroDec(), types.OpenState, sdk.NewInt(1000), true,
		}, // (10 * 100) + (10 * FEE) = 1000 <= 1100, where FEE=0
		{
			sdk.NewInt(10), maxPrices, sdk.NewDec(10), types.OpenState, sdk.NewInt(1100), true,
		}, // (10 * 100) + (10 * FEE) = 1100 <= 1100, where FEE=10
		{
			sdk.NewInt(10), maxPrices, sdk.NewDec(20), types.OpenState, sdk.NewInt(1200), false,
		}, // (10 * 100) + (10 * FEE) = 1200 > 1100, where FEE=20 [not fulfillable]
	}

	for _, tc := range testCases {
		// Create buy order
		amount := sdk.NewCoin(war.Token, tc.amount)
		bo := types.NewBuyOrder(buyerAddress, amount, tc.maxPrices)

		// Set transaction fee and state
		war.TxFeePercentage = tc.txFee
		war.State = tc.state
		app.WarsKeeper.SetWar(ctx, war.Token, war)

		// Calculate total prices
		reservePrice := buyPrices[0].Amount.MulInt(bo.Amount.Amount)
		reservePrices := sdk.DecCoins{sdk.NewDecCoinFromDec(buyPrices[0].Denom, reservePrice)}
		reservePricesRounded := types.RoundReserveReturns(reservePrices)
		txFees := war.GetTxFees(reservePrices)
		totalPrices := reservePricesRounded.Add(txFees...)

		// Check expected prices
		require.Equal(t, totalPrices.AmountOf(reserveToken), tc.expectedPrices)

		// Add reserve tokens paid by buyer to module account address
		moduleAcc := app.SupplyKeeper.GetModuleAccount(ctx, types.BatchesIntermediaryAccount)
		err := app.BankKeeper.SetCoins(ctx, moduleAcc.GetAddress(), tc.maxPrices)
		require.NoError(t, err)

		// Previous values
		prevSupplySDK := app.SupplyKeeper.GetSupply(ctx).GetTotal().AmountOf(war.Token)
		prevSupplyWars := app.WarsKeeper.MustGetWar(ctx, war.Token).CurrentSupply
		prevModuleAccBal := app.BankKeeper.GetCoins(ctx, moduleAcc.GetAddress())
		prevFeeAddrBal := app.BankKeeper.GetCoins(ctx, war.FeeAddress)
		prevBuyerBal := app.BankKeeper.GetCoins(ctx, buyerAddress)
		prevReserveBal := app.WarsKeeper.GetReserveBalances(ctx, war.Token)

		// Perform buy
		err = app.WarsKeeper.PerformBuyAtPrice(ctx, war.Token, bo, buyPrices)

		// Check if buy is fulfillable (i.e. if maxPrices >= totalPrices)
		if tc.fulfillable {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
			continue // app would panic at this stage
		}

		// Calculate increase in buyer balance
		remainderForBuyer := tc.maxPrices.Sub(totalPrices)
		tokensBought := sdk.NewCoin(war.Token, tc.amount)
		increaseInBuyerBal := sdk.Coins{tokensBought}.Add(remainderForBuyer...)

		// New values
		newSupplySDK := app.SupplyKeeper.GetSupply(ctx).GetTotal().AmountOf(war.Token)
		newSupplyWars := app.WarsKeeper.MustGetWar(ctx, war.Token).CurrentSupply
		newModuleAccBal := app.BankKeeper.GetCoins(ctx, moduleAcc.GetAddress())
		newFeeAddrBal := app.BankKeeper.GetCoins(ctx, war.FeeAddress)
		newBuyerBal := app.BankKeeper.GetCoins(ctx, buyerAddress)
		newReserveBal := app.WarsKeeper.GetReserveBalances(ctx, war.Token)

		require.Equal(t, prevSupplySDK.Add(tc.amount), newSupplySDK)
		require.Equal(t, prevSupplyWars.Add(tokensBought), newSupplyWars)
		require.Equal(t, prevModuleAccBal.Sub(tc.maxPrices), newModuleAccBal)
		require.Equal(t, txFees.IsZero(), tc.txFee.IsZero())
		require.Equal(t, prevBuyerBal.Add(increaseInBuyerBal...), newBuyerBal)
		if tc.state == types.HatchState {
			toInitialReserve, _ := sdk.NewDecCoinsFromCoins(reservePricesRounded...).MulDec(
				sdk.OneDec().Sub(args["theta"])).TruncateDecimal()
			toFundingPool := txFees.Add(reservePricesRounded.Sub(toInitialReserve)...)
			require.Equal(t, prevReserveBal.Add(toInitialReserve...), newReserveBal)
			require.Equal(t, prevFeeAddrBal.Add(toFundingPool...), newFeeAddrBal)
		} else {
			require.Equal(t, prevFeeAddrBal.Add(txFees...), newFeeAddrBal)
			require.Equal(t, prevReserveBal.Add(reservePricesRounded...), newReserveBal)
		}
	}
}

func TestPerformSellAtPrice(t *testing.T) {
	app, ctx := createTestApp(false)
	war := getValidWar()

	sellAmount := sdk.NewInt64Coin(war.Token, 10)
	sellPrices := sdk.DecCoins{sdk.NewInt64DecCoin(reserveToken, 100)}

	testCases := []struct {
		txFee           sdk.Dec
		exitFee         sdk.Dec
		expectedReturns sdk.Int
	}{
		{
			sdk.ZeroDec(), sdk.ZeroDec(), sdk.NewInt(1000),
		}, // (10 * 100) - (10 * FEE) = 1000, where FEE=0
		{
			sdk.NewDec(10), sdk.ZeroDec(), sdk.NewInt(900),
		}, // (10 * 100) - (10 * FEE) = 900, where FEE=10
		{
			sdk.ZeroDec(), sdk.NewDec(10), sdk.NewInt(900),
		}, // (10 * 100) - (10 * FEE) = 900, where FEE=10
		{
			sdk.NewDec(10), sdk.NewDec(10), sdk.NewInt(800),
		}, // (10 * 100) - (10 * FEE) = 800, where FEE=20
		{
			sdk.NewDec(100), sdk.NewDec(0), sdk.ZeroInt(),
		}, // (10 * 100) - (10 * FEE) = 0, where FEE=100
		{
			sdk.NewDec(100), sdk.NewDec(100), sdk.ZeroInt(),
		}, // (10 * 100) - (10 * FEE) = adjusted(-1000) = 0, where FEE=200
	}

	for _, tc := range testCases {
		// Create sell order
		so := types.NewSellOrder(sellerAddress, sellAmount)

		// Set transaction and exit fee and current supply
		war.TxFeePercentage = tc.txFee
		war.ExitFeePercentage = tc.exitFee
		war.CurrentSupply = sellAmount
		app.WarsKeeper.SetWar(ctx, war.Token, war)

		// Calculate total return
		reserveReturn := sellPrices[0].Amount.MulInt(so.Amount.Amount)
		reserveReturns := sdk.DecCoins{sdk.NewDecCoinFromDec(sellPrices[0].Denom, reserveReturn)}
		reserveReturnsRounded := types.RoundReserveReturns(reserveReturns)
		totalFees := war.GetTxFees(reserveReturns).Add(war.GetExitFees(reserveReturns)...)
		totalFees = types.AdjustFees(totalFees, reserveReturnsRounded)
		totalReturns := reserveReturnsRounded.Sub(totalFees)

		// Check expected returns
		require.Equal(t, totalReturns.AmountOf(reserveToken), tc.expectedReturns)

		// Add reserve tokens (freshly minted) paid by seller when buying to reserve
		err := app.SupplyKeeper.MintCoins(ctx, types.WarsMintBurnAccount, reserveReturnsRounded)
		require.Nil(t, err)
		err = app.WarsKeeper.DepositReserveFromModule(
			ctx, war.Token, types.WarsMintBurnAccount, reserveReturnsRounded)
		require.NoError(t, err)

		// Previous values
		prevSupplyWars := app.WarsKeeper.MustGetWar(ctx, war.Token).CurrentSupply
		prevReserveBal := app.WarsKeeper.GetReserveBalances(ctx, war.Token)
		prevFeeAddrBal := app.BankKeeper.GetCoins(ctx, war.FeeAddress)
		prevSellerBal := app.BankKeeper.GetCoins(ctx, sellerAddress)

		// Perform sell
		err = app.WarsKeeper.PerformSellAtPrice(ctx, war.Token, so, sellPrices)
		require.NoError(t, err)

		// New values
		newSupplyWars := app.WarsKeeper.MustGetWar(ctx, war.Token).CurrentSupply
		newReserveBal := app.WarsKeeper.GetReserveBalances(ctx, war.Token)
		newFeeAddrBal := app.BankKeeper.GetCoins(ctx, war.FeeAddress)
		newSellerBal := app.BankKeeper.GetCoins(ctx, sellerAddress)

		require.True(t, prevSupplyWars.Sub(so.Amount).IsEqual(newSupplyWars))
		require.Equal(t, prevReserveBal.Sub(reserveReturnsRounded), newReserveBal)
		require.Equal(t, totalFees.IsZero(), tc.txFee.IsZero() && tc.exitFee.IsZero())
		if totalFees.IsZero() {
			require.Equal(t, prevFeeAddrBal, newFeeAddrBal)
		} else {
			require.Equal(t, prevFeeAddrBal.Add(totalFees...), newFeeAddrBal)
		}
		require.Equal(t, prevSellerBal.Add(totalReturns...), newSellerBal)
	}
}

func TestPerformSwap(t *testing.T) {
	app, ctx := createTestApp(false)
	war := getValidSwapperWar()

	swapAmount := sdk.NewInt(100)

	res200 := sdk.NewInt64Coin(reserveToken, 200)
	res300 := sdk.NewInt64Coin(reserveToken, 300)

	rez200 := sdk.NewInt64Coin(reserveToken2, 200)
	rez300 := sdk.NewInt64Coin(reserveToken2, 300)

	testCases := []struct {
		fromToken              string
		toToken                string
		txFee                  sdk.Dec
		inReserve              sdk.Coin
		outReserve             sdk.Coin
		sanityRate             sdk.Dec
		sanityMarginPercentage sdk.Dec
		sanityRateViolated     bool
	}{
		{
			reserveToken, reserveToken2, sdk.ZeroDec(), res200, rez300,
			sdk.ZeroDec(), sdk.ZeroDec(), false,
		}, // 100res to rez (with initial reserve 200res,300rez) and no fee and no sanity rates
		{
			reserveToken2, reserveToken, sdk.ZeroDec(), rez200, res300,
			sdk.ZeroDec(), sdk.ZeroDec(), false,
		}, // 100rez to res (with initial reserve 300res,200rez) and no fee and no sanity rates

		{
			reserveToken, reserveToken2, sdk.NewDec(10), res200, rez300,
			sdk.ZeroDec(), sdk.ZeroDec(), false,
		}, // 100res to rez (with initial reserve 200res,300rez) with 10% fee and no sanity rates
		{
			reserveToken2, reserveToken, sdk.NewDec(10), rez200, res300,
			sdk.ZeroDec(), sdk.ZeroDec(), false,
		}, // 100rez to res (with initial reserve 300res,200rez) with 10% fee and no sanity rates

		{
			reserveToken, reserveToken2, sdk.ZeroDec(), res200, rez300,
			sdk.OneDec(), sdk.NewDec(50), false, // 1 +- 50%
		}, // 100res to rez (with initial reserve 200res,300rez) and no fee but with sanity rates not violated
		// Sanity rates are not violated since reserves will become 300res,200rez -> 300/200 -> 1.5 which is >= 1.50
		{
			reserveToken, reserveToken2, sdk.ZeroDec(), res200, rez300,
			sdk.OneDec(), sdk.NewDec(49), true, // 1 +- 49%
		}, // 100res to rez (with initial reserve 200res,300rez) and no fee but with sanity rates violated
		// Sanity rates are violated since reserves will become 300res,200rez -> 300/200 -> 1.5 which is > 1.49
	}

	for _, tc := range testCases {
		// Constant product
		cp := tc.inReserve.Amount.Mul(tc.outReserve.Amount).ToDec()

		// Create swap order
		fromAmount := sdk.NewCoin(tc.fromToken, swapAmount)
		fromAmounts := sdk.Coins{fromAmount}
		fromAmountsDec := sdk.DecCoins{sdk.NewDecCoinFromCoin(fromAmount)}
		so := types.NewSwapOrder(swapperAddress, fromAmount, tc.toToken)

		// Set transaction fee, sanity rates, and initial reserve balances
		war.TxFeePercentage = tc.txFee
		war.SanityRate = tc.sanityRate
		war.SanityMarginPercentage = tc.sanityMarginPercentage
		app.WarsKeeper.SetWar(ctx, war.Token, war)
		startingReserves := sdk.NewCoins(tc.inReserve, tc.outReserve)
		err := app.SupplyKeeper.MintCoins(ctx, types.WarsMintBurnAccount, startingReserves)
		require.Nil(t, err)
		err = app.WarsKeeper.DepositReserveFromModule(
			ctx, war.Token, types.WarsMintBurnAccount, startingReserves)
		require.NoError(t, err)

		// Add reserve tokens sent by swapper to module account address
		moduleAcc := app.SupplyKeeper.GetModuleAccount(ctx, types.BatchesIntermediaryAccount)
		err = app.BankKeeper.SetCoins(ctx, moduleAcc.GetAddress(), fromAmounts)
		require.NoError(t, err)

		// Calculations
		txFees := war.GetTxFees(fromAmountsDec)
		totalIns := fromAmounts.Sub(txFees) // into reserves
		newInReserveDec := tc.inReserve.Amount.Add(totalIns.AmountOf(tc.fromToken)).ToDec()
		newOutReserveDec := sdk.NewDecCoinFromDec(tc.toToken, cp.Quo(newInReserveDec))
		totalOuts := sdk.Coins{types.RoundReserveReturn(sdk.NewDecCoinFromCoin(tc.outReserve).Sub(newOutReserveDec))} // out of reserves (i.e. returns)

		// Previous values
		prevModuleAccBal := app.BankKeeper.GetCoins(ctx, moduleAcc.GetAddress())
		prevReserveBal := app.WarsKeeper.GetReserveBalances(ctx, war.Token)
		prevFeeAddrBal := app.BankKeeper.GetCoins(ctx, war.FeeAddress)
		prevSwapperBal := app.BankKeeper.GetCoins(ctx, swapperAddress)

		// Perform swap
		err, ok := app.WarsKeeper.PerformSwap(ctx, war.Token, so)
		require.True(t, ok)

		// Check if error due to violated sanity rate
		if tc.sanityRateViolated {
			require.Error(t, err)
			continue // app would panic at this stage
		} else {
			require.NoError(t, err)
		}

		// New values
		newModuleAccBal := app.BankKeeper.GetCoins(ctx, moduleAcc.GetAddress())
		newReserveBal := app.WarsKeeper.GetReserveBalances(ctx, war.Token)
		newFeeAddrBal := app.BankKeeper.GetCoins(ctx, war.FeeAddress)
		newSwapperBal := app.BankKeeper.GetCoins(ctx, swapperAddress)

		require.Equal(t, prevModuleAccBal.Sub(fromAmounts), newModuleAccBal)
		require.Equal(t, prevReserveBal.Add(totalIns...).Sub(totalOuts), newReserveBal)
		require.Equal(t, txFees.IsZero(), tc.txFee.IsZero())
		if txFees.IsZero() {
			require.Equal(t, prevFeeAddrBal, newFeeAddrBal)
		} else {
			require.Equal(t, prevFeeAddrBal.Add(txFees...), newFeeAddrBal)
		}
		require.Equal(t, prevSwapperBal.Add(totalOuts...), newSwapperBal)
	}
}

func TestPerformBuys(t *testing.T) {
	app, ctx := createTestApp(false)

	// Create war and batch (with no fees for simpler test)
	war := getValidWar()
	batch := getValidBatch()
	war.TxFeePercentage = sdk.ZeroDec()
	war.ExitFeePercentage = sdk.ZeroDec()
	app.WarsKeeper.SetWar(ctx, war.Token, war)
	app.WarsKeeper.SetBatch(ctx, war.Token, batch)

	buyPrices := sdk.DecCoins{sdk.NewInt64DecCoin(reserveToken, 100)}
	blankSellPrices := sdk.NewDecCoinsFromCoins() // blank
	maxPrices := sdk.Coins{sdk.NewInt64Coin(reserveToken, 2000)}

	testCases := []struct {
		amount    sdk.Int
		maxPrices sdk.Coins
	}{
		{
			sdk.NewInt(10), maxPrices,
		}, // 10 * 100 = 1000 <= 2000
		{
			sdk.NewInt(11), maxPrices,
		}, // 11 * 100 = 1100 <= 2000
		{
			sdk.NewInt(12), maxPrices,
		}, // 12 * 100 = 1200 <= 2000
	}

	// Add reserve tokens paid by buyer to module account address
	moduleAcc := app.SupplyKeeper.GetModuleAccount(ctx, types.BatchesIntermediaryAccount)

	globalTotalPrices := sdk.NewCoins()
	globalIncreaseInBuyerBal := sdk.NewCoins()
	globalTokensBought := sdk.NewCoin(war.Token, sdk.ZeroInt())

	// Add buy orders
	for _, tc := range testCases {
		// Create and add buy order
		amount := sdk.NewCoin(war.Token, tc.amount)
		bo := types.NewBuyOrder(buyerAddress, amount, tc.maxPrices)
		app.WarsKeeper.AddBuyOrder(ctx, token, bo, buyPrices, blankSellPrices)

		// Calculate total prices
		reservePrice := buyPrices[0].Amount.MulInt(bo.Amount.Amount)
		reservePrices := sdk.DecCoins{sdk.NewDecCoinFromDec(buyPrices[0].Denom, reservePrice)}
		totalPrices := types.RoundReserveReturns(reservePrices)
		globalTotalPrices = globalTotalPrices.Add(totalPrices...)

		// Add coins paid by buyer
		_, err := app.BankKeeper.AddCoins(ctx, moduleAcc.GetAddress(), tc.maxPrices)
		require.NoError(t, err)

		// Calculate increase in buyer balance
		remainderForBuyer := tc.maxPrices.Sub(totalPrices)
		tokensBought := sdk.NewCoin(war.Token, tc.amount)
		globalTokensBought = globalTokensBought.Add(tokensBought)
		globalIncreaseInBuyerBal = globalIncreaseInBuyerBal.Add(
			sdk.Coins{tokensBought}.Add(remainderForBuyer...)...)
	}

	// Perform buys
	app.WarsKeeper.PerformBuyOrders(ctx, token)

	// New values
	newSupplySDK := app.SupplyKeeper.GetSupply(ctx).GetTotal().AmountOf(war.Token)
	newSupplyWars := app.WarsKeeper.MustGetWar(ctx, war.Token).CurrentSupply
	newModuleAccBal := app.BankKeeper.GetCoins(ctx, moduleAcc.GetAddress())
	newBuyerBal := app.BankKeeper.GetCoins(ctx, buyerAddress)

	require.Equal(t, globalTokensBought.Amount, newSupplySDK)
	require.Equal(t, globalTokensBought, newSupplyWars)
	require.Equal(t, sdk.Coins(nil), newModuleAccBal)
	require.Equal(t, globalIncreaseInBuyerBal, newBuyerBal)
}

func TestPerformSells(t *testing.T) {
	app, ctx := createTestApp(false)

	// Create war and batch (with no fees for simpler test)
	war := getValidWar()
	batch := getValidBatch()
	war.TxFeePercentage = sdk.ZeroDec()
	war.ExitFeePercentage = sdk.ZeroDec()
	app.WarsKeeper.SetWar(ctx, war.Token, war)
	app.WarsKeeper.SetBatch(ctx, war.Token, batch)

	sellPrices := sdk.DecCoins{sdk.NewInt64DecCoin(reserveToken, 100)}
	blankBuyPrices := sdk.NewDecCoinsFromCoins() // blank

	testCases := []struct {
		amount sdk.Int
	}{
		{sdk.NewInt(10)}, // 10 * 100 = 1000
		{sdk.NewInt(11)}, // 11 * 100 = 1100
		{sdk.NewInt(12)}, // 12 * 100 = 1200
	}

	globalTotalReturns := sdk.NewCoins()

	// Add sell orders
	for _, tc := range testCases {
		// Create and add sell order
		amount := sdk.NewCoin(war.Token, tc.amount)
		so := types.NewSellOrder(sellerAddress, amount)
		app.WarsKeeper.AddSellOrder(ctx, token, so, blankBuyPrices, sellPrices)

		// Calculate total return
		reserveReturn := sellPrices[0].Amount.MulInt(so.Amount.Amount)
		reserveReturns := sdk.DecCoins{sdk.NewDecCoinFromDec(sellPrices[0].Denom, reserveReturn)}
		reserveReturnsRounded := types.RoundReserveReturns(reserveReturns)
		globalTotalReturns = globalTotalReturns.Add(reserveReturnsRounded...)

		// Add reserve tokens (freshly minted) paid by seller when buying to reserve
		err := app.SupplyKeeper.MintCoins(ctx, types.WarsMintBurnAccount, reserveReturnsRounded)
		require.Nil(t, err)
		err = app.WarsKeeper.DepositReserveFromModule(
			ctx, war.Token, types.WarsMintBurnAccount, reserveReturnsRounded)
		require.NoError(t, err)

		// Add increase in current supply due to a (simulated) buy
		war.CurrentSupply = war.CurrentSupply.Add(amount)
		app.WarsKeeper.SetCurrentSupply(ctx, war.Token, war.CurrentSupply)
	}

	// Perform sells
	app.WarsKeeper.PerformSellOrders(ctx, token)

	// New values
	newSupplyWars := app.WarsKeeper.MustGetWar(ctx, war.Token).CurrentSupply
	newReserveBal := app.WarsKeeper.GetReserveBalances(ctx, war.Token)
	newSellerBal := app.BankKeeper.GetCoins(ctx, sellerAddress)

	zeroWarTokens := sdk.NewCoin(war.Token, sdk.ZeroInt())
	require.Equal(t, zeroWarTokens, newSupplyWars)
	require.Equal(t, sdk.Coins(nil), newReserveBal)
	require.Equal(t, globalTotalReturns, newSellerBal)
}

func TestPerformSwaps(t *testing.T) {
	app, ctx := createTestApp(false)

	// Create war and batch (with no fees for simpler test)
	war := getValidSwapperWar()
	batch := getValidBatch()
	war.TxFeePercentage = sdk.ZeroDec()
	war.ExitFeePercentage = sdk.ZeroDec()
	war.SanityRate = sdk.OneDec()
	war.SanityMarginPercentage = sdk.NewDec(1000)
	app.WarsKeeper.SetWar(ctx, war.Token, war)
	app.WarsKeeper.SetBatch(ctx, war.Token, batch)

	// Set initial reserves
	initialInReserve := sdk.NewInt64Coin(reserveToken, 200)
	initialOutReserve := sdk.NewInt64Coin(reserveToken2, 300)
	initialReserves := sdk.NewCoins(initialInReserve, initialOutReserve)
	err := app.SupplyKeeper.MintCoins(ctx, types.WarsMintBurnAccount, initialReserves)
	require.Nil(t, err)
	err = app.WarsKeeper.DepositReserveFromModule(
		ctx, war.Token, types.WarsMintBurnAccount, initialReserves)
	require.NoError(t, err)

	testCases := []struct {
		fromToken        string
		toToken          string
		amount           sdk.Int
		willGetCancelled bool
	}{
		{
			reserveToken, reserveToken2, sdk.NewInt(100), false,
		}, // 100res to rez
		{
			reserveToken2, reserveToken, sdk.NewInt(200), false,
		}, // 200rez to res
		{
			reserveToken, reserveToken2, sdk.NewInt(300), false,
		}, // 300res to rez
		{
			reserveToken2, reserveToken, sdk.NewInt(400), false,
		}, // 400rez to res
		{
			reserveToken, reserveToken2, sdk.NewInt(1000), true,
		}, // 1000res to rez, to violate sanity
	}

	globalTotalReturns := sdk.NewCoins()
	globalIncreaseInSwapperBal := sdk.NewCoins()
	globalDecreaseInSwapperBal := sdk.NewCoins()
	globalReserveBal := initialReserves

	// Previous values
	moduleAcc := app.SupplyKeeper.GetModuleAccount(ctx, types.BatchesIntermediaryAccount)

	// Add swap orders
	for _, tc := range testCases {
		// Constant product
		inReserve := sdk.NewCoin(tc.fromToken, globalReserveBal.AmountOf(tc.fromToken))
		outReserve := sdk.NewCoin(tc.toToken, globalReserveBal.AmountOf(tc.toToken))
		cp := inReserve.Amount.Mul(outReserve.Amount).ToDec()

		// Create and add swap order
		fromAmount := sdk.NewCoin(tc.fromToken, tc.amount)
		fromAmounts := sdk.Coins{fromAmount}
		so := types.NewSwapOrder(swapperAddress, fromAmount, tc.toToken)
		app.WarsKeeper.AddSwapOrder(ctx, token, so)

		// Add reserve tokens sent by swapper to module account address
		_, err = app.BankKeeper.AddCoins(ctx, moduleAcc.GetAddress(), fromAmounts)
		require.NoError(t, err)

		if tc.willGetCancelled {
			// Reserve returned back to swapper
			globalIncreaseInSwapperBal = globalIncreaseInSwapperBal.Add(fromAmounts...)
		} else {
			// Calculations
			newInReserveDec := inReserve.Amount.Add(fromAmount.Amount).ToDec()
			newOutReserveDec := sdk.NewDecCoinFromDec(tc.toToken, cp.Quo(newInReserveDec))
			totalOuts := sdk.Coins{types.RoundReserveReturn(sdk.NewDecCoinFromCoin(outReserve).Sub(newOutReserveDec))} // out of reserves (i.e. returns)
			globalIncreaseInSwapperBal = globalIncreaseInSwapperBal.Add(totalOuts...)
			globalDecreaseInSwapperBal = globalDecreaseInSwapperBal.Add(fromAmounts...)
			globalReserveBal = globalReserveBal.Add(fromAmounts...).Sub(totalOuts)
		}
	}

	// Perform swaps
	app.WarsKeeper.PerformSwapOrders(ctx, token)

	// New balances
	newModuleAccBal := app.BankKeeper.GetCoins(ctx, moduleAcc.GetAddress())
	newReserveBal := app.WarsKeeper.GetReserveBalances(ctx, war.Token)
	newSwapperBal := app.BankKeeper.GetCoins(ctx, sellerAddress)

	require.Equal(t, sdk.Coins(nil), newModuleAccBal)
	require.Equal(t, globalReserveBal, newReserveBal)
	require.Equal(t, globalTotalReturns, newSwapperBal)
}

func TestOrderCancelled(t *testing.T) {
	// Create order and set as cancelled
	baseOrder := getValidBaseOrder()

	// Not cancelled by default
	require.False(t, baseOrder.IsCancelled())

	// Set as cancelled
	baseOrder.Cancelled = true

	// Check that cancelled
	require.True(t, baseOrder.IsCancelled())
}

func TestCheckIfBuyOrderFulfillableAtPrice(t *testing.T) {
	app, ctx := createTestApp(false)
	war := getValidWar()

	buyPrices := sdk.DecCoins{sdk.NewInt64DecCoin(reserveToken, 100)}
	maxPrices := sdk.Coins{sdk.NewInt64Coin(reserveToken, 1100)}

	testCases := []struct {
		amount           int64
		maxPrices        sdk.Coins
		txFee            sdk.Dec
		orderFulfillable bool
	}{
		{
			10, maxPrices, sdk.ZeroDec(), true,
		}, // (10 * 100) + (10 * FEE) = 1000 <= 1100, where FEE=0
		{
			11, maxPrices, sdk.ZeroDec(), true,
		}, // (11 * 100) + (11 * FEE) = 1100 <= 1100, where FEE=0
		{
			12, maxPrices, sdk.ZeroDec(), false,
		}, // (12 * 100) + (12 * FEE) = 1200 > 1100, where FEE=0
		{
			10, maxPrices, sdk.NewDec(10), true,
		}, // (10 * 100) + (10 * FEE) = 1100 <= 1100, where FEE=10
		{
			10, maxPrices, sdk.NewDec(20), false,
		}, // (10 * 100) + (10 * FEE) = 1200 > 1100, where FEE=20
	}
	for i, tc := range testCases {
		// Create buy order
		amount := sdk.NewCoin(war.Token, sdk.NewInt(tc.amount))
		bo := types.NewBuyOrder(buyerAddress, amount, tc.maxPrices)

		// Set transaction fee
		war.TxFeePercentage = tc.txFee
		app.WarsKeeper.SetWar(ctx, war.Token, war)

		err := app.WarsKeeper.CheckIfBuyOrderFulfillableAtPrice(
			ctx, war.Token, bo, buyPrices)
		require.Equal(t, tc.orderFulfillable, err == nil, "unexpected result for test case #%d", i)
	}
}

func TestCancelUnfulfillableBuys(t *testing.T) {
	app, ctx := createTestApp(false)
	war := getValidWar()

	buyPrices := sdk.DecCoins{sdk.NewInt64DecCoin(reserveToken, 100)}
	maxPrices := sdk.Coins{sdk.NewInt64Coin(reserveToken, 1100)}
	blankSellPrices := sdk.NewDecCoinsFromCoins() // blank
	zeroTokens := sdk.NewCoin(war.Token, sdk.ZeroInt())

	testCases := []struct {
		amount           int64
		maxPrices        sdk.Coins
		txFee            sdk.Dec
		orderFulfillable bool
	}{
		{
			10, maxPrices, sdk.ZeroDec(), true,
		}, // (10 * 100) + (10 * FEE) = 1000 <= 1100, where FEE=0
		{
			11, maxPrices, sdk.ZeroDec(), true,
		}, // (11 * 100) + (11 * FEE) = 1100 <= 1100, where FEE=0
		{
			12, maxPrices, sdk.ZeroDec(), false,
		}, // (12 * 100) + (12 * FEE) = 1200 > 1100, where FEE=0
		{
			10, maxPrices, sdk.NewDec(10), true,
		}, // (10 * 100) + (10 * FEE) = 1100 <= 1100, where FEE=10
		{
			10, maxPrices, sdk.NewDec(20), false,
		}, // (10 * 100) + (10 * FEE) = 1200 > 1100, where FEE=20
	}
	for _, tc := range testCases {
		// Set up war (with tx fee) and new batch
		war.TxFeePercentage = tc.txFee
		app.WarsKeeper.SetWar(ctx, war.Token, war)
		app.WarsKeeper.SetBatch(ctx, war.Token, getValidBatch())

		// Create and add buy order
		amount := sdk.NewCoin(war.Token, sdk.NewInt(tc.amount))
		bo := types.NewBuyOrder(buyerAddress, amount, tc.maxPrices)
		app.WarsKeeper.AddBuyOrder(ctx, war.Token, bo, buyPrices, blankSellPrices)

		// Add reserve tokens to module account address for return if cancel
		moduleAcc := app.SupplyKeeper.GetModuleAccount(ctx, types.BatchesIntermediaryAccount)
		_ = app.BankKeeper.SetCoins(ctx, moduleAcc.GetAddress(), tc.maxPrices)
		require.Equal(t, tc.maxPrices, app.BankKeeper.GetCoins(ctx, moduleAcc.GetAddress()))

		// Check that order added to batch and that it's not cancelled
		batch := app.WarsKeeper.MustGetBatch(ctx, war.Token)
		require.Equal(t, bo.Amount, batch.TotalBuyAmount)
		require.Len(t, batch.Buys, 1)
		require.False(t, batch.Buys[0].Cancelled)

		// Get account balance before possible cancellation
		balanceBefore := app.BankKeeper.GetCoins(ctx, buyerAddress)

		// Cancel unfulfillable buys and check amount of cancellations
		cancelledOrders := app.WarsKeeper.CancelUnfulfillableBuys(ctx, war.Token)
		if tc.orderFulfillable {
			require.Equal(t, 0, cancelledOrders)
		} else {
			require.Equal(t, 1, cancelledOrders)
		}

		// Check that batch is (un)changed based on order (un)fulfillability
		batch = app.WarsKeeper.MustGetBatch(ctx, war.Token)
		if tc.orderFulfillable {
			// Check that not cancelled
			require.Equal(t, bo.Amount, batch.TotalBuyAmount)
			require.False(t, batch.Buys[0].Cancelled)
			require.Equal(t, buyPrices, batch.BuyPrices)

			// Check that balances unchanged
			require.Equal(t, tc.maxPrices, app.BankKeeper.GetCoins(ctx, moduleAcc.GetAddress()))
			require.Equal(t, balanceBefore, app.BankKeeper.GetCoins(ctx, buyerAddress))
		} else {
			// Check that cancelled
			require.Equal(t, zeroTokens, batch.TotalBuyAmount)
			require.True(t, batch.Buys[0].Cancelled)
			require.Equal(t, buyPrices, batch.BuyPrices) // this changes only CancelUnfulfillableOrders is used

			// Check that reserve tokens returned to buyer
			newBalance := balanceBefore.Add(tc.maxPrices...)
			require.Empty(t, app.BankKeeper.GetCoins(ctx, moduleAcc.GetAddress()))
			require.Equal(t, newBalance, app.BankKeeper.GetCoins(ctx, buyerAddress))
		}
	}
}

func TestCancelUnfulfillableOrders(t *testing.T) {
	app, ctx := createTestApp(false)
	war := getValidWar()

	buyPrices := sdk.DecCoins{sdk.NewInt64DecCoin(reserveToken, 100)}
	maxPrices := sdk.Coins{sdk.NewInt64Coin(reserveToken, 1100)}
	blankSellPrices := sdk.NewDecCoinsFromCoins() // blank
	zeroTokens := sdk.NewCoin(war.Token, sdk.ZeroInt())

	testCases := []struct {
		amount           int64
		maxPrices        sdk.Coins
		txFee            sdk.Dec
		orderFulfillable bool
	}{
		{
			10, maxPrices, sdk.ZeroDec(), true,
		}, // (10 * 100) + (10 * FEE) = 1000 <= 1100, where FEE=0
		{
			11, maxPrices, sdk.ZeroDec(), true,
		}, // (11 * 100) + (11 * FEE) = 1100 <= 1100, where FEE=0
		{
			12, maxPrices, sdk.ZeroDec(), false,
		}, // (12 * 100) + (12 * FEE) = 1200 > 1100, where FEE=0
		{
			10, maxPrices, sdk.NewDec(10), true,
		}, // (10 * 100) + (10 * FEE) = 1100 <= 1100, where FEE=10
		{
			10, maxPrices, sdk.NewDec(20), false,
		}, // (10 * 100) + (10 * FEE) = 1200 > 1100, where FEE=20
	}
	for _, tc := range testCases {
		// Set up war (with tx fee and bumped-up supply) and new batch
		war.TxFeePercentage = tc.txFee
		war.CurrentSupply = sdk.NewInt64Coin(war.Token, 100)
		app.WarsKeeper.SetWar(ctx, war.Token, war)
		app.WarsKeeper.SetBatch(ctx, war.Token, getValidBatch())

		// Create and add buy order
		amount := sdk.NewCoin(war.Token, sdk.NewInt(tc.amount))
		bo := types.NewBuyOrder(buyerAddress, amount, tc.maxPrices)
		app.WarsKeeper.AddBuyOrder(ctx, war.Token, bo, buyPrices, blankSellPrices)

		// Add reserve tokens to module account address for return if cancel
		moduleAcc := app.SupplyKeeper.GetModuleAccount(ctx, types.BatchesIntermediaryAccount)
		_ = app.BankKeeper.SetCoins(ctx, moduleAcc.GetAddress(), tc.maxPrices)
		require.Equal(t, tc.maxPrices, app.BankKeeper.GetCoins(ctx, moduleAcc.GetAddress()))

		// Check that order added to batch and that it's not cancelled
		batch := app.WarsKeeper.MustGetBatch(ctx, war.Token)
		require.Equal(t, bo.Amount, batch.TotalBuyAmount)
		require.Len(t, batch.Buys, 1)
		require.False(t, batch.Buys[0].Cancelled)

		// Get account balance before possible cancellation
		balanceBefore := app.BankKeeper.GetCoins(ctx, buyerAddress)

		// Cancel unfulfillable buys and check amount of cancellations
		cancelledOrders := app.WarsKeeper.CancelUnfulfillableOrders(ctx, war.Token)
		if tc.orderFulfillable {
			require.Equal(t, 0, cancelledOrders)
		} else {
			require.Equal(t, 1, cancelledOrders)
		}

		// Check that batch is (un)changed based on order fulfillability
		batch = app.WarsKeeper.MustGetBatch(ctx, war.Token)
		if tc.orderFulfillable {
			// Check that not cancelled
			require.Equal(t, bo.Amount, batch.TotalBuyAmount)
			require.False(t, batch.Buys[0].Cancelled)
			require.Equal(t, buyPrices, batch.BuyPrices)

			// Check that balances unchanged
			require.Equal(t, tc.maxPrices, app.BankKeeper.GetCoins(ctx, moduleAcc.GetAddress()))
			require.Equal(t, balanceBefore, app.BankKeeper.GetCoins(ctx, buyerAddress))
		} else {
			// Check that cancelled
			require.Equal(t, zeroTokens, batch.TotalBuyAmount)
			require.True(t, batch.Buys[0].Cancelled)
			require.NotEqual(t, buyPrices, batch.BuyPrices)

			// Check that reserve tokens returned to buyer
			newBalance := balanceBefore.Add(tc.maxPrices...)
			require.Empty(t, app.BankKeeper.GetCoins(ctx, moduleAcc.GetAddress()))
			require.Equal(t, newBalance, app.BankKeeper.GetCoins(ctx, buyerAddress))
		}
	}
}
