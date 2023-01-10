package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"
	"github.com/mage-war/wars/x/wars/client"
	"github.com/mage-war/wars/x/wars/internal/types"
	"net/http"
	"strings"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/wars/create_war", createWarRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/wars/edit_war", editWarRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/wars/buy", buyRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/wars/sell", sellRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/wars/swap", swapRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/wars/make_outcome_payment", makeOutcomePaymentRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/wars/withdraw_share", withdrawShareRequestHandler(cliCtx)).Methods("POST")
}

type createWarReq struct {
	BaseReq                rest.BaseReq `json:"base_req" yaml:"base_req"`
	Token                  string       `json:"token" yaml:"token"`
	Name                   string       `json:"name" yaml:"name"`
	Description            string       `json:"description" yaml:"description"`
	FunctionType           string       `json:"function_type" yaml:"function_type"`
	FunctionParameters     string       `json:"function_parameters" yaml:"function_parameters"`
	ReserveTokens          string       `json:"reserve_tokens" yaml:"reserve_tokens"`
	TxFeePercentage        string       `json:"tx_fee_percentage" yaml:"tx_fee_percentage"`
	ExitFeePercentage      string       `json:"exit_fee_percentage" yaml:"exit_fee_percentage"`
	FeeAddress             string       `json:"fee_address" yaml:"fee_address"`
	MaxSupply              string       `json:"max_supply" yaml:"max_supply"`
	OrderQuantityLimits    string       `json:"order_quantity_limits" yaml:"order_quantity_limits"`
	SanityRate             string       `json:"sanity_rate" yaml:"sanity_rate"`
	SanityMarginPercentage string       `json:"sanity_margin_percentage" yaml:"sanity_margin_percentage"`
	AllowSells             string       `json:"allow_sells" yaml:"allow_sells"`
	Signers                string       `json:"signers" yaml:"signers"`
	BatchBlocks            string       `json:"batch_blocks" yaml:"batch_blocks"`
	OutcomePayment         string       `json:"outcome_payment" yaml:"outcome_payment"`
}

func createWarRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createWarReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		creator, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse function parameters
		functionParams, err := client.ParseFunctionParams(req.FunctionParameters)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse reserve tokens
		reserveTokens := strings.Split(req.ReserveTokens, ",")

		// Parse tx fee percentage
		txFeePercentageDec, err := sdk.NewDecFromStr(req.TxFeePercentage)
		if err != nil {
			err = sdkerrors.Wrap(types.ErrArgumentMissingOrNonFloat, "tx fee percentage")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse exit fee percentage
		exitFeePercentageDec, err := sdk.NewDecFromStr(req.ExitFeePercentage)
		if err != nil {
			err = sdkerrors.Wrap(types.ErrArgumentMissingOrNonFloat, "exit fee percentage")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse fee address
		feeAddress, err2 := sdk.AccAddressFromBech32(req.FeeAddress)
		if err2 != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err2.Error())
			return
		}

		// Parse max supply
		maxSupply, err2 := sdk.ParseCoin(req.MaxSupply)
		if err2 != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err2.Error())
			return
		}

		// Parse order quantity limits
		orderQuantityLimits, err2 := sdk.ParseCoins(req.OrderQuantityLimits)
		if err2 != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err2.Error())
			return
		}

		// Parse sanity rate
		sanityRate, err := sdk.NewDecFromStr(req.SanityRate)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse sanity margin percentage
		sanityMarginPercentage, err := sdk.NewDecFromStr(req.SanityMarginPercentage)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse allowSells
		var allowSells bool
		allowSellsStrLower := strings.ToLower(req.AllowSells)
		if allowSellsStrLower == "true" {
			allowSells = true
		} else if allowSellsStrLower == "false" {
			allowSells = false
		} else {
			err := sdkerrors.Wrap(types.ErrArgumentMissingOrNonBoolean, "allow_sells")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse signers
		signers, err := client.ParseSigners(req.Signers)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse batch blocks
		batchBlocks, err2 := sdk.ParseUint(req.BatchBlocks)
		if err2 != nil {
			err := sdkerrors.Wrap(types.ErrArgumentMissingOrNonUInteger, "max batch blocks")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse outcome payment
		outcomePayment, err2 := sdk.ParseCoins(req.OutcomePayment)
		if err2 != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err2.Error())
			return
		}

		msg := types.NewMsgCreateWar(req.Token, req.Name, req.Description,
			creator, req.FunctionType, functionParams, reserveTokens,
			txFeePercentageDec, exitFeePercentageDec, feeAddress, maxSupply,
			orderQuantityLimits, sanityRate, sanityMarginPercentage,
			allowSells, signers, batchBlocks, outcomePayment)

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type editWarRequestReq struct {
	BaseReq                rest.BaseReq `json:"base_req" yaml:"base_req"`
	Token                  string       `json:"token" yaml:"token"`
	Name                   string       `json:"name" yaml:"name"`
	Description            string       `json:"description" yaml:"description"`
	OrderQuantityLimits    string       `json:"order_quantity_limits" yaml:"order_quantity_limits"`
	SanityRate             string       `json:"sanity_rate" yaml:"sanity_rate"`
	SanityMarginPercentage string       `json:"sanity_margin_percentage" yaml:"sanity_margin_percentage"`
	Signers                string       `json:"signers" yaml:"signers"`
}

func editWarRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req editWarRequestReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		editor, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse signers
		signers, err := client.ParseSigners(req.Signers)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgEditWar(req.Token, req.Name, req.Description,
			req.OrderQuantityLimits, req.SanityRate, req.SanityMarginPercentage,
			editor, signers)

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type buyReq struct {
	BaseReq    rest.BaseReq `json:"base_req" yaml:"base_req"`
	WarToken  string       `json:"war_token" yaml:"war_token"`
	WarAmount string       `json:"war_amount" yaml:"war_amount"`
	MaxPrices  string       `json:"max_prices" yaml:"max_prices"`
}

func buyRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req buyReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		buyer, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		warCoin, err := client.ParseTwoPartCoin(req.WarAmount, req.WarToken)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		maxPrices, err := sdk.ParseCoins(req.MaxPrices)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgBuy(buyer, warCoin, maxPrices)
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type sellReq struct {
	BaseReq    rest.BaseReq `json:"base_req" yaml:"base_req"`
	WarToken  string       `json:"war_token" yaml:"war_token"`
	WarAmount string       `json:"war_amount" yaml:"war_amount"`
}

func sellRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req sellReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		seller, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		warCoin, err := client.ParseTwoPartCoin(req.WarAmount, req.WarToken)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgSell(seller, warCoin)
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type swapReq struct {
	BaseReq    rest.BaseReq `json:"base_req" yaml:"base_req"`
	WarToken  string       `json:"war_token" yaml:"war_token"`
	FromAmount string       `json:"from_amount" yaml:"from_amount"`
	FromToken  string       `json:"from_token" yaml:"from_token"`
	ToToken    string       `json:"to_token" yaml:"to_token"`
}

func swapRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req swapReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		swapper, err := sdk.AccAddressFromBech32(baseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Check that from amount and token can be parsed to a coin
		fromCoin, err := client.ParseTwoPartCoin(req.FromAmount, req.FromToken)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgSwap(swapper, req.WarToken, fromCoin, req.ToToken)
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type makeOutcomePaymentReq struct {
	BaseReq   rest.BaseReq `json:"base_req" yaml:"base_req"`
	WarToken string       `json:"war_token" yaml:"war_token"`
}

func makeOutcomePaymentRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req makeOutcomePaymentReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		sender, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgMakeOutcomePayment(sender, req.WarToken)
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type withdrawShareReq struct {
	BaseReq   rest.BaseReq `json:"base_req" yaml:"base_req"`
	WarToken string       `json:"war_token" yaml:"war_token"`
}

func withdrawShareRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req withdrawShareReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		recipient, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgWithdrawShare(recipient, req.WarToken)
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
