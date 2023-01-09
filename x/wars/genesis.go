package wars

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/warmage-sports/wars/x/wars/internal/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	// Initialise wars
	for _, b := range data.Wars {
		keeper.SetWar(ctx, b.Token, b)
	}

	// Initialise batches
	for _, b := range data.Batches {
		keeper.SetBatch(ctx, b.Token, b)
	}

	// Initialise params
	keeper.SetParams(ctx, data.Params)
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	// Export wars and batches
	var wars []types.War
	var batches []types.Batch
	iterator := k.GetWarIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		war := k.MustGetWarByKey(ctx, iterator.Key())
		batch := k.MustGetBatch(ctx, war.Token)
		wars = append(wars, war)
		batches = append(batches, batch)
	}

	// Export params
	params := k.GetParams(ctx)

	return GenesisState{
		Wars:   wars,
		Batches: batches,
		Params:  params,
	}
}
