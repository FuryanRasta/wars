package simulation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	defaultReserveTokens = []string{sdk.DefaultWarDenom}

	blankOrderQuantityLimits    = sdk.Coins{}
	blankOutcomePayment         = sdk.Coins{}
	blankSanityRate             = sdk.MustNewDecFromStr("0")
	blankSanityMarginPercentage = sdk.MustNewDecFromStr("0")

	tokenPrefix    = "token"
	totalWarCount = 0 // Updated for each war created
	maxWarCount   = 0 // Set during genesis creation
	gas            = uint64(100000000)

	swapperWars []string
)
