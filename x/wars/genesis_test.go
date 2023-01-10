package wars_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mage-war/wars/x/wars"
	"github.com/mage-war/wars/x/wars/internal/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestInitAndExportGenesis(t *testing.T) {
	app, ctx := createTestApp(false)
	genesisState := wars.DefaultGenesisState()
	require.Equal(t, 0, len(genesisState.Wars))
	require.Equal(t, 0, len(genesisState.Batches))

	token := "testtoken"
	name := "test token"
	description := "this is a test token"
	creator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	functionType := types.PowerFunction
	functionParameters := types.FunctionParams{
		types.NewFunctionParam("m", sdk.NewDec(12)),
		types.NewFunctionParam("n", sdk.NewDec(2)),
		types.NewFunctionParam("c", sdk.NewDec(100))}
	reserveTokens := []string{"reservetoken"}
	txFeePercentage := sdk.MustNewDecFromStr("0.1")
	exitFeePercentage := sdk.MustNewDecFromStr("0.2")
	feeAddress := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	maxSupply := sdk.NewInt64Coin(token, 10000)
	orderQuantityLimits := sdk.NewCoins(
		sdk.NewInt64Coin("token1", 1),
		sdk.NewInt64Coin("token2", 2),
		sdk.NewInt64Coin("token3", 3),
	)
	sanityRate := sdk.MustNewDecFromStr("0.3")
	sanityMarginPercentage := sdk.MustNewDecFromStr("0.4")
	allowSell := true
	signers := []sdk.AccAddress{creator}
	batchBlocks := sdk.NewUint(10)
	outcomePayment := sdk.NewCoins(
		sdk.NewInt64Coin("token1", 1),
		sdk.NewInt64Coin("token2", 2),
		sdk.NewInt64Coin("token3", 3),
	)
	state := "dummy_state"

	war := types.NewWar(token, name, description, creator, functionType,
		functionParameters, reserveTokens, txFeePercentage, exitFeePercentage,
		feeAddress, maxSupply, orderQuantityLimits, sanityRate, sanityMarginPercentage,
		allowSell, signers, batchBlocks, outcomePayment, state)
	batch := types.NewBatch(war.Token, war.BatchBlocks)

	genesisState = wars.NewGenesisState(
		[]types.War{war}, []types.Batch{batch}, types.DefaultParams())

	wars.InitGenesis(ctx, app.WarsKeeper, genesisState)

	returnedWar := app.WarsKeeper.MustGetWar(ctx, token)
	require.EqualValues(t, war, returnedWar)

	returnedBatch := app.WarsKeeper.MustGetBatch(ctx, token)
	require.Equal(t, batch, returnedBatch)

	exportedGenesisState := wars.ExportGenesis(ctx, app.WarsKeeper)
	require.Equal(t, genesisState.Wars, exportedGenesisState.Wars)
	require.Equal(t, genesisState.Batches, exportedGenesisState.Batches)
}
