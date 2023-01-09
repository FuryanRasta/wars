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

tx_from_s() {
  cmd=$1
  shift
  yes $PASSWORD | warscli tx wars "$cmd" --from shaun --keyring-backend=test -y --broadcast-mode block --gas-prices="$GAS_PRICES" "$@"
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

# d0 := 500.0   // initial raise (reserve)
# p0 := 0.01    // initial price (reserve per token)
# theta := 0.4  // initial allocation (percentage)
# kappa := 3.0  // degrees of polynomial (i.e. x^2, x^4, x^6)

# R0 = 300              // initial reserve (1-theta)*d0
# S0 = 50000            // initial supply
# V0 = 416666666666.667 // invariant

echo "Creating war..."
tx_from_m create-war \
  --token=abc \
  --name="A B C" \
  --description="Description about A B C" \
  --function-type=augmented_function \
  --function-parameters="d0:500.0,p0:0.01,theta:0.4,kappa:3.0" \
  --reserve-tokens=res \
  --tx-fee-percentage=0 \
  --exit-fee-percentage=0 \
  --fee-address="$FEE" \
  --max-supply=1000000abc \
  --order-quantity-limits="" \
  --sanity-rate="0" \
  --sanity-margin-percentage="0" \
  --allow-sells \
  --signers="$MIGUEL" \
  --batch-blocks=1 \
  --outcome-payment="100000res"
echo "Created war..."
warscli q wars war abc

echo "Miguel buys 20000abc..."
tx_from_m buy 20000abc 100000res
echo "Miguel's account..."
warscli q auth account "$MIGUEL"

echo "Francesco buys 20000abc..."
tx_from_f buy 20000abc 100000res
echo "Francesco's account..."
warscli q auth account "$FRANCESCO"

echo "Shaun cannot buy 10001abc..."
tx_from_s buy 10001abc 100000res
echo "Shaun cannot sell anything..."
tx_from_s sell 10000abc
echo "Shaun can buy 10000abc..."
tx_from_s buy 10000abc 100000res
echo "Shaun's account..."
warscli q auth account "$SHAUN"

echo "War state is now open..."  # since 50000 (S0) reached
warscli q wars war abc

echo "Miguel sells 20000abc..."
tx_from_m sell 20000abc
echo "Miguel's account..."
warscli q auth account "$MIGUEL"

echo "Francesco makes outcome payment..."
tx_from_f make-outcome-payment abc
echo "Francesco's account..."
warscli q auth account "$FRANCESCO"

echo "Francesco withdraws share..."
tx_from_f withdraw-share abc
echo "Francesco's account..."
warscli q auth account "$FRANCESCO"

echo "Shaun withdraws share..."
tx_from_s withdraw-share abc
echo "Shaun's account..."
warscli q auth account "$SHAUN"
