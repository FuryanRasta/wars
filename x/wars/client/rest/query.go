package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/mage-war/wars/x/wars/internal/types"
	"net/http"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	r.HandleFunc(
		"/wars", queryWarsHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/wars/{%s}", RestWarToken),
		queryWarHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/wars/{%s}/batch", RestWarToken),
		queryBatchHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/wars/{%s}/last_batch", RestWarToken),
		queryLastBatchHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/wars/{%s}/current_price", RestWarToken),
		queryCurrentPriceHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/wars/{%s}/current_reserve", RestWarToken),
		queryCurrentReserveHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/wars/{%s}/price/{%s}", RestWarToken, RestWarAmount),
		queryCustomPriceHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/wars/{%s}/buy_price/{%s}", RestWarToken, RestWarAmount),
		queryBuyPriceHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/wars/{%s}/sell_return/{%s}", RestWarToken, RestWarAmount),
		querySellReturnHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/wars/{%s}/swap_return/{%s}/{%s}", RestWarToken, RestFromTokenWithAmount, RestToToken),
		querySwapReturnHandler(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		"/wars/params",
		queryParamsRequestHandler(cliCtx),
	).Methods("GET")
}

func queryWarsHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/wars", queryRoute), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryWarHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		warToken := vars[RestWarToken]

		res, _, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/war/%s",
				queryRoute, warToken), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryBatchHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		warToken := vars[RestWarToken]

		res, _, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/batch/%s",
				queryRoute, warToken), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryLastBatchHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		warToken := vars[RestWarToken]

		res, _, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/last_batch/%s",
				queryRoute, warToken), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryCurrentPriceHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		warToken := vars[RestWarToken]

		res, _, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/current_price/%s",
				queryRoute, warToken), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryCurrentReserveHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		warToken := vars[RestWarToken]

		res, _, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/current_reserve/%s",
				queryRoute, warToken), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryCustomPriceHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		warToken := vars[RestWarToken]
		warAmount := vars[RestWarAmount]

		res, _, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/custom_price/%s/%s",
				queryRoute, warToken, warAmount), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryBuyPriceHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		warToken := vars[RestWarToken]
		warAmount := vars[RestWarAmount]

		res, _, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/buy_price/%s/%s",
				queryRoute, warToken, warAmount), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func querySellReturnHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		warToken := vars[RestWarToken]
		warAmount := vars[RestWarAmount]

		res, _, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/sell_return/%s/%s",
				queryRoute, warToken, warAmount), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func querySwapReturnHandler(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		warToken := vars[RestWarToken]
		fromTokenWithAmount := vars[RestFromTokenWithAmount]
		toToken := vars[RestToToken]

		reserveCoinWithAmount, err := sdk.ParseCoin(fromTokenWithAmount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/swap_return/%s/%s/%s/%s",
				queryRoute, warToken, reserveCoinWithAmount.Denom,
				reserveCoinWithAmount.Amount.String(), toToken), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryParamsRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/parameters", types.QuerierRoute)

		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
