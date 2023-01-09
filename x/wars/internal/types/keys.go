package types

const (
	// ModuleName is the name of this module
	ModuleName = "wars"

	// StoreKey is the default store key for this module
	StoreKey = ModuleName

	// DefaultParamspace is the default param space for this module
	DefaultParamspace = ModuleName

	// WarsMintBurnAccount the root string for the wars mint burn account address
	WarsMintBurnAccount = "wars_mint_burn_account"

	// BatchesIntermediaryAccount the root string for the batches account address
	BatchesIntermediaryAccount = "batches_intermediary_account"

	// WarsReserveAccount the root string for the wars reserve account address
	WarsReserveAccount = "wars_reserve_account"

	// QuerierRoute is the querier route for this module's store.
	QuerierRoute = ModuleName

	// RouterKey is the message route for this module
	RouterKey = ModuleName
)

// Wars and batches are stored as follow:
//
// - Wars: 0x00<war_token_bytes>
// - Batches: 0x01<war_token_bytes>
// - Last batches: 0x02<war_token_bytes>
var (
	WarsKeyPrefix       = []byte{0x00} // key for wars
	BatchesKeyPrefix     = []byte{0x01} // key for batches
	LastBatchesKeyPrefix = []byte{0x02} // key for last batches
)

func GetWarKey(token string) []byte {
	return append(WarsKeyPrefix, []byte(token)...)
}

func GetBatchKey(token string) []byte {
	return append(BatchesKeyPrefix, []byte(token)...)
}

func GetLastBatchKey(token string) []byte {
	return append(LastBatchesKeyPrefix, []byte(token)...)
}
