#!/usr/bin/env bash

wait() {
  echo "Waiting for chain to start..."
  while :; do
    RET=$(warscli status 2>&1)
    if [[ ($RET == ERROR*) || ($RET == *'"latest_block_height": "0"'*) ]]; then
      sleep 1
    else
      echo "A few more seconds..."
      sleep 6
      break
    fi
  done
}

rest_from_m() {
  curl -s -X POST localhost:1317/"$1" --data-binary "$2" >tx.json                   # generate
  yes $PASSWORD | warscli tx sign tx.json --from=miguel --output-document=tx.json  # sign
  python3 demo_rest_fix_tx_format.py                                                # reformat
  curl -s -X POST localhost:1317/txs --data-binary "$(<tx.json)"                    # broadcast
  rm tx.json
}

rest_from_f() {
  curl -s -X POST localhost:1317/"$1" --data-binary "$2" >tx.json                      # generate
  yes $PASSWORD | warscli tx sign tx.json --from=francesco --output-document=tx.json  # sign
  python3 demo_rest_fix_tx_format.py                                                   # reformat
  curl -s -X POST localhost:1317/txs --data-binary "$(<tx.json)"                       # broadcast
  rm tx.json
}

query_war() {
  curl -X GET localhost:1317/wars/"$1"
}

query_account() {
  curl -X GET localhost:1317/auth/accounts/"$1"
}

RET=$(warscli status 2>&1)
if [[ ($RET == ERROR*) || ($RET == *'"latest_block_height": "0"'*) ]]; then
  wait
fi

PASSWORD="12345678"
GAS_PRICES="0.025stake"
MIGUEL=$(yes $PASSWORD | warscli keys show miguel --keyring-backend=test -a)
FRANCESCO=$(yes $PASSWORD | warscli keys show francesco --keyring-backend=test -a)
SHAUN=$(yes $PASSWORD | warscli keys show shaun --keyring-backend=test -a)
FEE=$(yes $PASSWORD | warscli keys show fee --keyring-backend=test -a)

echo "Creating war..."
# shellcheck disable=SC2046
rest_from_m wars/create_war '{
                                "base_req":{
                                  "from":"'$MIGUEL'",
                                  "chain_id":"warschain-1",
                                  "gas_prices":[{"denom":"stake","amount":"0.025"}]
                                },
                                "token":"abc",
                                "name":"A B C",
                                "description":"Description about A B C",
                                "function_type":"swapper_function",
                                "function_parameters":"",
                                "reserve_tokens":"res,rez",
                                "tx_fee_percentage":"0.5",
                                "exit_fee_percentage":"0.1",
                                "fee_address":"'$FEE'",
                                "max_supply":"1000000abc",
                                "order_quantity_limits":"",
                                "sanity_rate":"0",
                                "sanity_margin_percentage":"0",
                                "allow_sells":"true",
                                "signers":"'$MIGUEL'",
                                "batch_blocks":"1"
                              }'
echo "Created war..."
query_war abc

echo "Editing war..."
# shellcheck disable=SC2046
rest_from_m wars/edit_war '{
                              "base_req":{
                                "from":"'$MIGUEL'",
                                "chain_id":"warschain-1",
                                "gas_prices":[{"denom":"stake","amount":"0.025"}]
                              },
                              "token":"abc",
                              "name":"New A B C",
                              "description":"New description about A B C",
                              "sanity_rate":"[do-not-modify]",
                              "sanity_margin_percentage":"[do-not-modify]",
                              "order_quantity_limits":"[do-not-modify]",
                              "signers":"'$MIGUEL'"
                            }'
echo "Edited war..."
query_war abc

echo "Miguel buys 1abc..."
# shellcheck disable=SC2046
rest_from_m wars/buy '{
                        "base_req":{
                          "from":"'$MIGUEL'",
                          "chain_id":"warschain-1",
                          "gas_prices":[{"denom":"stake","amount":"0.025"}]
                        },
                        "war_token":"abc",
                        "war_amount":"1",
                        "max_prices":"500res,1000rez"
                      }'
echo "Miguel's account..."
query_account "$MIGUEL"

echo "Francesco buys 10abc..."
# shellcheck disable=SC2046
rest_from_f wars/buy '{
                        "base_req":{
                          "from":"'$FRANCESCO'",
                          "chain_id":"warschain-1",
                          "gas_prices":[{"denom":"stake","amount":"0.025"}]
                        },
                        "war_token":"abc",
                        "war_amount":"10",
                        "max_prices":"10100res,10100rez"
                      }'
echo "Francesco's account..."
query_account "$FRANCESCO"

echo "Miguel swap 5000 res to rez..."
# shellcheck disable=SC2046
rest_from_m wars/swap '{
                        "base_req":{
                          "from":"'$MIGUEL'",
                          "chain_id":"warschain-1",
                          "gas_prices":[{"denom":"stake","amount":"0.025"}]
                        },
                        "war_token":"abc",
                        "from_amount":"5000",
                        "from_token":"res",
                        "to_token":"rez"
                      }'
echo "Miguel's account..."
warscli q auth account "$MIGUEL"

echo "Francesco swap 5000 rez to res..."
# shellcheck disable=SC2046
rest_from_f wars/swap '{
                          "base_req":{
                            "from":"'$FRANCESCO'",
                            "chain_id":"warschain-1",
                            "gas_prices":[{"denom":"stake","amount":"0.025"}]
                          },
                          "war_token":"abc",
                          "from_amount":"5000",
                          "from_token":"rez",
                          "to_token":"res"
                        }'
echo "Francesco's account..."
warscli q auth account "$FRANCESCO"

echo "Miguel sells 1abc..."
# shellcheck disable=SC2046
rest_from_m wars/sell '{
                          "base_req":{
                            "from":"'$MIGUEL'",
                            "chain_id":"warschain-1",
                            "gas_prices":[{"denom":"stake","amount":"0.025"}]
                          },
                          "war_token":"abc",
                          "war_amount":"1"
                        }'
echo "Miguel's account..."
query_account "$MIGUEL"

echo "Francesco sells 10abc..."
# shellcheck disable=SC2046
rest_from_f wars/sell '{
                          "base_req":{
                            "from":"'$FRANCESCO'",
                            "chain_id":"warschain-1",
                            "gas_prices":[{"denom":"stake","amount":"0.025"}]
                          },
                          "war_token":"abc",
                          "war_amount":"10"
                        }'
echo "Francesco's account..."
query_account "$FRANCESCO"
