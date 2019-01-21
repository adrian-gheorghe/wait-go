FROM bash:4.4

ARG RELEASE_VERSION

RUN echo ${RELEASE_VERSION} && \
    mkdir -p /app && \
    wget -O /app/wait-go "https://github.com/adrian-gheorghe/wait-go/releases/download/${RELEASE_VERSION}/wait-go-linux" && \
    chmod +x /app/wait-go && \
    cd /app && \
    ls -al . && \
    ./wait-go --version && \
    cp /app/wait-go /usr/local/bin/wait-go && \
    ls -al /usr/local/bin/
    
WORKDIR /app