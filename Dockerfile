FROM bash:4.4

ARG RELEASE_VERSION=0.2.0

RUN mkdir -p /app && \
    wget -O /app/wait-go "https://github.com/adrian-gheorghe/wait-go/releases/download/${RELEASE_VERSION}/wait-go-linux" && \
    chmod +x /app/wait-go && \
    cp /app/wait-go /usr/local/bin/wait-go
    
WORKDIR /app