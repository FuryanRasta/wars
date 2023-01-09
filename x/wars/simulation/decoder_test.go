package simulation

import (
	"fmt"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"

	"github.com/stretchr/testify/require"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/warmage-sports/wars/x/wars/internal/types"
)

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()
	sdk.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	return
}

func TestDecodeStore(t *testing.T) {
	cdc := makeTestCodec()

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
	lastBatch := types.NewBatch(war.Token, war.BatchBlocks)

	kvPairs := tmkv.Pairs{
		tmkv.Pair{Key: types.GetWarKey(token),
			Value: cdc.MustMarshalBinaryBare(war)},
		tmkv.Pair{Key: types.GetBatchKey(token),
			Value: cdc.MustMarshalBinaryBare(batch)},
		tmkv.Pair{Key: types.GetLastBatchKey(token),
			Value: cdc.MustMarshalBinaryBare(lastBatch)},
		tmkv.Pair{Key: []byte{0x99}, Value: []byte{0x99}},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"wars", fmt.Sprintf("%v\n%v", war, war)},
		{"batches", fmt.Sprintf("%v\n%v", batch, batch)},
		{"lastBatches", fmt.Sprintf("%v\n%v", lastBatch, lastBatch)},
		{"other", ""},
	}

	for i, tt := range tests {
		tt, i := tt, i
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() {
					DecodeStore(cdc, kvPairs[i], kvPairs[i])
				}, tt.name)
			default:
				require.Equal(t, tt.expectedLog,
					DecodeStore(cdc, kvPairs[i], kvPairs[i]), tt.name)
			}
		})
	}
}
