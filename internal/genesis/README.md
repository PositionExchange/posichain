# Genesis Accounts

1. Generate multiple BLS Keys:

```shell
# Using psc cli
./psc keys generate-bls-keys --shard-count=1 --shard=0 --count=25 --write-passphrase=true --passphrase

# Or using docker
docker run -it posichain/psc keys generate-bls-keys --shard-count=1 --shard=0 --count=25 --write-passphrase=true --passphrase
```

2. Generate Golang compatible accounts list

```shell
python keys.py -json-file=localnet-node-accounts.json -out=accounts-generated-go.txt
```
