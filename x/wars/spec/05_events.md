# Events

The wars module emits the following events:

## EndBlocker

| Type          | Attribute Key     | Attribute Value     |
|---------------|-------------------|---------------------|
| order_cancel  | war              | {token}             |
| order_cancel  | order_type        | {orderType}         |
| order_cancel  | address           | {address}           |
| order_cancel  | cancel_reason     | {cancelReason}      |
| order_fulfill | war              | {token}             |
| order_fulfill | order_type        | {orderType}         |
| order_fulfill | address           | {address}           |
| order_fulfill | tokensMinted      | {tokensMinted}      |
| order_fulfill | chargedPrices     | {chargedPrices}     |
| order_fulfill | chargedFees       | {chargedFees}       |
| order_fulfill | returnedToAddress | {returnedToAddress} |
| state_change  | war              | {token}             |
| state_change  | old_state         | {oldState}          |
| state_change  | new_state         | {newState}          |

## Handlers

### MsgCreateWar

| Type        | Attribute Key            | Attribute Value          |
|-------------|--------------------------|--------------------------|
| create_war | war                     | {token}                  |
| create_war | name                     | {name}                   |
| create_war | description              | {description}            |
| create_war | function_type            | {functionType}           |
| create_war | function_parameters [0]  | {functionParameters}     |
| create_war | reserve_tokens [1]       | {reserveTokens}          |
| create_war | tx_fee_percentage        | {txFeePercentage}        |
| create_war | exit_fee_percentage      | {exitFeePercentage}      |
| create_war | fee_address              | {feeAddress}             |
| create_war | max_supply               | {maxSupply}              |
| create_war | order_quantity_limits    | {orderQuantityLimits}    |
| create_war | sanity_rate              | {sanityRate}             |
| create_war | sanity_margin_percentage | {sanityMarginPercentage} |
| create_war | allow_sells              | {allowSells}             |
| create_war | signers [2]              | {signers}                |
| create_war | batch_blocks             | {batchBlocks}            |
| create_war | state                    | {state}                  |
| message     | module                   | wars                    |
| message     | action                   | create_war              |
| message     | sender                   | {senderAddress}          |

* [0] Example formatting: `"{m:12,n:2,c:100}"`
* [1] Example formatting: `"[res,rez]"`
* [2] Example formatting: `"[ADDR1,ADDR2]"`

### MsgEditWar

| Type      | Attribute Key            | Attribute Value          |
|-----------|--------------------------|--------------------------|
| edit_war | war                     | {token}                  |
| edit_war | name                     | {name}                   |
| edit_war | description              | {description}            |
| edit_war | order_quantity_limits    | {orderQuantityLimits}    |
| edit_war | sanity_rate              | {sanityRate}             |
| edit_war | sanity_margin_percentage | {sanityMarginPercentage} |
| message   | module                   | wars                    |
| message   | action                   | edit_war                |
| message   | sender                   | {senderAddress}          |

### MsgBuy

#### First Buy for Swapper Function War

| Type         | Attribute Key  | Attribute Value |
|--------------|----------------|-----------------|
| init_swapper | war           | {token}         |
| init_swapper | amount         | {amount}        |
| init_swapper | charged_prices | {chargedPrices} |
| message      | module         | wars           |
| message      | action         | buy             |
| message      | sender         | {senderAddress} |

#### Otherwise

| Type         | Attribute Key | Attribute Value |
|--------------|---------------|-----------------|
| buy          | war          | {token}         |
| buy          | amount        | {amount}        |
| buy          | max_prices    | {maxPrices}     |
| order_cancel | war          | {token}         |
| order_cancel | order_type    | {orderType}     |
| order_cancel | address       | {address}       |
| order_cancel | cancel_reason | {cancelReason}  |
| message      | module        | wars           |
| message      | action        | buy             |
| message      | sender        | {senderAddress} |

### MsgSell

| Type    | Attribute Key | Attribute Value |
|---------|---------------|-----------------|
| sell    | war          | {token}         |
| sell    | amount        | {amount}        |
| message | module        | wars           |
| message | action        | buy             |
| message | sender        | {senderAddress} |

### MsgSwap

| Type    | Attribute Key | Attribute Value |
|---------|---------------|-----------------|
| swap    | war          | {token}         |
| swap    | amount        | {amount}        |
| swap    | from_token    | {fromToken}     |
| swap    | to_token      | {toToken}       |
| message | module        | wars           |
| message | action        | swap            |
| message | sender        | {senderAddress} |

### MsgMakeOutcomePayment

| Type                 | Attribute Key | Attribute Value      |
|----------------------|---------------|----------------------|
| make_outcome_payment | war          | {token}              |
| make_outcome_payment | address       | {senderAddress}      |
| message              | module        | wars                |
| message              | action        | make_outcome_payment |
| message              | sender        | {senderAddress}      |

### MsgWithdrawShare

| Type           | Attribute Key | Attribute Value    |
|----------------|---------------|--------------------|
| withdraw_share | war          | {token}            |
| withdraw_share | address       | {recipientAddress} |
| withdraw_share | amount        | {reserveOwed}      |
| message        | module        | wars              |
| message        | action        | withdraw_share     |
| message        | sender        | {recipientAddress} |
