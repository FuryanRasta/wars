package keeper

// DONTCOVER

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/warmage-sports/wars/x/wars/internal/types"
)

// RegisterInvariants registers all supply invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "wars-supply",
		SupplyInvariant(k))
	ir.RegisterRoute(types.ModuleName, "wars-reserve",
		ReserveInvariant(k))
}

// AllInvariants runs all invariants of the wars module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		res, stop := SupplyInvariant(k)(ctx)
		if stop {
			return res, stop
		}
		return ReserveInvariant(k)(ctx)
	}
}

func SupplyInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var count int

		// Get supply of coins held in accounts (includes stake token)
		supplyInAccounts := sdk.Coins{}
		k.accountKeeper.IterateAccounts(ctx, func(acc exported.Account) bool {
			supplyInAccounts = supplyInAccounts.Add(acc.GetCoins()...)
			return false
		})

		iterator := k.GetWarIterator(ctx)
		for ; iterator.Valid(); iterator.Next() {
			war := k.MustGetWarByKey(ctx, iterator.Key())
			denom := war.Token
			batch := k.MustGetBatch(ctx, denom)

			// Add war current supply
			supplyInWarsAndBatches := war.CurrentSupply

			// Subtract amount to be burned (this amount was already burned
			// in handleMsgSell but is still a part of war's CurrentSupply)
			for _, s := range batch.Sells {
				if !s.Cancelled {
					supplyInWarsAndBatches = supplyInWarsAndBatches.Sub(
						s.Amount)
				}
			}

			// Check that amount matches supply in accounts
			inAccounts := supplyInAccounts.AmountOf(war.Token)
			if !supplyInWarsAndBatches.Amount.Equal(inAccounts) {
				count++
				msg += fmt.Sprintf("total %s supply invariance:\n"+
					"\ttotal %s supply: %s\n"+
					"\tsum of %s in accounts: %s\n",
					denom, denom, supplyInWarsAndBatches.Amount.String(),
					denom, inAccounts.String())
			}
		}

		broken := count != 0
		return sdk.FormatInvariant(types.ModuleName, "supply", fmt.Sprintf(
			"%d Wars supply invariants broken\n%s", count, msg)), broken
	}
}

func ReserveInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var count int

		iterator := k.GetWarIterator(ctx)
		for ; iterator.Valid(); iterator.Next() {
			war := k.MustGetWarByKey(ctx, iterator.Key())
			denom := war.Token

			if war.FunctionType == types.AugmentedFunction ||
				war.FunctionType == types.SwapperFunction {
				continue // Check does not apply to augmented/swapper functions
			}

			expectedReserve := war.ReserveAtSupply(war.CurrentSupply.Amount)
			expectedRounded := expectedReserve.Ceil().TruncateInt()
			actualReserve := k.GetReserveBalances(ctx, denom)

			for _, r := range actualReserve {
				if r.Amount.LT(expectedRounded) {
					count++
					msg += fmt.Sprintf("%s reserve invariance:\n"+
						"\texpected(ceil-rounded) %s reserve: %s\n"+
						"\tactual %s reserve: %s\n",
						denom, denom, expectedReserve.String(),
						denom, r.String())
				}
			}
		}

		broken := count != 0
		return sdk.FormatInvariant(types.ModuleName, "reserve", fmt.Sprintf(
			"%d Wars reserve invariants broken\n%s", count, msg)), broken
	}
}
