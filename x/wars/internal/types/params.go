package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Parameter store keys
var (
	KeyReservedWarTokens = []byte("ReservedWarTokens")
)

// wars parameters
type Params struct {
	ReservedWarTokens []string `json:"reserved_war_tokens" yaml:"reserved_war_tokens"`
}

// ParamTable for wars module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(reservedWarTokens []string) Params {
	return Params{
		ReservedWarTokens: reservedWarTokens,
	}

}

// default wars module parameters
func DefaultParams() Params {
	return Params{
		ReservedWarTokens: []string{}, // no reserved war tokens
	}
}

// validate params
func ValidateParams(params Params) error {
	return nil
}

func (p Params) String() string {
	return fmt.Sprintf(`Wars Params:
  Reserved War Tokens: %s
`,
		p.ReservedWarTokens)
}

func validateReservedWarTokens(i interface{}) error {
	_, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyReservedWarTokens, &p.ReservedWarTokens, validateReservedWarTokens},
	}
}
