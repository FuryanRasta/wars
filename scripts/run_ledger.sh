#!/usr/bin/env bash

PASSWORD="12345678"
GAS_PRICES="0.025stake"

wars init local --chain-id warschain-1

warscli keys add miguel --ledger
yes $PASSWORD | warscli keys add francesco
yes $PASSWORD | warscli keys add shaun
yes $PASSWORD | warscli keys add reserve
yes $PASSWORD | warscli keys add fee

wars add-genesis-account "$(warscli keys show miguel -a)" 100000000stake,1000000res,1000000rez
wars add-genesis-account "$(warscli keys show francesco -a)" 100000000stake,1000000res,1000000rez
wars add-genesis-account "$(warscli keys show shaun -a)" 100000000stake,1000000res,1000000rez

warscli config chain-id warschain-1
warscli config output json
warscli config indent true
warscli config trust-node true

echo "$PASSWORD" | wars gentx --name miguel

wars collect-gentxs
wars validate-genesis

wars start --pruning "everything" &
warscli rest-server --chain-id warschain-1 --trust-node && fg
