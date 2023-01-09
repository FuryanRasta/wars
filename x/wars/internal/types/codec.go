package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc is the codec for the module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&War{}, "wars/War", nil)
	cdc.RegisterConcrete(&FunctionParam{}, "wars/FunctionParam", nil)
	cdc.RegisterConcrete(&Batch{}, "wars/Batch", nil)
	cdc.RegisterConcrete(&BaseOrder{}, "wars/BaseOrder", nil)
	cdc.RegisterConcrete(&BuyOrder{}, "wars/BuyOrder", nil)
	cdc.RegisterConcrete(&SellOrder{}, "wars/SellOrder", nil)
	cdc.RegisterConcrete(&SwapOrder{}, "wars/SwapOrder", nil)
	cdc.RegisterConcrete(MsgCreateWar{}, "wars/MsgCreateWar", nil)
	cdc.RegisterConcrete(MsgEditWar{}, "wars/MsgEditWar", nil)
	cdc.RegisterConcrete(MsgBuy{}, "wars/MsgBuy", nil)
	cdc.RegisterConcrete(MsgSell{}, "wars/MsgSell", nil)
	cdc.RegisterConcrete(MsgSwap{}, "wars/MsgSwap", nil)
	cdc.RegisterConcrete(MsgMakeOutcomePayment{}, "wars/MsgMakeOutcomePayment", nil)
	cdc.RegisterConcrete(MsgWithdrawShare{}, "wars/MsgWithdrawShare", nil)
}
