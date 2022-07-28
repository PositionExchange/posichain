#!/bin/bash

if [[ "$ENABLE_REMOTE_DEBUG" == "true" ]]; then
  dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient --log=true --log-output=debugger,debuglineerr,gdbwire,lldbout,rpc exec ./bin/posichain -- --bls.pass.file=passphrases/passphrase.txt -c=./posichain.conf
else
  ./bin/posichain --bls.pass.file=passphrases/passphrase.txt -c=./posichain.conf
fi
