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

tx_from_m() {
  cmd=$1
  shift
  yes $PASSWORD | warscli tx wars "$cmd" --from miguel --keyring-backend=test -y --broadcast-mode block --gas-prices="$GAS_PRICES" "$@"
}

tx_from_f() {
  cmd=$1
  shift
  yes $PASSWORD | warscli tx wars "$cmd" --from francesco --keyring-backend=test -y --broadcast-mode block --gas-prices="$GAS_PRICES" "$@"
}

create_war_multisig() {
  warscli tx wars create-war \
    --token=abc \
    --name="A B C" \
    --description="Description about A B C" \
    --function-type=power_function \
    --function-parameters="m:12,n:2,c:100" \
    --reserve-tokens=res \
    --tx-fee-percentage=0.5 \
    --exit-fee-percentage=0.1 \
    --fee-address="$FEE" \
    --max-supply=1000000abc \
    --order-quantity-limits="" \
    --sanity-rate="0" \
    --sanity-margin-percentage="0" \
    --allow-sells \
    --signers="$(warscli keys show francesco --keyring-backend=test -a),$(warscli keys show shaun --keyring-backend=test -a)" \
    --batch-blocks=1 \
    --from="$MIGUEL" -y --broadcast-mode block --generate-only >multisig.json
  yes $PASSWORD | warscli tx sign multisig.json --from=francesco --output-document=multisig.json
  yes $PASSWORD | warscli tx sign multisig.json --from=shaun --output-document=multisig.json
  warscli tx broadcast multisig.json
  rm multisig.json
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
create_war_multisig
echo "Waiting a bit..."
sleep 5
echo "Created war..."
warscli q wars war abc

echo "Miguel buys 10abc..."
tx_from_m buy 10abc 1000000res
echo "Miguel's account..."
warscli q auth account "$MIGUEL"

echo "Francesco buys 10abc..."
tx_from_f buy 10abc 1000000res
echo "Francesco's account..."
warscli q auth account "$FRANCESCO"

echo "Miguel sells 10abc..."
tx_from_m sell 10abc
echo "Miguel's account..."
warscli q auth account "$MIGUEL"

echo "Francesco sells 10abc..."
tx_from_f sell 10abc
echo "Francesco's account..."
warscli q auth account "$FRANCESCO"
