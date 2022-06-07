#!/bin/bash

./bin/hmycli --node="http://localhost:9500" --chain-id mainnet staking create-validator \
    --validator-addr one13gya32xnff6xpege024jn30qmdvctmrjjt2jjz --amount 10000 \
    --bls-pubkeys 7e552cc0562b3c08220e120cc6fab96c9aa1c0a734720e8359c7e5d20ba006cb5bfa92aa2ac9309404972715888ca596 \
    --name "Node1" --identity "node1" --details "Node1 validator" \
    --security-contact "node1@posichain.com" --website "position.exchange" \
    --max-change-rate 0.1 --max-rate 0.1 --rate 0.1 \
    --max-total-delegation 100000000 --min-self-delegation 10000

./bin/hmycli --node="http://localhost:9500" --chain-id mainnet staking create-validator \
    --validator-addr one107rghkwx0c7g83h74v3t5q5ld8knu9w5aznufq --amount 10000 --gas-limit 5400000 \
    --bls-pubkeys f47b6d2b91eb37a5c0b35520803b86901a086367506e423f1651464191cc95b67b95470b8bc41ea4ad57ead0739fc180 \
    --name "Node3" --identity "node3" --details "Node3 validator" \
    --security-contact "node3@posichain.com" --website "position.exchange" \
    --max-change-rate 0.1 --max-rate 0.1 --rate 0.1 \
    --max-total-delegation 100000000 --min-self-delegation 10000