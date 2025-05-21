# icademo
**icademo** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).

## ICA Demo Guide

This guide will walk you through setting up and testing the Interchain Account (ICA) functionality.

### Prerequisites

- Go 1.19 or later
- Ignite CLI

### Setup and Initialization
```bash
1. Bootstrap two chains, configure the relayer and create an IBC connection (on top of clients that are created as well)


# go relayer
make init-golang-relayer

2. Start relayer:

```bash
#go relayer
make start-golang-rly
```

### Creating and Managing Interchain Accounts

0. Set up wallet (ref from scripts setup 2 chains)

```bash
# Store the following account addresses within the current shell env
export WALLET_1=$(icad keys show wallet1 -a --keyring-backend test --home ./data/test-1) && echo $WALLET_1;
```

1. Create an ICA account on chain-2 from chain-1:

```bash
icademod tx txdemo register-ica-account connection-0 "" --from $WALLET_1 --chain-id test-1 --home ./data/test-1 --node tcp://localhost:26657 --keyring-backend test --gas auto --gas-adjustment 1.3 -y
```

2. Query the ICA account:

```bash
# Query the ICA account on chain-1
icademod q txdemo interchain-account connection-0 $WALLET_1 --home ./data/test-1 --node tcp://localhost:16657

# Store interchain account address
export ICA_ADDR=$(icademod query txdemo interchain-account connection-0 $WALLET_1 --home ./data/test-1 --node tcp://localhost:26657 -o json | jq -r '.interchain_account_address') && echo $ICA_ADDR
# output: icademo1salycgx0lua8ekvn6ptewx20p3lkunl3a5mvkj2sxyf92dflthvqq36jx2
```

### Sending Interchain Account transactions

```bash
# Query check bank balance behalf of host chain
icademod q bank balances $ICA_ADDR --chain-id test-2 --node tcp://localhost:26667

# balances: []
# pagination:
#   next_key: null
#   total: "0"

# Send token from wallet 3 to ica
icademod tx bank send $WALLET_3 $ICA_ADDR 10000stake --chain-id test-2 --home ./data/test-2 --node tcp://localhost:26667 --keyring-backend test -y

# Query bank balances again 
icademod q bank balances $ICA_ADDR --chain-id test-2 --node tcp://localhost:26667

# balances:
# - amount: "10000"
#   denom: stake
# pagination:
#   next_key: null
#   total: "0"

```

#### Funding the Interchain Account wallet

- Example 1: Staking Delegation

```bash
# currently we not implement icq yet so we have to get host chain's validator by get it from genesis
cat ./data/test-2/config/genesis.json | jq -r '.app_state.genutil.gen_txs[0].body.messages[0].validator_address'
# output: icademovaloper1qnk2n4nlkpw9xfqntladh74w6ujtulwnlj904e

# Submit a staking delegation tx using the interchain account via ibc
icademod tx txdemo submit-tx ./test/staking_test.json connection-0 --from $WALLET_1 --chain-id test-1 --home ./data/test-1 --node tcp://localhost:26657 --keyring-backend test -y
# Wait until the relayer has relayed the packet

# Inspect the staking delegations on the host chain
icademod q staking delegations-to icademovaloper1qnk2n4nlkpw9xfqntladh74w6ujtulwnlj904e --home ./data/test-2 --node tcp://localhost:26657

# delegation_responses:
# - balance:
#     amount: "7000000000"
#     denom: stake
#   delegation:
#     delegator_address: icademo1qnk2n4nlkpw9xfqntladh74w6ujtulwne3l4t4
#     shares: "7000000000.000000000000000000"
#     validator_address: icademovaloper1qnk2n4nlkpw9xfqntladh74w6ujtulwnlj904e
# - balance:
#     amount: "1000"
#     denom: stake
#   delegation:
#     delegator_address: icademo1gwxrlec5u32vdwkyyhqjpc2u0736hks09s6jen6hfa0r93ykp3qs03xwmz
#     shares: "1000.000000000000000000"
#     validator_address: icademovaloper1qnk2n4nlkpw9xfqntladh74w6ujtulwnlj904e
# pagination:
#   next_key: null
#   total: "0

```

### Troubleshooting

If you encounter any issues:

1. Check the node logs in `./data/test-1.log` and `./data/test-2.log`
2. Ensure all ports are available and not in use

### Cleanup

To stop and clean up the demo:

```bash
# Stop the nodes
pkill icademod

# Remove the data directory
rm -rf ./data
```
