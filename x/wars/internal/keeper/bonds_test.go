package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mage-war/wars/x/wars/internal/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestWarExistsSetGet(t *testing.T) {
	app, ctx := createTestApp(false)

	// Try to get war
	_, found := app.WarsKeeper.GetWar(ctx, token)

	// War doesn't exist yet
	require.False(t, found)
	require.False(t, app.WarsKeeper.WarExists(ctx, token))

	// Add war
	warAdded := getValidWar()
	app.WarsKeeper.SetWar(ctx, token, warAdded)

	// War now exists
	require.True(t, app.WarsKeeper.WarExists(ctx, token))

	// Option 1: get war
	warFetched1, found := app.WarsKeeper.GetWar(ctx, token)
	// Option 2: must get war
	warFetched2 := app.WarsKeeper.MustGetWar(ctx, token)
	// Option 2: must get war
	warFetched3 := app.WarsKeeper.MustGetWarByKey(ctx, types.GetWarKey(token))

	// Batch fetched is equal to added batch
	require.EqualValues(t, warAdded, warFetched1)
	require.EqualValues(t, warAdded, warFetched2)
	require.EqualValues(t, warAdded, warFetched3)
	require.True(t, found)
}

func TestDepositReserve(t *testing.T) {
	app, ctx := createTestApp(false)

	// Add war
	war := getValidWar()
	app.WarsKeeper.SetWar(ctx, token, war)

	// Reserve is initially empty
	require.True(t, app.WarsKeeper.GetReserveBalances(ctx, token).IsZero())

	// Add tokens to an account
	amount, err := sdk.ParseCoins("12res1,34res2")
	require.Nil(t, err)
	address := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	err = app.BankKeeper.SetCoins(ctx, address, amount)
	require.Nil(t, err)

	// Deposit reserve
	err = app.WarsKeeper.DepositReserve(ctx, token, address, amount)
	require.Nil(t, err)

	// Reserve now equal to amount sent and address balance is zero
	war = app.WarsKeeper.MustGetWar(ctx, token)
	reserveBalances := app.WarsKeeper.GetReserveBalances(ctx, token)
	require.Equal(t, amount, reserveBalances)
	addressBalance := app.BankKeeper.GetCoins(ctx, address)
	require.Empty(t, addressBalance)

	// Also confirm that reserve module account has the actual amount
	moduleAddr := app.SupplyKeeper.GetModuleAddress(types.WarsReserveAccount)
	addressBalance = app.BankKeeper.GetCoins(ctx, moduleAddr)
	require.Equal(t, amount, addressBalance)
}

func TestDepositReserveFromModule(t *testing.T) {
	app, ctx := createTestApp(false)

	// Add war
	war := getValidWar()
	app.WarsKeeper.SetWar(ctx, token, war)

	// Reserve is initially empty
	require.True(t, app.WarsKeeper.GetReserveBalances(ctx, token).IsZero())

	// Mint tokens to a module
	amount, err := sdk.ParseCoins("12res1,34res2")
	require.Nil(t, err)
	err = app.SupplyKeeper.MintCoins(ctx, types.WarsMintBurnAccount, amount)
	require.Nil(t, err)

	// Deposit reserve
	err = app.WarsKeeper.DepositReserveFromModule(
		ctx, token, types.WarsMintBurnAccount, amount)
	require.Nil(t, err)

	// Reserve now equal to amount sent and module address balance is zero
	war = app.WarsKeeper.MustGetWar(ctx, token)
	reserveBalances := app.WarsKeeper.GetReserveBalances(ctx, token)
	require.Equal(t, amount, reserveBalances)
	moduleAddr := app.SupplyKeeper.GetModuleAddress(types.WarsMintBurnAccount)
	addressBalance := app.BankKeeper.GetCoins(ctx, moduleAddr)
	require.Empty(t, addressBalance)

	// Also confirm that reserve module account has the actual amount
	moduleAddr = app.SupplyKeeper.GetModuleAddress(types.WarsReserveAccount)
	addressBalance = app.BankKeeper.GetCoins(ctx, moduleAddr)
	require.Equal(t, amount, addressBalance)
}

func TestWithdrawReserve(t *testing.T) {
	app, ctx := createTestApp(false)

	// Add war
	war := getValidWar()
	app.WarsKeeper.SetWar(ctx, token, war)

	// Reserve is initially empty
	require.True(t, app.WarsKeeper.GetReserveBalances(ctx, token).IsZero())

	// Simulate depositing reserve
	amount, err := sdk.ParseCoins("12res1,34res2")
	require.Nil(t, err)
	err = app.SupplyKeeper.MintCoins(ctx, types.WarsMintBurnAccount, amount)
	require.Nil(t, err)
	err = app.SupplyKeeper.SendCoinsFromModuleToModule(
		ctx, types.WarsMintBurnAccount, types.WarsReserveAccount, amount)
	require.Nil(t, err)
	war.CurrentReserve = amount
	app.WarsKeeper.SetWar(ctx, token, war)

	// Withdraw reserve
	address := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	err = app.WarsKeeper.WithdrawReserve(ctx, token, address, amount)
	require.Nil(t, err)

	// Reserve now zero sent and address balance is equal to amount
	war = app.WarsKeeper.MustGetWar(ctx, token)
	reserveBalances := app.WarsKeeper.GetReserveBalances(ctx, token)
	require.Empty(t, reserveBalances)
	addressBalance := app.BankKeeper.GetCoins(ctx, address)
	require.Equal(t, amount, addressBalance)

	// Also confirm that reserve module account is now empty
	moduleAddr := app.SupplyKeeper.GetModuleAddress(types.WarsReserveAccount)
	addressBalance = app.BankKeeper.GetCoins(ctx, moduleAddr)
	require.Empty(t, addressBalance)
}

func TestGetReserveBalances(t *testing.T) {
	app, ctx := createTestApp(false)

	// Add war
	war := getValidWar()
	app.WarsKeeper.SetWar(ctx, token, war)

	// Reserve is initially empty
	require.True(t, app.WarsKeeper.GetReserveBalances(ctx, token).IsZero())

	// Set war reserve
	var err error
	war.CurrentReserve, err = sdk.ParseCoins("12res1,34res2")
	require.Nil(t, err)
	app.WarsKeeper.SetWar(ctx, token, war)

	// Reserve now equal to amount sent
	reserveBalances := app.WarsKeeper.GetReserveBalances(ctx, token)
	require.Equal(t, war.CurrentReserve, reserveBalances)
}

func TestGetSupplyAdjustedForBuy(t *testing.T) {
	app, ctx := createTestApp(false)

	// Add war and batch
	war := getValidWar()
	batch := getValidBatch()
	app.WarsKeeper.SetWar(ctx, token, war)
	app.WarsKeeper.SetBatch(ctx, token, batch)

	// Supply is initially zero
	require.True(t, app.WarsKeeper.GetSupplyAdjustedForBuy(ctx, token).IsZero())

	// Increase current supply
	increaseInSupply := sdk.NewInt64Coin(token, 100)
	war.CurrentSupply = increaseInSupply
	app.WarsKeeper.SetWar(ctx, token, war)

	// Supply has increased
	supply := app.WarsKeeper.GetSupplyAdjustedForBuy(ctx, token)
	require.Equal(t, increaseInSupply, supply)

	// Increase supply by adding a buy order
	increaseDueToOrder := sdk.NewInt64Coin(token, 11)
	buyOrder := getValidBuyOrder()
	buyOrder.Amount = increaseDueToOrder
	app.WarsKeeper.AddBuyOrder(ctx, token, buyOrder, nil, nil)

	// Supply has increased
	expectedSupply := increaseInSupply.Add(increaseDueToOrder)
	supply = app.WarsKeeper.GetSupplyAdjustedForBuy(ctx, token)
	require.Equal(t, expectedSupply, supply)

	// Adding sell order does not affect supply
	sellOrder := getValidSellOrder()
	sellOrder.Amount = sdk.NewInt64Coin(token, 100)
	app.WarsKeeper.AddSellOrder(ctx, token, sellOrder, nil, nil)

	// Supply has not increased
	supply = app.WarsKeeper.GetSupplyAdjustedForBuy(ctx, token)
	require.Equal(t, expectedSupply, supply)

	// Adding swap order does not affect supply
	app.WarsKeeper.AddSwapOrder(ctx, token, getValidSwapOrder())

	// Supply has not increased
	supply = app.WarsKeeper.GetSupplyAdjustedForBuy(ctx, token)
	require.Equal(t, expectedSupply, supply)
}

func TestGetSupplyAdjustedForSell(t *testing.T) {
	app, ctx := createTestApp(false)

	// Add war and batch
	war := getValidWar()
	batch := getValidBatch()
	app.WarsKeeper.SetWar(ctx, token, war)
	app.WarsKeeper.SetBatch(ctx, token, batch)

	// Supply is initially zero
	require.True(t, app.WarsKeeper.GetSupplyAdjustedForSell(ctx, token).IsZero())

	// Increase current supply
	increaseInSupply := sdk.NewInt64Coin(token, 100)
	war.CurrentSupply = increaseInSupply
	app.WarsKeeper.SetWar(ctx, token, war)

	// Supply has increased
	supply := app.WarsKeeper.GetSupplyAdjustedForSell(ctx, token)
	require.Equal(t, increaseInSupply, supply)

	// Decrease supply by adding a sell order
	decreaseDueToOrder := sdk.NewInt64Coin(token, 11)
	sellOrder := getValidSellOrder()
	sellOrder.Amount = decreaseDueToOrder
	app.WarsKeeper.AddSellOrder(ctx, token, sellOrder, nil, nil)

	// Supply has decreased
	expectedSupply := increaseInSupply.Sub(decreaseDueToOrder)
	supply = app.WarsKeeper.GetSupplyAdjustedForSell(ctx, token)
	require.Equal(t, expectedSupply, supply)

	// Adding buy order does not affect supply
	buyOrder := getValidBuyOrder()
	buyOrder.Amount = sdk.NewInt64Coin(token, 100)
	app.WarsKeeper.AddBuyOrder(ctx, token, buyOrder, nil, nil)

	// Supply has not increased
	supply = app.WarsKeeper.GetSupplyAdjustedForSell(ctx, token)
	require.Equal(t, expectedSupply, supply)

	// Adding swap order does not affect supply
	app.WarsKeeper.AddSwapOrder(ctx, token, getValidSwapOrder())

	// Supply has not increased
	supply = app.WarsKeeper.GetSupplyAdjustedForSell(ctx, token)
	require.Equal(t, expectedSupply, supply)
}

func TestSetCurrentSupply(t *testing.T) {
	app, ctx := createTestApp(false)

	// Add war
	war := getValidWar()
	app.WarsKeeper.SetWar(ctx, token, war)

	// Supply is initially zero
	require.True(t, app.WarsKeeper.MustGetWar(ctx, token).CurrentSupply.IsZero())

	// Change supply
	newSupply := sdk.NewInt64Coin(token, 100)
	app.WarsKeeper.SetCurrentSupply(ctx, token, newSupply)

	// Check that supply changed
	supplyFetched := app.WarsKeeper.MustGetWar(ctx, token).CurrentSupply
	require.Equal(t, newSupply, supplyFetched)

	// Change supply again
	newSupply = sdk.NewInt64Coin(token, 50)
	app.WarsKeeper.SetCurrentSupply(ctx, token, newSupply)

	// Check that supply changed
	supplyFetched = app.WarsKeeper.MustGetWar(ctx, token).CurrentSupply
	require.Equal(t, newSupply, supplyFetched)
}

func TestSetWarState(t *testing.T) {
	app, ctx := createTestApp(false)

	// Add war
	war := getValidWar()
	app.WarsKeeper.SetWar(ctx, token, war)

	// State is initially "initState"
	require.Equal(t, initState, app.WarsKeeper.MustGetWar(ctx, token).State)

	// Change state
	newState := "some_other_state"
	app.WarsKeeper.SetWarState(ctx, token, newState)

	// Check that state changed
	stateFetched := app.WarsKeeper.MustGetWar(ctx, token).State
	require.Equal(t, newState, stateFetched)

	// Change supply again
	newState = "yet another state"
	app.WarsKeeper.SetWarState(ctx, token, newState)

	// Check that supply changed
	stateFetched = app.WarsKeeper.MustGetWar(ctx, token).State
	require.Equal(t, newState, stateFetched)
}
