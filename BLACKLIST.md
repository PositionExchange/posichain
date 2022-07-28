# Blacklist info

The black list is a newline delimited file of wallet addresses. It can also support comments with the `#` character.

## Default Location

By default, the posichain binary looks for the file `./.psc/blaklist.txt`.

## Example File
```
0x806171f95C5a74371a19e8a312c9e5Cb4E1D24f6
0xE1217E2a4861DD5D50983DaD32474Bbfd6A7333F  # This is a comment
0x1D44424803e7D258D3B5F160807c3dF1ec2F0BF8

```

## Details

Each transaction added to the tx-pool has its `to` and `from` address checked against this blacklist.
If there is a hit, the transaction is considered invalid and is dropped from the tx-pool.
