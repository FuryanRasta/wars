package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

// REST variable names
//noinspection GoNameStartsWithPackageName
const (
	RestWarToken           = "war_token"
	RestWarAmount          = "war_amount"
	RestFromTokenWithAmount = "from_token_with_amount"
	RestToToken             = "to_token"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, queryRoute string) {
	registerQueryRoutes(cliCtx, r, queryRoute)
	registerTxRoutes(cliCtx, r)
}
