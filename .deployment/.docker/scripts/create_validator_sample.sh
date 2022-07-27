#!/bin/bash

# Node 11
./bin/psc --node="http://172.189.0.13:9500" --chain-id devnet staking create-validator \
    --validator-addr one1v87jtng53yjwnppw6msk9valp7sr3v0l7fencu --amount 10000 --gas-limit 5400000 \
    --bls-pubkeys 0cc8db483787a944e2075702c860633b1ac856ff647dec570862e7eef7b670ff66158ede030599babde44c6b85c38804 \
    --name "Danny Node11" --identity "danny-node11" --details "Validator (node11) by Danny" \
    --security-contact "danny@position.exchange" --website "https://position.exchange" \
    --max-change-rate 0.1 --max-rate 0.1 --rate 0.1 \
    --max-total-delegation 100000000 --min-self-delegation 10000

./bin/psc --node="http://172.189.0.13:9500" --chain-id devnet staking edit-validator \
   --validator-addr one1v87jtng53yjwnppw6msk9valp7sr3v0l7fencu --active true --passphrase

# Node 12
./bin/psc --node="http://172.189.0.13:9500" --chain-id devnet staking create-validator \
    --validator-addr one1cgl23vk0rmctcjp5a7kt0jn00qa8zpeap6jhfj --amount 10000 --gas-limit 5400000 \
    --bls-pubkeys e95d68d4aff05156ce9dd9a987bab6ae26677faf4e8770a66f991030d517b02fbeb526d9634bd7732344ebb40d5beb02 \
    --name "Danny Node12" --identity "danny-node12" --details "Validator (node12) by Danny" \
    --security-contact "danny@position.exchange" --website "https://position.exchange" \
    --max-change-rate 0.1 --max-rate 0.1 --rate 0.1 \
    --max-total-delegation 100000000 --min-self-delegation 10000

./bin/psc --node="http://172.189.0.13:9500" --chain-id devnet staking edit-validator \
   --validator-addr one1cgl23vk0rmctcjp5a7kt0jn00qa8zpeap6jhfj --active true --passphrase

# Node 13
./bin/psc --node="http://172.189.0.13:9500" --chain-id devnet staking create-validator \
    --validator-addr one1kre48u2nacpsv7hxe2l0ud346pluw9e5ugth2r --amount 25000 --gas-limit 5400000 \
    --bls-pubkeys 52d853265ff6014e8693d331b116d1a331a10a1e3f74c9229179b92e739e36f69688ecbc5fb2911d1603cf9274ca2e06 \
    --name "Danny Node13" --identity "danny-node13" --details "Validator (node13) by Danny" \
    --security-contact "danny@position.exchange" --website "https://position.exchange" \
    --max-change-rate 0.1 --max-rate 0.1 --rate 0.1 \
    --max-total-delegation 100000000 --min-self-delegation 10000

./bin/psc --node="http://172.189.0.13:9500" --chain-id devnet staking edit-validator \
   --validator-addr one1kre48u2nacpsv7hxe2l0ud346pluw9e5ugth2r --active true --passphrase
