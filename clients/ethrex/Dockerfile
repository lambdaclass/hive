ARG baseimage=ethrex
ARG tag=latest

FROM $baseimage:$tag as builder

# Install script tools.
RUN apt-get update -y
RUN apt-get install -y bash curl jq clang pkg-config libssl-dev linux-tools-generic

# Installs Cargo
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y --no-modify-path && \
    ln -s $HOME/.cargo/bin/* /usr/local/bin/

RUN cargo install flamegraph

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
# RUN /usr/local/bin/ethrex --version | sed -e 's/ethrex \(.*\)/\1/' > /version.txt
RUN echo "ethrex/0.1.0" > /version.txt

# Export the usual networking ports to allow outside access to the node.
EXPOSE 8545 8546 30303 30303/udp

ENTRYPOINT ["/ethrex.sh"]
