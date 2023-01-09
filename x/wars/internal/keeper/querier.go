package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/warmage-sports/wars/x/wars/client"
	"github.com/warmage-sports/wars/x/wars/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryWars          = "wars"
	QueryWar           = "war"
	QueryBatch          = "batch"
	QueryLastBatch      = "last_batch"
	QueryCurrentPrice   = "current_price"
	QueryCurrentReserve = "current_reserve"
	QueryCustomPrice    = "custom_price"
	QueryBuyPrice       = "buy_price"
	QuerySellReturn     = "sell_return"
	QuerySwapReturn     = "swap_return"
	QueryParams         = "params"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryWars:
			return queryWars(ctx, keeper)
		case QueryWar:
			return queryWar(ctx, path[1:], keeper)
		case QueryBatch:
			return queryBatch(ctx, path[1:], keeper)
		case QueryLastBatch:
			return queryLastBatch(ctx, path[1:], keeper)
		case QueryCurrentPrice:
			return queryCurrentPrice(ctx, path[1:], keeper)
		case QueryCurrentReserve:
			return queryCurrentReserve(ctx, path[1:], keeper)
		case QueryCustomPrice:
			return queryCustomPrice(ctx, path[1:], keeper)
		case QueryBuyPrice:
			return queryBuyPrice(ctx, path[1:], keeper)
		case QuerySellReturn:
			return querySellReturn(ctx, path[1:], keeper)
		case QuerySwapReturn:
			return querySwapReturn(ctx, path[1:], keeper)
		case QueryParams:
			return queryParams(ctx, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown wars query endpoint")
		}
	}
}

func zeroReserveTokensIfEmpty(reserveCoins sdk.Coins, war types.War) sdk.Coins {
	if reserveCoins.IsZero() {
		zeroes, _ := war.GetNewReserveDecCoins(sdk.OneDec()).TruncateDecimal()
		for i := range zeroes {
			zeroes[i].Amount = sdk.ZeroInt()
		}
		reserveCoins = zeroes
	}
	return reserveCoins
}

func zeroReserveTokensIfEmptyDec(reserveCoins sdk.DecCoins, war types.War) sdk.DecCoins {
	if reserveCoins.IsZero() {
		zeroes := war.GetNewReserveDecCoins(sdk.OneDec())
		for i := range zeroes {
			zeroes[i].Amount = sdk.ZeroDec()
		}
		reserveCoins = zeroes
	}
	return reserveCoins
}

func queryWars(ctx sdk.Context, keeper Keeper) (res []byte, err error) {
	var warsList types.QueryWars
	iterator := keeper.GetWarIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		var war types.War
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &war)
		warsList = append(warsList, war.Token)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, warsList)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func queryWar(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err error) {
	warToken := path[0]

	war, found := keeper.GetWar(ctx, warToken)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "war '%s' does not exist", warToken)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, war)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func queryBatch(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err error) {
	warToken := path[0]

	if !keeper.BatchExists(ctx, warToken) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "batch for '%s' does not exist", warToken)
	}

	batch := keeper.MustGetBatch(ctx, warToken)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, batch)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func queryLastBatch(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err error) {
	warToken := path[0]

	if !keeper.LastBatchExists(ctx, warToken) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "last batch for '%s' does not exist", warToken)
	}

	batch := keeper.MustGetLastBatch(ctx, warToken)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, batch)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func queryCurrentPrice(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err error) {
	warToken := path[0]

	war, found := keeper.GetWar(ctx, warToken)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "war '%s' does not exist", warToken)
	}

	reserveBalances := keeper.GetReserveBalances(ctx, warToken)
	reservePrices, err := war.GetCurrentPricesPT(reserveBalances)
	if err != nil {
		return nil, err
	}
	reservePrices = zeroReserveTokensIfEmptyDec(reservePrices, war)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, reservePrices)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func queryCurrentReserve(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err error) {
	warToken := path[0]

	war, found := keeper.GetWar(ctx, warToken)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "war '%s' does not exist", warToken)
	}

	reserveBalances := zeroReserveTokensIfEmpty(war.CurrentReserve, war)
	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, reserveBalances)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func queryCustomPrice(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err error) {
	warToken := path[0]
	warAmount := path[1]

	war, found := keeper.GetWar(ctx, warToken)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "war '%s' does not exist", warToken)
	}

	warCoin, err2 := client.ParseTwoPartCoin(warAmount, war.Token)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, err2.Error())
	}

	reservePrices, err := war.GetPricesAtSupply(warCoin.Amount)
	if err != nil {
		return nil, err
	}
	reservePrices = zeroReserveTokensIfEmptyDec(reservePrices, war)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, reservePrices)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func queryBuyPrice(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err error) {
	warToken := path[0]
	warAmount := path[1]

	war, found := keeper.GetWar(ctx, warToken)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("war '%s' does not exist", warToken))
	}

	warCoin, err2 := client.ParseTwoPartCoin(warAmount, warToken)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, err2.Error())
	}

	// Max supply cannot be less than supply (max supply >= supply)
	adjustedSupply := keeper.GetSupplyAdjustedForBuy(ctx, warToken)
	if war.MaxSupply.IsLT(adjustedSupply.Add(warCoin)) {
		return nil, sdkerrors.Wrap(types.ErrCannotMintMoreThanMaxSupply, war.MaxSupply.String())
	}

	reserveBalances := keeper.GetReserveBalances(ctx, warToken)
	reservePrices, err := war.GetPricesToMint(warCoin.Amount, reserveBalances)
	if err != nil {
		return nil, err
	}
	reservePricesRounded := types.RoundReservePrices(reservePrices)
	txFee := war.GetTxFees(reservePrices)

	var result types.QueryBuyPrice
	result.AdjustedSupply = adjustedSupply
	result.Prices = zeroReserveTokensIfEmpty(reservePricesRounded, war)
	result.TxFees = zeroReserveTokensIfEmpty(txFee, war)
	result.TotalPrices = zeroReserveTokensIfEmpty(reservePricesRounded.Add(txFee...), war)
	result.TotalFees = zeroReserveTokensIfEmpty(txFee, war)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, result)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func querySellReturn(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err error) {
	warToken := path[0]
	warAmount := path[1]

	war, found := keeper.GetWar(ctx, warToken)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "war '%s' does not exist", warToken)
	}

	warCoin, err2 := client.ParseTwoPartCoin(warAmount, warToken)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, err2.Error())
	}

	if !war.AllowSells {
		return nil, sdkerrors.Wrap(types.ErrWarDoesNotAllowSelling, war.Name)
	}

	// Cannot burn more tokens than what exists
	adjustedSupply := keeper.GetSupplyAdjustedForSell(ctx, warToken)
	if adjustedSupply.IsLT(warCoin) {
		return nil, sdkerrors.Wrap(types.ErrCannotBurnMoreThanSupply, adjustedSupply.String())
	}

	reserveBalances := keeper.GetReserveBalances(ctx, warToken)
	reserveReturns := war.GetReturnsForBurn(warCoin.Amount, reserveBalances)
	reserveReturnsRounded := types.RoundReserveReturns(reserveReturns)

	txFees := war.GetTxFees(reserveReturns)
	exitFees := war.GetExitFees(reserveReturns)
	totalFees := types.AdjustFees(txFees.Add(exitFees...), reserveReturnsRounded)

	var result types.QuerySellReturn
	result.AdjustedSupply = adjustedSupply
	result.Returns = zeroReserveTokensIfEmpty(reserveReturnsRounded, war)
	result.TxFees = zeroReserveTokensIfEmpty(txFees, war)
	result.ExitFees = zeroReserveTokensIfEmpty(exitFees, war)
	result.TotalReturns = zeroReserveTokensIfEmpty(reserveReturnsRounded.Sub(totalFees), war)
	result.TotalFees = zeroReserveTokensIfEmpty(totalFees, war)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, result)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func querySwapReturn(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err error) {
	warToken := path[0]
	fromToken := path[1]
	fromAmount := path[2]
	toToken := path[3]

	fromCoin, err2 := client.ParseTwoPartCoin(fromAmount, fromToken)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, err2.Error())
	}

	war, found := keeper.GetWar(ctx, warToken)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrWarDoesNotExist, warToken)
	}

	reserveBalances := keeper.GetReserveBalances(ctx, warToken)
	reserveReturns, txFee, err := war.GetReturnsForSwap(fromCoin, toToken, reserveBalances)
	if err != nil {
		return nil, err
	}

	if reserveReturns.Empty() {
		reserveReturns = sdk.Coins{sdk.Coin{Denom: toToken, Amount: sdk.ZeroInt()}}
	}

	var result types.QuerySwapReturn
	result.TotalFees = sdk.Coins{txFee}
	result.TotalReturns = reserveReturns

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, result)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(k.cdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
