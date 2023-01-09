package types

type GenesisState struct {
	Wars   []War  `json:"wars" yaml:"wars"`
	Batches []Batch `json:"batches" yaml:"batches"`
	Params  Params  `json:"params" yaml:"params"`
}

func NewGenesisState(wars []War, batches []Batch, params Params) GenesisState {
	return GenesisState{
		Wars:   wars,
		Batches: batches,
		Params:  params,
	}
}

func ValidateGenesis(data GenesisState) error {
	err := ValidateParams(data.Params)
	if err != nil {
		return err
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Wars:   nil,
		Batches: nil,
		Params:  DefaultParams(),
	}
}
