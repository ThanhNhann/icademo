#!/bin/bash

BINARY=icademod
CHAIN_DIR=./data
CHAINID_1=test-1
CHAINID_2=test-2
GRPCPORT_1=8090
GRPCPORT_2=9090
GRPCWEB_1=8091
GRPCWEB_2=9091
RPCPORT_1=26657
RPCPORT_2=26667
APIPORT_1=1317
APIPORT_2=1318

# Function to wait for a node to be ready
wait_for_node() {
    local port=$1
    local max_attempts=30
    local attempt=1
    
    echo "Waiting for node on port $port to be ready..."
    while [ $attempt -le $max_attempts ]; do
        if curl -s "http://localhost:$port/status" > /dev/null; then
            echo "Node on port $port is ready!"
            return 0
        fi
        echo "Attempt $attempt/$max_attempts: Node not ready yet, waiting..."
        sleep 2
        attempt=$((attempt + 1))
    done
    echo "Node failed to start within the timeout period"
    return 1
}

# Function to check if a process is running
check_process() {
    local pid=$1
    if ps -p $pid > /dev/null; then
        return 0
    else
        return 1
    fi
}

echo "Starting $CHAINID_1 in $CHAIN_DIR..."
echo "Creating log file at $CHAIN_DIR/$CHAINID_1.log"
$BINARY start --log_level trace --log_format json --home $CHAIN_DIR/$CHAINID_1 --pruning=nothing \
    --grpc.address="0.0.0.0:$GRPCPORT_1" \
    --grpc-web.address="0.0.0.0:$GRPCWEB_1" \
    --rpc.laddr="tcp://0.0.0.0:$RPCPORT_1" \
    --api.address="tcp://0.0.0.0:$APIPORT_1" > $CHAIN_DIR/$CHAINID_1.log 2>&1 &
PID1=$!

echo "Starting $CHAINID_2 in $CHAIN_DIR..."
echo "Creating log file at $CHAIN_DIR/$CHAINID_2.log"
$BINARY start --log_level trace --log_format json --home $CHAIN_DIR/$CHAINID_2 --pruning=nothing \
    --grpc.address="0.0.0.0:$GRPCPORT_2" \
    --grpc-web.address="0.0.0.0:$GRPCWEB_2" \
    --rpc.laddr="tcp://0.0.0.0:$RPCPORT_2" \
    --api.address="tcp://0.0.0.0:$APIPORT_2" > $CHAIN_DIR/$CHAINID_2.log 2>&1 &
PID2=$!

# Wait for both nodes to be ready
if ! wait_for_node $RPCPORT_1; then
    echo "Error: First node failed to start"
    if check_process $PID1; then
        kill $PID1
    fi
    if check_process $PID2; then
        kill $PID2
    fi
    exit 1
fi

if ! wait_for_node $RPCPORT_2; then
    echo "Error: Second node failed to start"
    if check_process $PID1; then
        kill $PID1
    fi
    if check_process $PID2; then
        kill $PID2
    fi
    exit 1
fi

echo "Both nodes are ready!"
