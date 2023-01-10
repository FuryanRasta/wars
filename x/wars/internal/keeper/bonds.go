package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mage-war/wars/x/wars/internal/types"
)

func (k Keeper) GetWarIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.WarsKeyPrefix)
}

func (k Keeper) GetWar(ctx sdk.Context, token string) (war types.War, found bool) {
	store := ctx.KVStore(k.storeKey)
	if !k.WarExists(ctx, token) {
		return
	}
	bz := store.Get(types.GetWarKey(token))
	k.cdc.MustUnmarshalBinaryBare(bz, &war)
	return war, true
}

func (k Keeper) MustGetWar(ctx sdk.Context, token string) types.War {
	war, found := k.GetWar(ctx, token)
	if !found {
		panic(fmt.Sprintf("war '%s' not found\n", token))
	}
	return war
}

func (k Keeper) MustGetWarByKey(ctx sdk.Context, key []byte) types.War {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(key) {
		panic("war not found")
	}

	bz := store.Get(key)
	var war types.War
	k.cdc.MustUnmarshalBinaryBare(bz, &war)

	return war
}

func (k Keeper) WarExists(ctx sdk.Context, token string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetWarKey(token))
}

func (k Keeper) SetWar(ctx sdk.Context, token string, war types.War) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetWarKey(token), k.cdc.MustMarshalBinaryBare(war))
}

func (k Keeper) DepositReserve(ctx sdk.Context, token string, from sdk.AccAddress, amount sdk.Coins) error {
	// Send tokens to wars reserve account
	err := k.SupplyKeeper.SendCoinsFromAccountToModule(
		ctx, from, types.WarsReserveAccount, amount)
	if err != nil {
		return err
	}

	// Update war reserve
	k.setReserveBalances(ctx, token,
		k.MustGetWar(ctx, token).CurrentReserve.Add(amount...))
	return nil
}

func (k Keeper) DepositReserveFromModule(ctx sdk.Context, token string,
	fromModule string, amount sdk.Coins) error {

	// Send tokens to wars reserve account
	err := k.SupplyKeeper.SendCoinsFromModuleToModule(
		ctx, fromModule, types.WarsReserveAccount, amount)
	if err != nil {
		return err
	}

	// Update war reserve
	k.setReserveBalances(ctx, token,
		k.MustGetWar(ctx, token).CurrentReserve.Add(amount...))
	return nil
}

func (k Keeper) WithdrawReserve(ctx sdk.Context, token string,
	to sdk.AccAddress, amount sdk.Coins) error {

	// Send tokens from wars reserve account
	err := k.SupplyKeeper.SendCoinsFromModuleToAccount(
		ctx, types.WarsReserveAccount, to, amount)
	if err != nil {
		return err
	}

	// Update war reserve
	k.setReserveBalances(ctx, token,
		k.MustGetWar(ctx, token).CurrentReserve.Sub(amount))
	return nil
}

func (k Keeper) setReserveBalances(ctx sdk.Context, token string, balance sdk.Coins) {
	war := k.MustGetWar(ctx, token)
	war.CurrentReserve = balance
	k.SetWar(ctx, token, war)
}

func (k Keeper) GetReserveBalances(ctx sdk.Context, token string) sdk.Coins {
	return k.MustGetWar(ctx, token).CurrentReserve
}

func (k Keeper) GetSupplyAdjustedForBuy(ctx sdk.Context, token string) sdk.Coin {
	war := k.MustGetWar(ctx, token)
	batch := k.MustGetBatch(ctx, token)
	supply := war.CurrentSupply
	return supply.Add(batch.TotalBuyAmount)
}

func (k Keeper) GetSupplyAdjustedForSell(ctx sdk.Context, token string) sdk.Coin {
	war := k.MustGetWar(ctx, token)
	batch := k.MustGetBatch(ctx, token)
	supply := war.CurrentSupply
	return supply.Sub(batch.TotalSellAmount)
}

func (k Keeper) SetCurrentSupply(ctx sdk.Context, token string, currentSupply sdk.Coin) {
	if currentSupply.IsNegative() {
		panic("current supply cannot be negative")
	}
	war := k.MustGetWar(ctx, token)
	war.CurrentSupply = currentSupply
	k.SetWar(ctx, token, war)
}

func (k Keeper) SetWarState(ctx sdk.Context, token string, newState string) {
	war := k.MustGetWar(ctx, token)
	previousState := war.State
	war.State = newState
	k.SetWar(ctx, token, war)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("updated state for %s from %s to %s", war.Token, previousState, newState))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeStateChange,
		sdk.NewAttribute(types.AttributeKeyWar, war.Token),
		sdk.NewAttribute(types.AttributeKeyOldState, previousState),
		sdk.NewAttribute(types.AttributeKeyNewState, newState),
	))
}

func (k Keeper) ReservedWarToken(ctx sdk.Context, warToken string) bool {
	reservedWarTokens := k.GetParams(ctx).ReservedWarTokens
	for _, rbt := range reservedWarTokens {
		if warToken == rbt {
			return true
		}
	}
	return false
}
