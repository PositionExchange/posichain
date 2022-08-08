FROM alpine

RUN apk add --no-cache bash libstdc++ gmp-dev libc6-compat bind-tools jq && ln -s libcrypto.so.1.1 /lib/libcrypto.so.10

# default base port, rpc port and rest port
EXPOSE 9000/tcp 9500/tcp 9800/tcp 6000/tcp

VOLUME ["/posichain/posichain_db_0","/posichain/log","/posichain/.psc"]

# Default BN for S3/mainnet
ENV NODE_BN_MNET "/ip4/34.91.46.25/tcp/9876/p2p/QmeeirwM67Q7pDf96KTJnML4w5EDMfYwupaatCXPVAUBio,/ip4/34.90.230.57/tcp/9876/p2p/QmQ62LRJciHiwxRW2iibAcHxmociPBDQCttsw71P5EWZMd"

ENV NODE_PORT "9000"
ENV NODE_BLSKEY ""
ENV NODE_BLSPASS ""
ENV NODE_DNS_ZONE "posichain.org"
ENV NODE_RPC "true"
ENV NODE_BLACKLIST ""
ENV NODE_NETWORK_TYPE "mainnet"

ENTRYPOINT ["/bin/run"]
WORKDIR /posichain

COPY run /bin/run
COPY posichain /bin/
