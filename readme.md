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
