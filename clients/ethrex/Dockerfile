ARG baseimage=ghcr.io/lambdaclass/ethrex
ARG tag=latest

FROM $baseimage:$tag AS builder

FROM ubuntu:24.04

WORKDIR /usr/local/bin
COPY --from=builder /usr/local/bin/ethrex .

# Install script tools.
RUN apt-get update -y
RUN apt-get install -y bash curl jq

# Add genesis mapper script.
ADD genesis.json /genesis.json
ADD mapper.jq /mapper.jq

# Add the startup script.
ADD ethrex.sh /ethrex.sh
RUN chmod +x /ethrex.sh

# Add the enode URL retriever script.
ADD enode.sh /hive-bin/enode.sh
RUN chmod +x /hive-bin/enode.sh

# Create version.txt
RUN ethrex --version | sed -e 's/ethrex \(.*\)/\1/' > /version.txt

# Export the usual networking ports to allow outside access to the node.
EXPOSE 8545 8546 30303 30303/udp

ENTRYPOINT ["/ethrex.sh"]
