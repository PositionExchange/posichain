#!/bin/bash

./bin/hmycli --node="http://localhost:9500" --chain-id mainnet staking create-validator \
    --validator-addr one1v87jtng53yjwnppw6msk9valp7sr3v0l7fencu --amount 10000 --gas-limit 5400000 \
    --bls-pubkeys 0cc8db483787a944e2075702c860633b1ac856ff647dec570862e7eef7b670ff66158ede030599babde44c6b85c38804 \
    --name "Node11" --identity "node11" --details "Node11 validator" \
    --security-contact "node11@posichain.com" --website "position.exchange" \
    --max-change-rate 0.1 --max-rate 0.1 --rate 0.1 \
    --max-total-delegation 100000000 --min-self-delegation 10000