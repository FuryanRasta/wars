# Events

The wars module emits the following events:

## EndBlocker

| Type | Attribute Key | Attribute Value |
| :--- | :--- | :--- |
| order\_cancel | war | {token} |
| order\_cancel | order\_type | {orderType} |
| order\_cancel | address | {address} |
| order\_cancel | cancel\_reason | {cancelReason} |
| order\_fulfill | war | {token} |
| order\_fulfill | order\_type | {orderType} |
| order\_fulfill | address | {address} |
| order\_fulfill | tokensMinted | {tokensMinted} |
| order\_fulfill | chargedPrices | {chargedPrices} |
| order\_fulfill | chargedFees | {chargedFees} |
| order\_fulfill | returnedToAddress | {returnedToAddress} |
| state\_change | war | {token} |
| state\_change | old\_state | {oldState} |
| state\_change | new\_state | {newState} |

## Handlers

### MsgCreateWar

| Type | Attribute Key | Attribute Value |
| :--- | :--- | :--- |
| create\_war | war | {token} |
| create\_war | name | {name} |
| create\_war | description | {description} |
| create\_war | function\_type | {functionType} |
| create\_war | function\_parameters \[0\] | {functionParameters} |
| create\_war | reserve\_tokens \[1\] | {reserveTokens} |
| create\_war | tx\_fee\_percentage | {txFeePercentage} |
| create\_war | exit\_fee\_percentage | {exitFeePercentage} |
| create\_war | fee\_address | {feeAddress} |
| create\_war | max\_supply | {maxSupply} |
| create\_war | order\_quantity\_limits | {orderQuantityLimits} |
| create\_war | sanity\_rate | {sanityRate} |
| create\_war | sanity\_margin\_percentage | {sanityMarginPercentage} |
| create\_war | allow\_sells | {allowSells} |
| create\_war | signers \[2\] | {signers} |
| create\_war | batch\_blocks | {batchBlocks} |
| create\_war | state | {state} |
| message | module | wars |
| message | action | create\_war |
| message | sender | {senderAddress} |

* \[0\] Example formatting: `"{m:12,n:2,c:100}"`
* \[1\] Example formatting: `"[res,rez]"`
* \[2\] Example formatting: `"[ADDR1,ADDR2]"`

### MsgEditWar

| Type | Attribute Key | Attribute Value |
| :--- | :--- | :--- |
| edit\_war | war | {token} |
| edit\_war | name | {name} |
| edit\_war | description | {description} |
| edit\_war | order\_quantity\_limits | {orderQuantityLimits} |
| edit\_war | sanity\_rate | {sanityRate} |
| edit\_war | sanity\_margin\_percentage | {sanityMarginPercentage} |
| message | module | wars |
| message | action | edit\_war |
| message | sender | {senderAddress} |

### MsgBuy

#### First Buy for Swapper Function War

| Type | Attribute Key | Attribute Value |
| :--- | :--- | :--- |
| init\_swapper | war | {token} |
| init\_swapper | amount | {amount} |
| init\_swapper | charged\_prices | {chargedPrices} |
| message | module | wars |
| message | action | buy |
| message | sender | {senderAddress} |

#### Otherwise

| Type | Attribute Key | Attribute Value |
| :--- | :--- | :--- |
| buy | war | {token} |
| buy | amount | {amount} |
| buy | max\_prices | {maxPrices} |
| order\_cancel | war | {token} |
| order\_cancel | order\_type | {orderType} |
| order\_cancel | address | {address} |
| order\_cancel | cancel\_reason | {cancelReason} |
| message | module | wars |
| message | action | buy |
| message | sender | {senderAddress} |

### MsgSell

| Type | Attribute Key | Attribute Value |
| :--- | :--- | :--- |
| sell | war | {token} |
| sell | amount | {amount} |
| message | module | wars |
| message | action | buy |
| message | sender | {senderAddress} |

### MsgSwap

| Type | Attribute Key | Attribute Value |
| :--- | :--- | :--- |
| swap | war | {token} |
| swap | amount | {amount} |
| swap | from\_token | {fromToken} |
| swap | to\_token | {toToken} |
| message | module | wars |
| message | action | swap |
| message | sender | {senderAddress} |

### MsgMakeOutcomePayment

| Type | Attribute Key | Attribute Value |
| :--- | :--- | :--- |
| make\_outcome\_payment | war | {token} |
| make\_outcome\_payment | address | {senderAddress} |
| message | module | wars |
| message | action | make\_outcome\_payment |
| message | sender | {senderAddress} |

### MsgWithdrawShare

| Type | Attribute Key | Attribute Value |
| :--- | :--- | :--- |
| withdraw\_share | war | {token} |
| withdraw\_share | address | {recipientAddress} |
| withdraw\_share | amount | {reserveOwed} |
| message | module | wars |
| message | action | withdraw\_share |
| message | sender | {recipientAddress} |

