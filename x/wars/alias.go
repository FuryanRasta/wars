package wars

// nolint
// autogenerated code using github.com/haasted/alias-generator.
// based on functionality in github.com/rigelrozanski/multitool

import (
	"github.com/mage-war/wars/x/wars/client"
	"github.com/mage-war/wars/x/wars/internal/keeper"
	"github.com/mage-war/wars/x/wars/internal/types"
)

const (
	PowerFunction     = types.PowerFunction
	SigmoidFunction   = types.SigmoidFunction
	SwapperFunction   = types.SwapperFunction
	AugmentedFunction = types.AugmentedFunction

	HatchState  = types.HatchState
	OpenState   = types.OpenState
	SettleState = types.SettleState

	DoNotModifyField = types.DoNotModifyField

	AnyNumberOfReserveTokens = types.AnyNumberOfReserveTokens

	DefaultCodespace = types.DefaultCodespace

	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	StoreKey          = types.StoreKey

	WarsMintBurnAccount       = types.WarsMintBurnAccount
	BatchesIntermediaryAccount = types.BatchesIntermediaryAccount
	WarsReserveAccount        = types.WarsReserveAccount

	QuerierRoute = types.QuerierRoute
	RouterKey    = types.RouterKey
)

var (
	// functions aliases

	NewQuerier = keeper.NewQuerier
	NewKeeper  = keeper.NewKeeper

	RegisterInvariants = keeper.RegisterInvariants
	AllInvariants      = keeper.AllInvariants
	SupplyInvariant    = keeper.SupplyInvariant
	ReserveInvariant   = keeper.ReserveInvariant

	RegisterCodec = types.RegisterCodec

	NewBatch         = types.NewBatch
	NewBaseOrder     = types.NewBaseOrder
	NewBuyOrder      = types.NewBuyOrder
	NewSellOrder     = types.NewSellOrder
	NewSwapOrder     = types.NewSwapOrder
	NewFunctionParam = types.NewFunctionParam
	NewWar          = types.NewWar

	RoundReservePrice     = types.RoundReservePrice
	RoundReserveReturn    = types.RoundReserveReturn
	RoundFee              = types.RoundFee
	RoundReservePrices    = types.RoundReservePrices
	RoundReserveReturns   = types.RoundReserveReturns
	MultiplyDecCoinByInt  = types.MultiplyDecCoinByInt
	MultiplyDecCoinsByInt = types.MultiplyDecCoinsByInt
	MultiplyDecCoinByDec  = types.MultiplyDecCoinByDec
	MultiplyDecCoinsByDec = types.MultiplyDecCoinsByDec
	DivideDecCoinByDec    = types.DivideDecCoinByDec
	DivideDecCoinsByDec   = types.DivideDecCoinsByDec
	AdjustFees            = types.AdjustFees

	NewGenesisState     = types.NewGenesisState
	ValidateGenesis     = types.ValidateGenesis
	DefaultGenesisState = types.DefaultGenesisState

	GetWarKey      = types.GetWarKey
	GetBatchKey     = types.GetBatchKey
	GetLastBatchKey = types.GetLastBatchKey

	NewMsgCreateWar         = types.NewMsgCreateWar
	NewMsgEditWar           = types.NewMsgEditWar
	NewMsgBuy                = types.NewMsgBuy
	NewMsgSell               = types.NewMsgSell
	NewMsgSwap               = types.NewMsgSwap
	NewMsgMakeOutcomePayment = types.NewMsgMakeOutcomePayment
	NewMsgWithdrawShare      = types.NewMsgWithdrawShare

	ParseFunctionParams = client.ParseFunctionParams
	ParseSigners        = client.ParseSigners
	ParseTwoPartCoin    = client.ParseTwoPartCoin

	// variable aliases

	ModuleCdc = types.ModuleCdc

	RequiredParamsForFunctionType    = types.RequiredParamsForFunctionType
	NoOfReserveTokensForFunctionType = types.NoOfReserveTokensForFunctionType
	ExtraParameterRestrictions       = types.ExtraParameterRestrictions

	ErrArgumentMustBePositive               = types.ErrArgumentMustBePositive
	ErrArgumentMustBeInteger                = types.ErrArgumentMustBeInteger
	ErrArgumentMustBeBetween                = types.ErrArgumentMustBeBetween
	ErrArgumentCannotBeEmpty                = types.ErrArgumentCannotBeEmpty
	ErrArgumentCannotBeNegative             = types.ErrArgumentCannotBeNegative
	ErrArgumentMissingOrNonFloat            = types.ErrArgumentMissingOrNonFloat
	ErrWarDoesNotExist                     = types.ErrWarDoesNotExist
	ErrWarAlreadyExists                    = types.ErrWarAlreadyExists
	ErrWarTokenCannotBeStakingToken        = types.ErrWarTokenCannotBeStakingToken
	ErrInvalidStateForAction                = types.ErrInvalidStateForAction
	ErrReserveDenomsMismatch                = types.ErrReserveDenomsMismatch
	ErrOrderQuantityLimitExceeded           = types.ErrOrderQuantityLimitExceeded
	ErrValuesViolateSanityRate              = types.ErrValuesViolateSanityRate
	ErrWarDoesNotAllowSelling              = types.ErrWarDoesNotAllowSelling
	ErrFunctionNotAvailableForFunctionType  = types.ErrFunctionNotAvailableForFunctionType
	ErrCannotMakeZeroOutcomePayment         = types.ErrCannotMakeZeroOutcomePayment
	ErrNoWarTokensOwned                    = types.ErrNoWarTokensOwned
	ErrCannotBurnMoreThanSupply             = types.ErrCannotBurnMoreThanSupply
	ErrFeesCannotBeOrExceed100Percent       = types.ErrFeesCannotBeOrExceed100Percent
	ErrFromAndToCannotBeTheSameToken        = types.ErrFromAndToCannotBeTheSameToken
	ErrCannotMintMoreThanMaxSupply          = types.ErrCannotMintMoreThanMaxSupply
	ErrMaxPriceExceeded                     = types.ErrMaxPriceExceeded
	ErrInsufficientReserveToBuy             = types.ErrInsufficientReserveToBuy
	ErrIncorrectNumberOfFunctionParameters  = types.ErrIncorrectNumberOfFunctionParameters
	ErrFunctionParameterMissingOrNonFloat   = types.ErrFunctionParameterMissingOrNonFloat
	ErrFunctionRequiresNonZeroCurrentSupply = types.ErrFunctionRequiresNonZeroCurrentSupply
	ErrTokenIsNotAValidReserveToken         = types.ErrTokenIsNotAValidReserveToken
	ErrSwapAmountTooSmallToGiveAnyReturn    = types.ErrSwapAmountTooSmallToGiveAnyReturn
	ErrSwapAmountCausesReserveDepletion     = types.ErrSwapAmountCausesReserveDepletion
	ErrInvalidCoinDenomination              = types.ErrInvalidCoinDenomination
	ErrMaxSupplyDenomDoesNotMatchTokenDenom = types.ErrMaxSupplyDenomDoesNotMatchTokenDenom
	ErrDidNotEditAnything                   = types.ErrDidNotEditAnything
	ErrWarTokenCannotAlsoBeReserveToken    = types.ErrWarTokenCannotAlsoBeReserveToken
	ErrDuplicateReserveToken                = types.ErrDuplicateReserveToken
	ErrUnrecognizedFunctionType             = types.ErrUnrecognizedFunctionType
	ErrIncorrectNumberOfReserveTokens       = types.ErrIncorrectNumberOfReserveTokens
	ErrInvalidFunctionParameter             = types.ErrInvalidFunctionParameter
	ErrArgumentMissingOrNonUInteger         = types.ErrArgumentMissingOrNonUInteger
	ErrArgumentMissingOrNonBoolean          = types.ErrArgumentMissingOrNonBoolean

	WarsKeyPrefix       = types.WarsKeyPrefix
	BatchesKeyPrefix     = types.BatchesKeyPrefix
	LastBatchesKeyPrefix = types.LastBatchesKeyPrefix
)

type (
	Keeper = keeper.Keeper

	Batch     = types.Batch
	BaseOrder = types.BaseOrder
	BuyOrder  = types.BuyOrder
	SellOrder = types.SellOrder
	SwapOrder = types.SwapOrder

	FunctionParamRestrictions = types.FunctionParamRestrictions
	FunctionParam             = types.FunctionParam
	FunctionParams            = types.FunctionParams

	War = types.War

	GenesisState = types.GenesisState

	MsgCreateWar         = types.MsgCreateWar
	MsgEditWar           = types.MsgEditWar
	MsgBuy                = types.MsgBuy
	MsgSell               = types.MsgSell
	MsgSwap               = types.MsgSwap
	MsgMakeOutcomePayment = types.MsgMakeOutcomePayment
	MsgWithdrawShare      = types.MsgWithdrawShare
)
