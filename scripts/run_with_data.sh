#!/usr/bin/env bash

PASSWORD="12345678"
GAS_PRICES="0.025stake"

wars init local --chain-id warschain-1

yes $PASSWORD | warscli keys delete miguel --keyring-backend=test --force
yes $PASSWORD | warscli keys delete francesco --keyring-backend=test --force
yes $PASSWORD | warscli keys delete shaun --keyring-backend=test --force
yes $PASSWORD | warscli keys delete fee --keyring-backend=test --force
yes $PASSWORD | warscli keys delete fee2 --keyring-backend=test --force
yes $PASSWORD | warscli keys delete fee3 --keyring-backend=test --force
yes $PASSWORD | warscli keys delete fee4 --keyring-backend=test --force
yes $PASSWORD | warscli keys delete fee5 --keyring-backend=test --force

yes $PASSWORD | warscli keys add miguel --keyring-backend=test
yes $PASSWORD | warscli keys add francesco --keyring-backend=test
yes $PASSWORD | warscli keys add shaun --keyring-backend=test
yes $PASSWORD | warscli keys add fee --keyring-backend=test
yes $PASSWORD | warscli keys add fee2 --keyring-backend=test
yes $PASSWORD | warscli keys add fee3 --keyring-backend=test
yes $PASSWORD | warscli keys add fee4 --keyring-backend=test
yes $PASSWORD | warscli keys add fee5 --keyring-backend=test

# Note: important to add 'miguel' as a genesis-account since this is the chain's validator
yes $PASSWORD | wars add-genesis-account $(warscli keys show miguel --keyring-backend=test -a) 200000000stake,1000000res,1000000rez
yes $PASSWORD | wars add-genesis-account $(warscli keys show francesco --keyring-backend=test -a) 100000000stake,1000000res,1000000rez
yes $PASSWORD | wars add-genesis-account $(warscli keys show shaun --keyring-backend=test -a) 100000000stake,1000000res,1000000rez

# Set min-gas-prices
FROM="minimum-gas-prices = \"\""
TO="minimum-gas-prices = \"0.025stake\""
sed -i "s/$FROM/$TO/" "$HOME"/.wars/config/app.toml

warscli config chain-id warschain-1
warscli config output json
warscli config indent true
warscli config trust-node true
warscli config keyring-backend test

yes $PASSWORD | wars gentx --name miguel --keyring-backend=test

wars collect-gentxs
wars validate-genesis

# Uncomment the below to broadcast node RPC endpoint
#FROM="laddr = \"tcp:\/\/127.0.0.1:26657\""
#TO="laddr = \"tcp:\/\/0.0.0.0:26657\""
#sed -i "s/$FROM/$TO/" "$HOME"/.wars/config/config.toml

# Uncomment the below to broadcast REST endpoint
# Do not forget to comment the bottom lines !!
# wars start --pruning "everything" &
# warscli rest-server --chain-id warschain-1 --laddr="tcp://0.0.0.0:1317" --trust-node && fg

wars start --pruning "everything" &
warscli rest-server --chain-id warschain-1 --trust-node && fg
