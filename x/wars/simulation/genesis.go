package simulation

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/warmage-sports/wars/x/wars/internal/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"math/rand"
)

// Simulation parameters constants
const (
	InitialWars            = "initial_wars"
	MaxWars                = "max_wars"
	MaxNumberOfInitialWars = 100
	MaxNumberOfWars        = 100000
)

// GenInitialNumberOfWars randomized initial number of wars
func GenInitialNumberOfWars(r *rand.Rand) (initialWars uint64) {
	return uint64(r.Int63n(MaxNumberOfInitialWars) + 1)
}

// GenMaxNumberOfWars randomized max number of wars
func GenMaxNumberOfWars(r *rand.Rand) (maxWars uint64) {
	return uint64(r.Int63n(MaxNumberOfWars-MaxNumberOfInitialWars) + MaxNumberOfInitialWars + 1)
}

// RandomizedGenState generates a random GenesisState
func RandomizedGenState(simState *module.SimulationState) {
	r := simState.Rand

	// Generate a random number of initial wars and maximum wars
	var initialWars, maxWars uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, InitialWars, &initialWars, simState.Rand,
		func(r *rand.Rand) { initialWars = GenInitialNumberOfWars(r) },
	)
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MaxWars, &maxWars, simState.Rand,
		func(r *rand.Rand) { maxWars = GenMaxNumberOfWars(r) },
	)

	if initialWars > maxWars {
		panic("initialWars > maxWars")
	}
	maxWarCount = int(maxWars)

	var wars []types.War
	var batches []types.Batch
	for i := 0; i < int(initialWars); i++ {
		simAccount, _ := simulation.RandomAcc(r, simState.Accounts)
		address := simAccount.Address

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
				initialWars -= 1 // Ignore this iteration
				continue
			}
			reserveTokens = []string{reserveToken1, reserveToken2}
		default:
			reserveTokens = defaultReserveTokens
		}
		functionParameters := getRandomFunctionParameters(r, functionType, true)

		// Max fee is 100, so exit fee uses 100-txFee as max
		txFeePercentage := simulation.RandomDecAmount(r, sdk.NewDec(100))
		exitFeePercentage := simulation.RandomDecAmount(r, sdk.NewDec(100).Sub(txFeePercentage))

		// Addresses
		feeAddress := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

		// Max supply, allow sells, batch blocks
		maxSupply := sdk.NewCoin(token, sdk.NewInt(int64(
			simulation.RandIntBetween(r, 1000000, 1000000000))))
		allowSells := getRandomAllowSellsValue(r)
		batchBlocks := sdk.NewUint(uint64(
			simulation.RandIntBetween(r, 1, 10)))
		outcomePayment := sdk.Coins(nil)
		state := getInitialWarState(functionType)

		war := types.NewWar(token, name, desc, creator, functionType,
			functionParameters, reserveTokens, txFeePercentage,
			exitFeePercentage, feeAddress, maxSupply, blankOrderQuantityLimits,
			blankSanityRate, blankSanityMarginPercentage, allowSells, signers,
			batchBlocks, outcomePayment, state)
		batch := types.NewBatch(war.Token, war.BatchBlocks)

		wars = append(wars, war)
		batches = append(batches, batch)
		incrementWarCount()
		if war.FunctionType == types.SwapperFunction {
			newSwapperWar(war.Token)
		}
	}

	warsGenesis := types.NewGenesisState(wars, batches,
		types.Params{ReservedWarTokens: defaultReserveTokens})

	fmt.Printf("Selected randomly generated wars genesis state:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, warsGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(warsGenesis)
}
