# Build
FROM golang:1.18 AS build

RUN apt update -y && \
    apt install libgmp-dev libssl-dev git -y && \
    apt -y clean all

ENV GOPATH=/root/go
ENV GO111MODULE=on
ENV HMY_PATH=${GOPATH}/src/github.com/PositionExchange
RUN mkdir -p $HMY_PATH

WORKDIR $HMY_PATH

RUN git clone https://github.com/PositionExchange/posichain.git && \
    git clone https://github.com/PositionExchange/bls.git && \
    git clone https://github.com/PositionExchange/mcl.git

WORKDIR $HMY_PATH/posichain

RUN go mod download

RUN go mod tidy

RUN make linux_static && \
    cp ./bin/posichain /root/posichain && \
    cp ./rosetta/infra/run.sh /root/run.sh && \
    cp ./rosetta/infra/rclone.conf /root/rclone.conf

# TODO copy testnet conf
RUN cp ./rosetta/infra/posichain-mainnet.conf /root/posichain-mainnet.conf && \
    cp ./.psc/rosetta_local_fix.csv /root/rosetta_local_fix.csv

# Execution
FROM ubuntu:20.04

RUN apt update -y && \
    apt install libgmp-dev libssl-dev ca-certificates rclone -y && \
    apt -y clean all \

WORKDIR /root

COPY --from=build /root/posichain /root/posichain
COPY --from=build /root/run.sh /root/run.sh
COPY --from=build /root/rclone.conf /root/.config/rclone/rclone.conf
# TODO copy testnet conf
COPY --from=build /root/posichain-mainnet.conf /root/posichain-mainnet.conf
COPY --from=build /root/rosetta_local_fix.csv /root/rosetta_local_fix.csv

EXPOSE 9500/tcp
EXPOSE 9700/tcp
EXPOSE 9800/tcp
ENTRYPOINT ["/bin/bash","/root/run.sh"]
