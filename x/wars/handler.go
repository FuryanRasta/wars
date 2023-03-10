package wars

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mage-war/wars/x/wars/internal/keeper"
	"github.com/mage-war/wars/x/wars/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"strconv"
	"strings"
)

func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgCreateWar:
			return handleMsgCreateWar(ctx, keeper, msg)
		case types.MsgEditWar:
			return handleMsgEditWar(ctx, keeper, msg)
		case types.MsgBuy:
			return handleMsgBuy(ctx, keeper, msg)
		case types.MsgSell:
			return handleMsgSell(ctx, keeper, msg)
		case types.MsgSwap:
			return handleMsgSwap(ctx, keeper, msg)
		case types.MsgMakeOutcomePayment:
			return handleMsgMakeOutcomePayment(ctx, keeper, msg)
		case types.MsgWithdrawShare:
			return handleMsgWithdrawShare(ctx, keeper, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Unrecognized wars Msg type: %v", msg.Type())
		}
	}
}

func EndBlocker(ctx sdk.Context, keeper keeper.Keeper) []abci.ValidatorUpdate {

	iterator := keeper.GetWarIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		war := keeper.MustGetWarByKey(ctx, iterator.Key())
		batch := keeper.MustGetBatch(ctx, war.Token)

		// Subtract one block
		batch.BlocksRemaining = batch.BlocksRemaining.SubUint64(1)
		keeper.SetBatch(ctx, war.Token, batch)

		// If blocks remaining > 0 do not perform orders
		if !batch.BlocksRemaining.IsZero() {
			continue
		}

		// Perform orders
		keeper.PerformOrders(ctx, war.Token)

		// Get war again just in case current supply was updated
		// Get batch again just in case orders were cancelled
		war = keeper.MustGetWar(ctx, war.Token)
		batch = keeper.MustGetBatch(ctx, war.Token)

		// For augmented, if hatch phase and newSupply >= S0, go to open phase
		if war.FunctionType == types.AugmentedFunction &&
			war.State == types.HatchState {
			args := war.FunctionParameters.AsMap()
			if war.CurrentSupply.Amount.ToDec().GTE(args["S0"]) {
				keeper.SetWarState(ctx, war.Token, types.OpenState)
				war = keeper.MustGetWar(ctx, war.Token) // get war again
				war.AllowSells = true                     // enable sells
				keeper.SetWar(ctx, war.Token, war)      // update war
			}
		}

		// Save current batch as last batch and reset current batch
		keeper.SetLastBatch(ctx, war.Token, batch)
		keeper.SetBatch(ctx, war.Token, types.NewBatch(war.Token, war.BatchBlocks))
	}
	return []abci.ValidatorUpdate{}
}

func handleMsgCreateWar(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgCreateWar) (*sdk.Result, error) {
	if keeper.BankKeeper.BlacklistedAddr(msg.FeeAddress) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive transactions", msg.FeeAddress)
	}

	// Check that war and war DID do not already exist
	if keeper.WarExists(ctx, msg.Token) {
		return nil, sdkerrors.Wrap(types.ErrWarAlreadyExists, msg.Token)
	} else if msg.Token == keeper.StakingKeeper.GetParams(ctx).WarDenom {
		return nil, sdkerrors.Wrap(types.ErrWarTokenCannotBeStakingToken, msg.Token)
	}

	// Check that war token not reserved
	if keeper.ReservedWarToken(ctx, msg.Token) {
		return nil, types.ErrReservedWarToken
	}

	// Set state to open by default (overridden below if augmented function)
	state := types.OpenState

	// If augmented, add R0, S0, V0 as parameters for quick access
	// Also, override AllowSells and set to False if S0 > 0
	if msg.FunctionType == types.AugmentedFunction {
		paramsMap := msg.FunctionParameters.AsMap()
		d0, _ := paramsMap["d0"]
		p0, _ := paramsMap["p0"]
		theta, _ := paramsMap["theta"]
		kappa, _ := paramsMap["kappa"]

		R0 := d0.Mul(sdk.OneDec().Sub(theta))
		S0 := d0.Quo(p0)
		V0 := types.Invariant(R0, S0, kappa.TruncateInt64())
		// TODO: consider calculating these on-the-fly, especially R0 and S0

		msg.FunctionParameters = append(msg.FunctionParameters,
			types.FunctionParams{
				types.NewFunctionParam("R0", R0),
				types.NewFunctionParam("S0", S0),
				types.NewFunctionParam("V0", V0),
			}...)

		// Set state to Hatch and disable sells. Note that it is never the case
		// that we start with OpenState because S0>0, since S0=d0/p0 and d0>0
		state = types.HatchState
		msg.AllowSells = false
	}

	war := types.NewWar(msg.Token, msg.Name, msg.Description, msg.Creator,
		msg.FunctionType, msg.FunctionParameters, msg.ReserveTokens,
		msg.TxFeePercentage, msg.ExitFeePercentage, msg.FeeAddress,
		msg.MaxSupply, msg.OrderQuantityLimits, msg.SanityRate,
		msg.SanityMarginPercentage, msg.AllowSells, msg.Signers,
		msg.BatchBlocks, msg.OutcomePayment, state)

	keeper.SetWar(ctx, msg.Token, war)
	keeper.SetBatch(ctx, msg.Token, types.NewBatch(war.Token, msg.BatchBlocks))

	logger := keeper.Logger(ctx)
	logger.Info(fmt.Sprintf("war %s [%s] with reserve(s) [%s] created by %s", msg.Token,
		msg.FunctionType, strings.Join(war.ReserveTokens, ","), msg.Creator.String()))

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateWar,
			sdk.NewAttribute(types.AttributeKeyWar, msg.Token),
			sdk.NewAttribute(types.AttributeKeyName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyDescription, msg.Description),
			sdk.NewAttribute(types.AttributeKeyFunctionType, msg.FunctionType),
			sdk.NewAttribute(types.AttributeKeyFunctionParameters, msg.FunctionParameters.String()),
			sdk.NewAttribute(types.AttributeKeyReserveTokens, types.StringsToString(msg.ReserveTokens)),
			sdk.NewAttribute(types.AttributeKeyTxFeePercentage, msg.TxFeePercentage.String()),
			sdk.NewAttribute(types.AttributeKeyExitFeePercentage, msg.ExitFeePercentage.String()),
			sdk.NewAttribute(types.AttributeKeyFeeAddress, msg.FeeAddress.String()),
			sdk.NewAttribute(types.AttributeKeyMaxSupply, msg.MaxSupply.String()),
			sdk.NewAttribute(types.AttributeKeyOrderQuantityLimits, msg.OrderQuantityLimits.String()),
			sdk.NewAttribute(types.AttributeKeySanityRate, msg.SanityRate.String()),
			sdk.NewAttribute(types.AttributeKeySanityMarginPercentage, msg.SanityMarginPercentage.String()),
			sdk.NewAttribute(types.AttributeKeyAllowSells, strconv.FormatBool(msg.AllowSells)),
			sdk.NewAttribute(types.AttributeKeySigners, types.AccAddressesToString(msg.Signers)),
			sdk.NewAttribute(types.AttributeKeyBatchBlocks, msg.BatchBlocks.String()),
			sdk.NewAttribute(types.AttributeKeyOutcomePayment, msg.OutcomePayment.String()),
			sdk.NewAttribute(types.AttributeKeyState, state),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgEditWar(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgEditWar) (*sdk.Result, error) {

	war, found := keeper.GetWar(ctx, msg.Token)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrWarDoesNotExist, msg.Token)
	}

	if !war.SignersEqualTo(msg.Signers) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "list of signers does not match the one in the war")
	}

	if msg.Name != types.DoNotModifyField {
		war.Name = msg.Name
	}
	if msg.Description != types.DoNotModifyField {
		war.Description = msg.Description
	}

	if msg.OrderQuantityLimits != types.DoNotModifyField {
		orderQuantityLimits, err := sdk.ParseCoins(msg.OrderQuantityLimits)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, err.Error())
		}
		war.OrderQuantityLimits = orderQuantityLimits
	}

	if msg.SanityRate != types.DoNotModifyField {
		var sanityRate, sanityMarginPercentage sdk.Dec
		if msg.SanityRate == "" {
			sanityRate = sdk.ZeroDec()
			sanityMarginPercentage = sdk.ZeroDec()
		} else {
			parsedSanityRate, err := sdk.NewDecFromStr(msg.SanityRate)
			if err != nil {
				return nil, sdkerrors.Wrap(types.ErrArgumentMissingOrNonFloat, "sanity rate")
			} else if parsedSanityRate.IsNegative() {
				return nil, sdkerrors.Wrap(types.ErrArgumentCannotBeNegative, "sanity rate")
			}
			parsedSanityMarginPercentage, err := sdk.NewDecFromStr(msg.SanityMarginPercentage)
			if err != nil {
				return nil, sdkerrors.Wrap(types.ErrArgumentMissingOrNonFloat, "sanity margin percentage")
			} else if parsedSanityMarginPercentage.IsNegative() {
				return nil, sdkerrors.Wrap(types.ErrArgumentCannotBeNegative, "sanity margin percentage")
			}
			sanityRate = parsedSanityRate
			sanityMarginPercentage = parsedSanityMarginPercentage
		}
		war.SanityRate = sanityRate
		war.SanityMarginPercentage = sanityMarginPercentage
	}

	logger := keeper.Logger(ctx)
	logger.Info(fmt.Sprintf("war %s edited by %s",
		msg.Token, msg.Editor.String()))

	keeper.SetWar(ctx, msg.Token, war)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditWar,
			sdk.NewAttribute(types.AttributeKeyWar, msg.Token),
			sdk.NewAttribute(types.AttributeKeyName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyDescription, msg.Description),
			sdk.NewAttribute(types.AttributeKeyOrderQuantityLimits, msg.OrderQuantityLimits),
			sdk.NewAttribute(types.AttributeKeySanityRate, msg.SanityRate),
			sdk.NewAttribute(types.AttributeKeySanityMarginPercentage, msg.SanityMarginPercentage),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Editor.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgBuy(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBuy) (*sdk.Result, error) {

	token := msg.Amount.Denom
	war, found := keeper.GetWar(ctx, token)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrWarDoesNotExist, token)
	}

	// Check current state is HATCH/OPEN, max prices, order quantity limits
	if war.State != types.OpenState && war.State != types.HatchState {
		return nil, sdkerrors.Wrap(types.ErrInvalidStateForAction, war.State)
	} else if !war.ReserveDenomsEqualTo(msg.MaxPrices) {
		return nil, sdkerrors.Wrapf(types.ErrReserveDenomsMismatch, "%s do not match reserve; expected: %s", msg.MaxPrices.String(), strings.Join(war.ReserveTokens, ","))
	} else if war.AnyOrderQuantityLimitsExceeded(sdk.Coins{msg.Amount}) {
		return nil, sdkerrors.Wrap(types.ErrOrderQuantityLimitExceeded, msg.Amount.String())
	}

	// For the swapper, the first buy is the initialisation of the reserves
	// The max prices are used as the actual prices and one token is minted
	// The amount of token serves to define the price of adding more liquidity
	if war.CurrentSupply.IsZero() && war.FunctionType == types.SwapperFunction {
		return performFirstSwapperFunctionBuy(ctx, keeper, msg)
	}

	// Take max that buyer is willing to pay (enforces maxPrice <= balance)
	err := keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Buyer,
		types.BatchesIntermediaryAccount, msg.MaxPrices)
	if err != nil {
		return nil, err
	}

	// Create order
	order := types.NewBuyOrder(msg.Buyer, msg.Amount, msg.MaxPrices)

	// Get buy price and check if can add buy order to batch
	buyPrices, sellPrices, err := keeper.GetUpdatedBatchPricesAfterBuy(ctx, token, order)
	if err != nil {
		return nil, err
	}

	// Add buy order to batch
	keeper.AddBuyOrder(ctx, token, order, buyPrices, sellPrices)

	// Cancel unfulfillable orders
	keeper.CancelUnfulfillableOrders(ctx, token)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBuy,
			sdk.NewAttribute(types.AttributeKeyWar, msg.Amount.Denom),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyMaxPrices, msg.MaxPrices.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Buyer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func performFirstSwapperFunctionBuy(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBuy) (*sdk.Result, error) {

	// TODO: investigate effect that a high amount has on future buyers' ability to buy.

	token := msg.Amount.Denom
	war, found := keeper.GetWar(ctx, token)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrWarDoesNotExist, token)
	}

	// Check if initial liquidity violates sanity rate
	if war.ReservesViolateSanityRate(msg.MaxPrices) {
		return nil, sdkerrors.Wrap(types.ErrValuesViolateSanityRate, msg.MaxPrices.String())
	}

	// Use max prices as the amount to send to the liquidity pool (i.e. price)
	err := keeper.DepositReserve(ctx, war.Token, msg.Buyer, msg.MaxPrices)
	if err != nil {
		return nil, err
	}

	// Mint war tokens
	err = keeper.SupplyKeeper.MintCoins(ctx, types.WarsMintBurnAccount,
		sdk.Coins{msg.Amount})
	if err != nil {
		return nil, err
	}

	// Send war tokens to buyer
	err = keeper.SupplyKeeper.SendCoinsFromModuleToAccount(ctx,
		types.WarsMintBurnAccount, msg.Buyer, sdk.Coins{msg.Amount})
	if err != nil {
		return nil, err
	}

	// Update supply
	keeper.SetCurrentSupply(ctx, war.Token, war.CurrentSupply.Add(msg.Amount))

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeInitSwapper,
			sdk.NewAttribute(types.AttributeKeyWar, msg.Amount.Denom),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyChargedPrices, msg.MaxPrices.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Buyer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSell(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgSell) (*sdk.Result, error) {

	token := msg.Amount.Denom
	war, found := keeper.GetWar(ctx, token)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrWarDoesNotExist, token)
	}

	// Check sells allowed, current state is OPEN, and order limits not exceeded
	if !war.AllowSells {
		return nil, sdkerrors.Wrap(types.ErrWarDoesNotAllowSelling, token)
	} else if war.State != types.OpenState {
		return nil, sdkerrors.Wrap(types.ErrInvalidStateForAction, war.State)
	} else if war.AnyOrderQuantityLimitsExceeded(sdk.Coins{msg.Amount}) {
		return nil, sdkerrors.Wrap(types.ErrOrderQuantityLimitExceeded, msg.Amount.String())
	}

	// Send coins to be burned from seller (enforces sellAmount <= balance)
	err := keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Seller,
		types.WarsMintBurnAccount, sdk.Coins{msg.Amount})
	if err != nil {
		return nil, err
	}

	// Burn war tokens to be sold
	err = keeper.SupplyKeeper.BurnCoins(ctx, types.WarsMintBurnAccount,
		sdk.Coins{msg.Amount})
	if err != nil {
		return nil, err
	}

	// Create order
	order := types.NewSellOrder(msg.Seller, msg.Amount)

	// Get sell price and check if can add sell order to batch
	buyPrices, sellPrices, err := keeper.GetUpdatedBatchPricesAfterSell(ctx, token, order)
	if err != nil {
		return nil, err
	}

	// Add sell order to batch
	keeper.AddSellOrder(ctx, token, order, buyPrices, sellPrices)

	//// Cancel unfulfillable orders (Note: no need)
	//keeper.CancelUnfulfillableOrders(ctx, token)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSell,
			sdk.NewAttribute(types.AttributeKeyWar, msg.Amount.Denom),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Seller.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSwap(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgSwap) (*sdk.Result, error) {

	war, found := keeper.GetWar(ctx, msg.WarToken)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrWarDoesNotExist, msg.WarToken)
	}

	// Confirm that function type is swapper_function and state is OPEN
	if war.FunctionType != types.SwapperFunction {
		return nil, sdkerrors.Wrap(types.ErrFunctionNotAvailableForFunctionType, war.FunctionType)
	} else if war.State != types.OpenState {
		return nil, sdkerrors.Wrap(types.ErrInvalidStateForAction, war.State)
	}

	// Check that from and to use reserve token names
	fromAndTo := sdk.NewCoins(msg.From, sdk.NewCoin(msg.ToToken, sdk.OneInt()))
	fromAndToDenoms := msg.From.Denom + "," + msg.ToToken
	if !war.ReserveDenomsEqualTo(fromAndTo) {
		return nil, sdkerrors.Wrapf(types.ErrReserveDenomsMismatch, "%s do not match reserve; expected: %s", fromAndToDenoms, war.ReserveTokens)
	}

	// Check if order quantity limit exceeded
	if war.AnyOrderQuantityLimitsExceeded(sdk.Coins{msg.From}) {
		return nil, sdkerrors.Wrap(types.ErrOrderQuantityLimitExceeded, msg.From.String())
	}

	// Take coins to be swapped from swapper (enforces swapAmount <= balance)
	err := keeper.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Swapper,
		types.BatchesIntermediaryAccount, sdk.Coins{msg.From})
	if err != nil {
		return nil, err
	}

	// Create order
	order := types.NewSwapOrder(msg.Swapper, msg.From, msg.ToToken)

	// Add swap order to batch
	keeper.AddSwapOrder(ctx, msg.WarToken, order)

	//// Cancel unfulfillable orders (Note: no need)
	//keeper.CancelUnfulfillableOrders(ctx, token)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSwap,
			sdk.NewAttribute(types.AttributeKeyWar, msg.WarToken),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.From.Amount.String()),
			sdk.NewAttribute(types.AttributeKeySwapFromToken, msg.From.Denom),
			sdk.NewAttribute(types.AttributeKeySwapToToken, msg.ToToken),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Swapper.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgMakeOutcomePayment(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgMakeOutcomePayment) (*sdk.Result, error) {

	war, found := keeper.GetWar(ctx, msg.WarToken)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrWarDoesNotExist, msg.WarToken)
	}

	// Confirm that state is OPEN and that outcome payment is not nil
	if war.State != types.OpenState {
		return nil, sdkerrors.Wrap(types.ErrInvalidStateForAction, war.State)
	} else if war.OutcomePayment.Empty() {
		return nil, types.ErrCannotMakeZeroOutcomePayment
	}

	// Send outcome payment to reserve
	err := keeper.DepositReserve(ctx, war.Token, msg.Sender, war.OutcomePayment)
	if err != nil {
		return nil, err
	}

	// Set war state to SETTLE
	keeper.SetWarState(ctx, war.Token, types.SettleState)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMakeOutcomePayment,
			sdk.NewAttribute(types.AttributeKeyWar, msg.WarToken),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Sender.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgWithdrawShare(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgWithdrawShare) (*sdk.Result, error) {

	war, found := keeper.GetWar(ctx, msg.WarToken)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrWarDoesNotExist, msg.WarToken)
	}

	// Check that state is SETTLE
	if war.State != types.SettleState {
		return nil, sdkerrors.Wrap(types.ErrInvalidStateForAction, war.State)
	}

	// Get number of war tokens owned by the recipient
	warTokensOwnedAmount := keeper.BankKeeper.GetCoins(ctx, msg.Recipient).AmountOf(msg.WarToken)
	if warTokensOwnedAmount.IsZero() {
		return nil, sdkerrors.Wrap(types.ErrNoWarTokensOwned, warTokensOwnedAmount.String())
	}
	warTokensOwned := sdk.NewCoin(msg.WarToken, warTokensOwnedAmount)

	// Send coins to be burned from recipient
	err := keeper.SupplyKeeper.SendCoinsFromAccountToModule(
		ctx, msg.Recipient, types.WarsMintBurnAccount, sdk.NewCoins(warTokensOwned))
	if err != nil {
		return nil, err
	}

	// Burn war tokens
	err = keeper.SupplyKeeper.BurnCoins(ctx, types.WarsMintBurnAccount,
		sdk.NewCoins(sdk.NewCoin(msg.WarToken, warTokensOwnedAmount)))
	if err != nil {
		return nil, err
	}

	// Calculate amount owned
	remainingReserve := keeper.GetReserveBalances(ctx, war.Token)
	warTokensShare := warTokensOwnedAmount.ToDec().QuoInt(war.CurrentSupply.Amount)
	reserveOwedDec := sdk.NewDecCoinsFromCoins(remainingReserve...).MulDec(warTokensShare)
	reserveOwed, _ := reserveOwedDec.TruncateDecimal()

	// Send coins owed to recipient
	err = keeper.WithdrawReserve(ctx, war.Token, msg.Recipient, reserveOwed)
	if err != nil {
		return nil, err
	}

	// Update supply
	keeper.SetCurrentSupply(ctx, war.Token, war.CurrentSupply.Sub(warTokensOwned))

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawShare,
			sdk.NewAttribute(types.AttributeKeyWar, msg.WarToken),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Recipient.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, reserveOwed.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Recipient.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
