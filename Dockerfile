# syntax=docker/dockerfile:1
FROM debian:stable-slim AS build

WORKDIR /usr/src

# Install dependencies
RUN apt update && apt install -y \
    git \
    curl \
    wget \
    && wget https://go.dev/dl/go1.20.5.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.20.5.linux-amd64.tar.gz \
    && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install --no-install-recommends -y nodejs \
    && corepack enable && corepack prepare pnpm@latest --activate \
    && git clone --depth 1 https://github.com/lnzx/strnx.git \
    && rm -rf /var/lib/apt/lists/* \
    && rm go1.20.5.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /usr/src/strnx
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app .

WORKDIR /usr/src/strnx/web
RUN pnpm install && pnpm build

FROM debian:stable-slim

WORKDIR /usr/src/app

# 复制 app 和静态文件
COPY --from=build /usr/src/strnx/app .
COPY --from=build /usr/src/strnx/dist/ ./dist/

CMD ["./app"]


