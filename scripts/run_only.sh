#!/usr/bin/env bash

# Uncomment the below to broadcast REST endpoint
# Do not forget to comment the bottom lines !!
# wars start --pruning "everything" &
# warscli rest-server --chain-id warschain-1 --laddr="tcp://0.0.0.0:1317" --trust-node && fg

wars start --pruning "everything" &
warscli rest-server --chain-id warschain-1 --trust-node && fg
