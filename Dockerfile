FROM ubuntu:18.04

ARG TARGETARCH
ARG GOLANG_VERSION="1.17"

SHELL ["/bin/bash", "-c"]

ENV GOPATH=/root/go
ENV GO111MODULE=on
ENV BASE_PATH=${GOPATH}/src/github.com/PositionExchange
ENV OPENSSL_DIR=/usr/lib/ssl
ENV MCL_DIR=${BASE_PATH}/mcl
ENV BLS_DIR=${BASE_PATH}/bls
ENV CGO_CFLAGS="-I${BLS_DIR}/include -I${MCL_DIR}/include"
ENV CGO_LDFLAGS="-L${BLS_DIR}/lib"
ENV LD_LIBRARY_PATH=${BLS_DIR}/lib:${MCL_DIR}/lib
ENV GIMME_GO_VERSION=${GOLANG_VERSION}
ENV PATH="/root/bin:${PATH}"

RUN apt update && apt upgrade -y && \
	apt install libgmp-dev libssl-dev curl git \
	psmisc dnsutils jq make gcc g++ bash tig tree sudo vim \
	silversearcher-ag unzip emacs-nox nano bash-completion -y

RUN mkdir ~/bin && \
	curl -sL -o ~/bin/gimme \
	https://raw.githubusercontent.com/travis-ci/gimme/master/gimme && \
	chmod +x ~/bin/gimme

RUN eval "$(~/bin/gimme ${GIMME_GO_VERSION})"

RUN touch /root/.bash_profile && \
	gimme ${GIMME_GO_VERSION} >> /root/.bash_profile && \
	echo "GIMME_GO_VERSION='${GIMME_GO_VERSION}'" >> /root/.bash_profile && \
	echo "GO111MODULE='on'" >> /root/.bash_profile && \
	echo ". ~/.bash_profile" >> /root/.profile && \
	echo ". ~/.bash_profile" >> /root/.bashrc

ENV PATH="/root/.gimme/versions/go${GIMME_GO_VERSION}.linux.${TARGETARCH:-amd64}/bin:${GOPATH}/bin:${PATH}"

RUN . ~/.bash_profile;
RUN go get -u golang.org/x/tools/cmd/goimports
RUN go get -u golang.org/x/lint/golint
RUN go get -u github.com/rogpeppe/godef
RUN go get -u github.com/go-delve/delve/cmd/dlv
RUN go get -u github.com/golang/mock/mockgen
RUN go get -u github.com/stamblerre/gocode
RUN go get -u golang.org/x/tools/...
RUN go get -u honnef.co/go/tools/cmd/staticcheck/...

RUN git clone https://github.com/PositionExchange/bls.git ${BASE_PATH}/bls
RUN git clone https://github.com/PositionExchange/mcl.git ${BASE_PATH}/mcl
RUN cd ${BASE_PATH}/bls && make -j8 BLS_SWAP_G=1

RUN git clone https://github.com/PositionExchange/posichain-gosdk.git ${BASE_PATH}/posichain-gosdk
RUN cd ${BASE_PATH}/posichain-gosdk && make -j8 && cp psc /root/bin

WORKDIR ${BASE_PATH}/posichain
COPY . .
RUN scripts/install_build_tools.sh
RUN go mod tidy
RUN scripts/go_executable_build.sh -S

ARG K1=one1tq4hy947c9gr8qzv06yxz4aeyhc9vn78al4rmu
ARG K2=one1y5gmmzumajkm5mx3g2qsxtza2d3haq0zxyg47r
ARG K3=one1qrqcfek6sc29sachs3glhs4zny72mlad76lqcp

ARG KS1=8d222cffa99eb1fb86c581d9dfe7d60dd40ec62aa29056b7ff48028385270541
ARG KS2=da1800da5dedf02717696675c7a7e58383aff90b1014dfa1ab5b7bd1ce3ef535
ARG KS3=f4267bb5a2f0e65b8f5792bb6992597fac2b35ebfac9885ce0f4152c451ca31a

RUN psc keys import-private-key ${KS1} && \
	psc keys import-private-key ${KS2} && \
	psc keys import-private-key ${KS3} && \
	psc keys generate-bls-key > keys.json

RUN jq  '.["encrypted-private-key-path"]' -r keys.json > /root/keypath && cp keys.json /root && \
	echo "export BLS_KEY_PATH=$(cat /root/keypath)" >> /root/.bashrc && \
	echo "export BLS_KEY=$(jq '.["public-key"]' -r keys.json)" >> /root/.bashrc && \
	echo "printf '${K1}, ${K2}, ${K3} are imported accounts in psc for local dev\n\n'" >> /root/.bashrc && \
	echo "printf 'test with: psc blockchain validator information ${K1}\n\n'" >> /root/.bashrc && \
	echo "echo "$(jq '.["public-key"]' -r keys.json)" is an extern bls key" >> /root/.bashrc && \
	echo ". /etc/bash_completion" >> /root/.bashrc && \
	echo ". <(psc completion)" >> /root/.bashrc
