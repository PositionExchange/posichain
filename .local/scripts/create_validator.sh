#!/bin/bash

./bin/hmycli --node="http://localhost:9500" --chain-id mainnet staking create-validator \
    --validator-addr one13gya32xnff6xpege024jn30qmdvctmrjjt2jjz --amount 10000 \
    --bls-pubkeys 7e552cc0562b3c08220e120cc6fab96c9aa1c0a734720e8359c7e5d20ba006cb5bfa92aa2ac9309404972715888ca596 \
    --name "Node1" --identity "node1" --details "Node1 validator" \
    --security-contact "node1@posichain.com" --website "position.exchange" \
    --max-change-rate 0.1 --max-rate 0.1 --rate 0.1 \
    --max-total-delegation 100000000 --min-self-delegation 10000

./bin/hmycli --node="http://localhost:9500" --chain-id mainnet staking create-validator \
    --validator-addr one1axv87jgexskmnh0fqpguaxecp4nu30mscwrw42 --amount 10000 --gas-limit 5400000 \
    --bls-pubkeys 64d4359bd948fe8895d97997551894cb1217c84584f844003cf19e54057d95fe537f18acb88f3ca9e1023ce68261ef00 \
    --name "Node3" --identity "node3" --details "Node3 validator" \
    --security-contact "node3@posichain.com" --website "position.exchange" \
    --max-change-rate 0.1 --max-rate 0.1 --rate 0.1 \
    --max-total-delegation 100000000 --min-self-delegation 10000