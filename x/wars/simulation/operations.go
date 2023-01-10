package simulation

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/mage-war/wars/x/wars/internal/keeper"
	"github.com/mage-war/wars/x/wars/internal/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"math/rand"
)

// Simulation operation weights constants
const (
	OpWeightMsgCreateWar = "op_weight_msg_create_war"
	OpWeightMsgEditWar   = "op_weight_msg_edit_war"
	OpWeightMsgBuy        = "op_weight_msg_buy"
	OpWeightMsgSell       = "op_weight_msg_sell"
	OpWeightMsgSwap       = "op_weight_msg_swap"

	DefaultWeightMsgCreateWar = 5
	DefaultWeightMsgEditWar   = 5
	DefaultWeightMsgBuy        = 100
	DefaultWeightMsgSell       = 100
	DefaultWeightMsgSwap       = 100
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simulation.AppParams, cdc *codec.Codec,
	ak auth.AccountKeeper, k keeper.Keeper,
) simulation.WeightedOperations {

	var weightMsgCreateWar int
	appParams.GetOrGenerate(cdc, OpWeightMsgCreateWar, &weightMsgCreateWar, nil,
		func(_ *rand.Rand) {
			weightMsgCreateWar = DefaultWeightMsgCreateWar
		},
	)

	var weightMsgEditWar int
	appParams.GetOrGenerate(cdc, OpWeightMsgEditWar, &weightMsgEditWar, nil,
		func(_ *rand.Rand) {
			weightMsgEditWar = DefaultWeightMsgEditWar
		},
	)

	var weightMsgBuy int
	appParams.GetOrGenerate(cdc, OpWeightMsgBuy, &weightMsgBuy, nil,
		func(_ *rand.Rand) {
			weightMsgBuy = DefaultWeightMsgBuy
		},
	)

	var weightMsgSell int
	appParams.GetOrGenerate(cdc, OpWeightMsgSell, &weightMsgSell, nil,
		func(_ *rand.Rand) {
			weightMsgSell = DefaultWeightMsgSell
		},
	)

	var weightMsgSwap int
	appParams.GetOrGenerate(cdc, OpWeightMsgSwap, &weightMsgSwap, nil,
		func(_ *rand.Rand) {
			weightMsgSwap = DefaultWeightMsgSwap
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreateWar,
			SimulateMsgCreateWar(ak),
		),
		simulation.NewWeightedOperation(
			weightMsgEditWar,
			SimulateMsgEditWar(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgBuy,
			SimulateMsgBuy(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgSell,
			SimulateMsgSell(ak, k),
		),
		simulation.NewWeightedOperation(
			weightMsgSwap,
			SimulateMsgSwap(ak, k),
		),
	}
}

func SimulateMsgCreateWar(ak auth.AccountKeeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, chainID string) (opMsg simulation.OperationMsg, fOpt []simulation.FutureOperation, err error) {

		if totalWarCount >= maxWarCount {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		simAccount, _ := simulation.RandomAcc(r, accs)
		address := simAccount.Address
		account := ak.GetAccount(ctx, address)

		token := getNextWarName()
		name := getRandomNonEmptyString(r)
		desc := getRandomNonEmptyString(r)

		creator := address
		signers := []sdk.AccAddress{creator}

		functionType := getRandomFunctionType(r)

		var reserveTokens []string
		switch functionType {
		case types.SwapperFunction:
			reserveToken1, ok1 := getRandomWarName(r)
			reserveToken2, ok2 := getRandomWarNameExcept(r, reserveToken1)
			if !ok1 || !ok2 {
				return simulation.NoOpMsg(types.ModuleName), nil, nil
			}
			reserveTokens = []string{reserveToken1, reserveToken2}
		default:
			reserveTokens = defaultReserveTokens
		}
		functionParameters := getRandomFunctionParameters(r, functionType, false)

		// Max fee is 100, so exit fee uses 100-txFee as max
		txFeePercentage := simulation.RandomDecAmount(r, sdk.NewDec(100))
		exitFeePercentage := simulation.RandomDecAmount(r, sdk.NewDec(100).Sub(txFeePercentage))

		// Since 100 is not allowed, a small number is subtracted from one of the fees
		if txFeePercentage.Add(exitFeePercentage).Equal(sdk.NewDec(100)) {
			if txFeePercentage.GT(sdk.ZeroDec()) {
				txFeePercentage = txFeePercentage.Sub(sdk.MustNewDecFromStr("0.000000000000000001"))
			} else {
				exitFeePercentage = exitFeePercentage.Sub(sdk.MustNewDecFromStr("0.000000000000000001"))
			}
		}

		// Addresses
		feeAddress := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

		// Max supply, allow sells, batch blocks
		maxSupply := sdk.NewCoin(token, sdk.NewInt(int64(
			simulation.RandIntBetween(r, 1000000, 1000000000))))
		allowSells := getRandomAllowSellsValue(r)
		batchBlocks := sdk.NewUint(uint64(
			simulation.RandIntBetween(r, 1, 10)))

		msg := types.NewMsgCreateWar(token, name, desc, creator, functionType,
			functionParameters, reserveTokens, txFeePercentage, exitFeePercentage,
			feeAddress, maxSupply, blankOrderQuantityLimits, blankSanityRate,
			blankSanityMarginPercentage, allowSells, signers, batchBlocks, blankOutcomePayment)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(types.ModuleName), nil,
				fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			sdk.Coins{},
			gas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		incrementWarCount() // since successfully created
		if msg.FunctionType == types.SwapperFunction {
			newSwapperWar(msg.Token)
		}
		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func SimulateMsgEditWar(ak auth.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, chainID string) (opMsg simulation.OperationMsg, fOpt []simulation.FutureOperation, err error) {

		// Get random war
		token, ok := getRandomWarName(r)
		if !ok {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		war, found := k.GetWar(ctx, token)
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		name := getRandomNonEmptyString(r)
		desc := getRandomNonEmptyString(r)

		simAccount, _ := simulation.FindAccount(accs, war.Creator)
		address := simAccount.Address
		account := ak.GetAccount(ctx, address)

		editor := address
		signers := []sdk.AccAddress{editor}

		msg := types.NewMsgEditWar(token, name, desc,
			types.DoNotModifyField, types.DoNotModifyField,
			types.DoNotModifyField, editor, signers)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(types.ModuleName), nil,
				fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			sdk.Coins{},
			gas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func getBuyIntoSwapper(r *rand.Rand, ctx sdk.Context, k keeper.Keeper,
	war types.War, account exported.Account) (msg types.MsgBuy, err error, ok bool) {
	address := account.GetAddress()
	spendable := account.SpendableCoins(ctx.BlockTime())

	// Come up with max prices based on what is spendable
	spendableReserve1 := spendable.AmountOf(war.ReserveTokens[0])
	spendableReserve2 := spendable.AmountOf(war.ReserveTokens[1])
	maxPriceInt1, err := simulation.RandPositiveInt(r, spendableReserve1)
	if err != nil {
		return types.MsgBuy{}, err, false
	}
	maxPriceInt2, err := simulation.RandPositiveInt(r, spendableReserve2)
	if err != nil {
		return types.MsgBuy{}, err, false
	}
	maxPrices := sdk.NewCoins(
		sdk.NewCoin(war.ReserveTokens[0], maxPriceInt1),
		sdk.NewCoin(war.ReserveTokens[1], maxPriceInt2),
	)

	// Get lesser of max possible increase in supply and max order quantity
	var maxBuyAmount sdk.Int
	maxIncreaseInSupply := war.MaxSupply.Sub(war.CurrentSupply).Amount
	maxOrderQuantity := war.OrderQuantityLimits.AmountOf(war.Token)
	if maxOrderQuantity.IsZero() {
		maxBuyAmount = maxIncreaseInSupply
	} else {
		maxBuyAmount = sdk.MinInt(maxIncreaseInSupply, maxOrderQuantity)
	}

	if maxBuyAmount.IsZero() {
		return types.MsgBuy{}, nil, false
	}

	toBuyInt, err := simulation.RandPositiveInt(r, maxBuyAmount)
	if err != nil {
		return types.MsgBuy{}, err, false
	}
	amountToBuy := sdk.NewCoin(war.Token, toBuyInt)

	// If not the first buy, create order and check if can afford
	if war.CurrentSupply.IsPositive() {
		_, _, err = k.GetUpdatedBatchPricesAfterBuy(ctx, war.Token,
			types.NewBuyOrder(address, amountToBuy, maxPrices))
		if err != nil {
			return types.MsgBuy{}, err, true
		}
	}

	return types.NewMsgBuy(address, amountToBuy, maxPrices), nil, true
}

func getBuyIntoNonSwapper(r *rand.Rand, ctx sdk.Context, k keeper.Keeper,
	war types.War, account exported.Account) (msg types.MsgBuy, err error, ok bool) {
	address := account.GetAddress()
	spendable := account.SpendableCoins(ctx.BlockTime())

	// Come up with max price based on what is spendable
	spendableReserve := spendable.AmountOf(war.ReserveTokens[0])
	maxPriceInt, err := simulation.RandPositiveInt(r, spendableReserve)
	if err != nil {
		return types.MsgBuy{}, err, false
	}
	maxPrices := sdk.Coins{sdk.NewCoin(war.ReserveTokens[0], maxPriceInt)}

	// Get lesser of max possible increase in supply and max order quantity
	var maxBuyAmount sdk.Int
	maxIncreaseInSupply := war.MaxSupply.Sub(war.CurrentSupply).Amount
	maxOrderQuantity := war.OrderQuantityLimits.AmountOf(war.Token)
	if maxOrderQuantity.IsZero() {
		maxBuyAmount = maxIncreaseInSupply
	} else {
		maxBuyAmount = sdk.MinInt(maxIncreaseInSupply, maxOrderQuantity)
	}

	if maxBuyAmount.IsZero() {
		return types.MsgBuy{}, nil, false
	}

	// For augmented function in hatch state, flip a coin to decide whether
	// to just buy all the remaining tokens up to S0 to go to the open state.
	// This is only possible if the account balance is >= remaining up to S0.
	// Otherwise, pick a random amount with the remaining to S0 as the max.
	//
	// If not an augmented function in hatch state, just pick a random amount.
	var toBuyInt sdk.Int
	if war.FunctionType == types.AugmentedFunction && war.State == types.HatchState {
		S0 := war.FunctionParameters.AsMap()["S0"].Ceil().TruncateInt()
		remainingForS0 := S0.Sub(war.CurrentSupply.Amount)
		if remainingForS0.LTE(maxBuyAmount) && simulation.RandIntBetween(r, 1, 2) == 1 {
			toBuyInt = remainingForS0
		} else if remainingForS0.GT(sdk.ZeroInt()) {
			toBuyInt, err = simulation.RandPositiveInt(r, remainingForS0)
			if err != nil {
				return types.MsgBuy{}, err, false
			}
		} else {
			panic("current cannot be equal to S0 in hatch phase")
		}
	} else {
		toBuyInt, err = simulation.RandPositiveInt(r, maxBuyAmount)
		if err != nil {
			return types.MsgBuy{}, err, false
		}
	}
	toBuy := sdk.NewCoin(war.Token, toBuyInt)

	// Create order and check if can afford
	_, _, err = k.GetUpdatedBatchPricesAfterBuy(ctx, war.Token,
		types.NewBuyOrder(address, toBuy, maxPrices))
	if err != nil {
		return types.MsgBuy{}, err, true
	}

	return types.NewMsgBuy(address, toBuy, maxPrices), nil, true
}

func SimulateMsgBuy(ak auth.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, chainID string) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		// Get random war
		token, ok := getRandomWarName(r)
		if !ok {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		war, found := k.GetWar(ctx, token)
		if !found {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		// Get accounts that have ALL the reserve tokens
		var filteredAccs []simulation.Account
		dummyNonZeroReserve := getDummyNonZeroReserve(war.ReserveTokens)
		for _, a := range accs {
			coins := ak.GetAccount(ctx, a.Address).SpendableCoins(ctx.BlockTime())
			if dummyNonZeroReserve.DenomsSubsetOf(coins) {
				filteredAccs = append(filteredAccs, a)
			}
		}

		if len(filteredAccs) == 0 {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		simAccount, _ := simulation.RandomAcc(r, filteredAccs)
		account := ak.GetAccount(ctx, simAccount.Address)

		var msg types.MsgBuy
		if war.FunctionType == types.SwapperFunction {
			msg, err, ok = getBuyIntoSwapper(r, ctx, k, war, account)
		} else {
			msg, err, ok = getBuyIntoNonSwapper(r, ctx, k, war, account)
		}

		// If ok, err is not something that should stop the simulation
		if (err != nil && ok) || (err == nil && !ok) {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		} else if err != nil && !ok {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		} else if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			sdk.Coins{},
			gas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func SimulateMsgSell(ak auth.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, chainID string) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		// Get random war
		token, ok := getRandomWarName(r)
		if !ok {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		war, found := k.GetWar(ctx, token)
		if !found || !war.AllowSells ||
			war.CurrentSupply.IsZero() ||
			war.State == types.HatchState {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		// Get accounts that have the token to be sold
		var filteredAccs []simulation.Account
		for _, a := range accs {
			coins := ak.GetAccount(ctx, a.Address).SpendableCoins(ctx.BlockTime())
			if coins.AmountOf(war.Token).IsPositive() {
				filteredAccs = append(filteredAccs, a)
			}
		}

		if len(filteredAccs) == 0 {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		simAccount, _ := simulation.RandomAcc(r, filteredAccs)
		address := simAccount.Address
		account := ak.GetAccount(ctx, address)
		amount := account.SpendableCoins(ctx.BlockTime()).AmountOf(war.Token)

		toSellInt, err := simulation.RandPositiveInt(r, amount)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		amountToSell := sdk.NewCoin(war.Token, toSellInt)

		msg := types.NewMsgSell(address, amountToSell)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			sdk.Coins{},
			gas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

func SimulateMsgSwap(ak auth.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, chainID string) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		// Get swapper function wars with some reserve
		var filteredWars []string
		for _, sbToken := range swapperWars {
			if !k.GetReserveBalances(ctx, sbToken).IsZero() {
				filteredWars = append(filteredWars, sbToken)
			}
		}

		if len(filteredWars) == 0 {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		// Get random war
		token := filteredWars[simulation.RandIntBetween(r, 0, len(filteredWars))]
		war := k.MustGetWar(ctx, token)

		fromIndex := simulation.RandIntBetween(r, 0, 1)
		toIndex := 1 - fromIndex

		fromToken := war.ReserveTokens[fromIndex]
		toToken := war.ReserveTokens[toIndex]

		// Get accounts that have the token to be swapped
		var filteredAccs []simulation.Account
		for _, a := range accs {
			coins := ak.GetAccount(ctx, a.Address).SpendableCoins(ctx.BlockTime())
			if coins.AmountOf(fromToken).IsPositive() {
				filteredAccs = append(filteredAccs, a)
			}
		}

		if len(filteredAccs) == 0 {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		simAccount, _ := simulation.RandomAcc(r, filteredAccs)
		address := simAccount.Address
		account := ak.GetAccount(ctx, address)
		fromBalance := account.SpendableCoins(ctx.BlockTime()).AmountOf(fromToken)

		toSwapInt, err := simulation.RandPositiveInt(r, fromBalance)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		amountToSwap := sdk.NewCoin(fromToken, toSwapInt)

		msg := types.NewMsgSwap(address, token, amountToSwap, toToken)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			sdk.Coins{},
			gas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		_, _, err = app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}
