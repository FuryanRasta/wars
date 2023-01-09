# Messages

In this section we describe the processing of the wars messages and the corresponding updates to the state. All created/modified state objects specified by each message are defined within the [state](./02_state.md) section.

## MsgCreateWar

Wars can be created by any address using `MsgCreateWar`.

| **Field**              | **Type**           | **Description** |
|:-----------------------|:-------------------|:----------------|
| Token                  | `string`           | The denomination of the war's tokens (e.g. `abc`, `mytoken1`)
| Name                   | `string`           | A friendly name as a title for the war (e.g. `A B C`, `My Token`)
| Description            | `string`           | A description of what the war represents or its purpose
| FunctionType           | `string`           | The type of function that will define the waring curve (`power_function`, `sigmoid_function`, or `swapper_function`)
| FunctionParameters     | `FunctionParams`   | The parameters of the function defining the waring curve (e.g. `m:12,n:2,c:100`)
| Creator                | `sdk.AccAddress`   | The address of the account creating the war
| ReserveTokens          | `[]string`         | The token denominations that will be used as reserve (e.g. `res,rez`)
| TxFeePercentage        | `sdk.Dec`          | The percentage fee charged for buys/sells/swaps (e.g. `0.3`)
| ExitFeePercentage      | `sdk.Dec`          | The percentage fee charged for sells on top of the tx fee (e.g. `0.2`)
| FeeAddress             | `sdk.AccAddress`   | The address of the account that will store charged fees
| MaxSupply              | `sdk.Coin`         | The maximum number of war tokens that can be minted
| OrderQuantityLimits    | `sdk.Coins`        | The maximum number of tokens that one can buy/sell/swap in a single order (e.g. `100abc,200res,300rez`)
| SanityRate             | `sdk.Dec`          | For a swapper, restricts conversion rate (`r1/r2`) to `sanity rate ± sanity margin percentage`. `0` for no sanity checks.
| SanityMarginPercentage | `sdk.Dec`          | Used as described above. `0` for no sanity checks
| AllowSells             | `bool`             | Whether or not selling is allowed
| Signers                | `[]sdk.AccAddress` | The addresses of the accounts that must sign this message and any future message that edits the war's parameters.
| BatchBlocks            | `sdk.Uint`         | The lifespan of each orders batch in blocks
| OutcomePayment         | `sdk.Coins`        | The payment required to be made in order to transition a war from OPEN to SETTLE

```go
type MsgCreateWar struct {
	Token                  string
	Name                   string
	Description            string
	FunctionType           string
	FunctionParameters     FunctionParams
	Creator                sdk.AccAddress
	ReserveTokens          []string
	TxFeePercentage        sdk.Dec
	ExitFeePercentage      sdk.Dec
	FeeAddress             sdk.AccAddress
	MaxSupply              sdk.Coin
	OrderQuantityLimits    sdk.Coins
	SanityRate             sdk.Dec
	SanityMarginPercentage sdk.Dec
	AllowSells             bool
	Signers                []sdk.AccAddress
	BatchBlocks            sdk.Uint
	OutcomePayment         sdk.Coins
}
```

This message is expected to fail if:
- another war with this token is already registered, the token is the staking token, or the token is not a valid denomination
- name or description is an empty string
- function type is not one of the defined function types (`power_function`, `sigmoid_function`, `swapper_function`, `augmented_function`)
- function parameters are negative or invalid for the selected function type:
  - Valid example for `power_function`: `"m:12.5,n:2,c:100.12"` \
    (i.e. `m=12`, `n=2`, `n=100.12`)
  - Valid example for `sigmoid_function`: `"a:3.5,b:5.4,c:1.3"` \
    (i.e. `a=3.5`, `b=5.4`, `c=1.3`)
  - Valid example for `augmented_function`: `"d0:500.0,p0:0.01,theta:0.4,kappa:3.0"` \
    (i.e. `d0=500.0`, `p0=0.01`, `theta=0.4`, `kappa=3.0`)
  - For `swapper_function`: `""` (no parameters)
- function parameters do not satisfy the extra parameter restrictions
  - `power_function`: `n` must be an integer
  - `sigmoid_function`: `c != 0`
  - `augmented_function`:
    - `d0 != 0` and must be an integer
    - `p0 != 0`
    - `0 <= theta < 1`
    - `kappa != 0` and must be an integer
- reserve tokens list is invalid. Valid inputs are:
  - For `swapper_function`: two valid comma-separated denominations, e.g. `res,rez`
  - Otherwise: one or more valid comma-separated denominations, e.g. `res,rez,rex`
- tx or exit fee percentage is negative
- sum of tx and exit fee percentages exceeds 100%
- order quantity limits is not one or more valid comma-separated amount
  - Valid example: `"100res,200rez"`
- max supply value is not in the war token denomination
- sanity rate is neither an empty string nor a valid decimal
- sanity margin percentage is neither an empty string nor a valid decimal
- sanity rate is not an empty string and sanity margin percentage is an empty string (in other words, sanity rate is defined but sanity margin percentage is not)
- signers is not one or more valid comma-separated account addresses
- any field is empty, except for order quantity limits, sanity rate, sanity margin percentage, and function parameters for `swapper_function`

This message creates and stores the `War` object at appropriate indexes. Note that the sanity rate and sanity margin percentage are only used in the case of the `swapper_function`, but no error is raised if these are set for other function types.

## MsgEditWar

The owner of a war can edit some of the war's parameters using `MsgEditWar`.

| **Field**              | **Type**           | **Description** |
|:-----------------------|:-------------------|:----------------|
| Token                  | `string`           | The war to be edited
| Name                   | `string`           | Refer to MsgCreateWar
| Description            | `string`           | Refer to MsgCreateWar
| FunctionType           | `string`           | Refer to MsgCreateWar
| OrderQuantityLimits    | `sdk.Coins`        | Refer to MsgCreateWar
| SanityRate             | `sdk.Dec`          | Refer to MsgCreateWar
| SanityMarginPercentage | `sdk.Dec`          | Refer to MsgCreateWar
| Editor                 | `sdk.AccAddress`   | The account address of the user editing the war
| Signers                | `[]sdk.AccAddress` | Refer to MsgCreateWar

This message is expected to fail if:
- any editable field violates the restrictions set for the same field in `MsgCreateWar`
- all editable fields are `"[do-not-modify]"`
- signers list is not equal to the war's signers list

```go
type MsgEditWar struct {
	Token                  string
	Name                   string
	Description            string
	OrderQuantityLimits    string
	SanityRate             string
	SanityMarginPercentage string
	Editor                 sdk.AccAddress
	Signers                []sdk.AccAddress
}
```

This message stores the updated `War` object.

## MsgBuy

Any address that holds tokens that a war uses as its reserve can buy tokens from that war in exchange for reserve tokens. Rather than performing the buy itself, the `MsgBuy` handler registers a buy order in the current orders batch and cancels any other orders that become unfulfillable. Any order in that batch gets fulfilled at the end of the batch's lifespan. The `MsgBuy` handler also locks away the `MaxPrices` value (`< Balance`) indicated by the address so that these are not used elsewhere whilst the batch is being processed.

A buy order is cancelled if the max prices are exceeded at any point during the lifespan of the batch. Otherwise, the buy order is fulfilled. The number of tokens requested are minted on the fly and any remaining tokens from the locked `MaxPrices`, minus the transaction fee specified by the war, are returned to the user. The actual price in reserve tokens charged to the address is determined from the war function, but is also influenced by any other buys and sells in the same orders batch, as a means to prevent front-running.

In the case of `augmented_function` wars, if the war state is `HATCH`, a fixed price-per-token `p0` is used. This value (`p0`) is one of the function parameters required for this function type.

| **Field** | **Type**         | **Description** |
|:----------|:-----------------|:----------------|
| Buyer     | `sdk.AccAddress` | The account address of the user buying the tokens
| Amount    | `sdk.Coin`       | The amount of war tokens to be bought
| MaxPrices | `sdk.Coins`      | The max price to pay in reserve tokens

This message is expected to fail if:
- amount is not an amount of an existing war
- war state is not HATCH or OPEN
- max prices is greater than the balance of the buyer
- max prices are not amounts of the war's reserve tokens
- denominations in max prices are not the war's reserve tokens
- buyer does not afford to buy the tokens at the current price
- amount causes the war's batch-adjusted current supply to exceed the max supply
- amount violates an order quantity limit defined by the war

The batch-adjusted current supply in the case of buys is the current supply of the war plus any uncancelled buy amounts in the current batch. 

```go
type MsgBuy struct {
	Buyer     sdk.AccAddress
	Amount    sdk.Coin
	MaxPrices sdk.Coins
}
```

This message adds the buy order to the current batch.

### MsgBuy for Swapper Function Wars

In general, but especially in the case of swapper function wars, buying tokens from a war can be seen as adding liquidity to that war's token. To add liquidity to a swapper function, the current exchange rate is used to determine how much of each reserve token makes up the price. Otherwise, the price is an equal number of each of the reserve tokens according to the function type.

Moreover, in the case of the swapper function, the first `MsgBuy` performed is special and plays a very important role in specifying the price of the war token. Since we have no price reference for the first buy in a swapper function, the `MaxPrices` specified are used as the actual price, with no fees charged.

This effectively means that if the user requested `n` war tokens with max prices `aR1` and `bR2` (for reserve tokens `R1` and `R2`), the next buyers will have to pay `(a/n)R1` and `(b/n)R2` tokens per war token requested. Specifying high `a` and `b` prices for a small `n` (say `n=1`) means that the next buyers will have to pay at most `aR1` and `bR2` per war token. **Thus, it is important that the first buy is well-calculated and performed carefully.**

## MsgSell

Any address that holds previously bought war tokens can, at any point, sell the tokens back to the war in exchange for reserve tokens. Similar to the `MsgBuy`, the `MsgSell` handler just registers a sell order in the current orders batch which then gets fulfilled at the end of the batch's lifespan.

Once the sell order is fulfilled, the number of tokens to be sold are burned on the fly and the address gets reserve tokens in return, minus the transaction and exit fees specified by the war. The actual number of reserve tokens given to the address in return is determined from the war function, but is also influenced by any other buys and sells in the same orders batch, as a means to prevent front-running. A sell order cannot be cancelled.

In general, but especially in the case of swapper function wars, buying tokens from a war can be seen as adding liquidity for that war. To add liquidity to a swapper function, the current exchange rate is used to determine how much of each reserve token makes up the price. Otherwise, the price is an equal number of each of the reserve tokens according to the function type.

| **Field** | **Type**         | **Description** |
|:----------|:-----------------|:----------------|
| Seller    | `sdk.AccAddress` | The account address of the user selling the tokens
| Amount    | `sdk.Coin`       | The amount of war tokens to be sold

This message is expected to fail if:
- amount is not an amount of an existing war
- war state is not OPEN
- amount is greater than the balance of the seller
- amount is greater than the war's current supply
- amount causes the war's batch-adjusted current supply to become negative
- amount violates an order quantity limit defined by the war
- war function type is `augmented_function` and war state is `HATCH`

The batch-adjusted current supply in the case of sells is the current supply of the war minus any uncancelled sell amounts in the current batch.

```go
type MsgSell struct {
	Seller sdk.AccAddress
	Amount sdk.Coin
}
```

This message adds the sell order to the current batch.

## MsgSwap

Any address that holds tokens (_t1_) that a swapper function war uses as one of its two reserves (_t1_ and _t2_) can swap the tokens in exchange for reserve tokens of the other type (_t2_). Similar to the `MsgBuy` and `MsgSell`, the `MsgSwap` handler just registers a swap order in the current orders batch which then gets fulfilled at the end of the batch's lifespan.

Once the swap order is fulfilled, 

| **Field** | **Type**         | **Description** |
|:----------|:-----------------|:----------------|
| Swapper   | `sdk.AccAddress` | The account address of the user swapping the tokens
| WarToken | `string`         | The swapper function war to use to perform the swap
| From      | `sdk.Coin`       | The amount of reserve tokens to be swapped
| ToToken   | `string`         | The token denomination that will be given in return

This message is expected to fail if:
- war does not exist, is not swapper function, or war state is not OPEN
- from amount is greater than the balance of the swapper
- from and to tokens are the same token
- from and to tokens are not the swapper function's reserve tokens
- from amount violates an order quantity limit defined by the war

```go
type MsgSwap struct {
	Swapper   sdk.AccAddress
	WarToken string
	From      sdk.Coin
	ToToken   string
}
```

This message adds the swap order to the current batch.

## MsgMakeOutcomePayment

If a war was created with an outcome payment field, then any token holder can make an outcome payment to the war. If the token holder has enough tokens to pay the outcome payment, the tokens are sent to the war's reserve and the war's state gets set to SETTLE. The only action possible by war token holders after the outcome payment has been made is a share withdrawal (using [MsgWithdrawShare](#MsgWithdrawShare)).

| **Field** | **Type**         | **Description**                                                                                               |
|:----------|:-----------------|:--------------------------------------------------------------------------------------------------------------|
| Sender    | `sdk.AccAddress` | The account address of the user making the outcome payment |
| WarToken | `string`         | The war to make the outcome payment to                    |

This message is expected to fail if:
- war does not exist or war state is not OPEN
- war outcome payment is empty (meaning the feature is disabled)
- war outcome payment is greater than the balance of the sender

```go
type MsgMakeOutcomePayment struct {
	Sender    sdk.AccAddress
	WarToken string
}
```

## MsgWithdrawShare

If a war's outcome payment was paid, any war token holder can use this message to get their share of the reserve. The amount owed to the war token holder is calculated by considering the percentage of war tokens owned as a fraction of the _remaining_ war token supply. Examples:

- If the war token holder owns 100% of all war tokens and the reserve has 1000 reserve tokens, then the war token holder gets all 1000 reserve tokens.
- If three war token holders each own 1/3 of all war tokens and the reserve has 1000 reserve tokens, then:
  - The first token holder to withdraw gets `1000/3 = 333 tokens` (notice the rounding down from 333.33)
  - The second token holder to withdraw gets `667/2 = 333 tokens` (notice the current supply is now 2)
  - The third token holder to withdraw gets `334/1 = 334 tokens` (because of rounding, the last holder got an extra token)

| **Field** | **Type**         | **Description**                                                                                               |
|:----------|:-----------------|:--------------------------------------------------------------------------------------------------------------|
| Recipient | `sdk.AccAddress` | The account address of the user withdrawing their share |
| WarToken | `string`         | The war to withdraw the share from                     |

This message is expected to fail if:
- war does not exist or war state is not SETTLE
- recipient does not own any war tokens

```go
type MsgWithdrawShare struct {
	Recipient sdk.AccAddress
	WarToken string
}
```
