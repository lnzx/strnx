# syntax=docker/dockerfile:1
ARG GO_NAME="go1.20.5.linux-amd64.tar.gz"

FROM debian:stable-slim AS build

ARG GO_NAME

WORKDIR /usr/src

# Install dependencies
RUN apt update && apt install -y \
    git \
    curl \
    && curl -sLO https://go.dev/dl/$GO_NAME \
    && tar -C /usr/local -xzf $GO_NAME \
    && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt install --no-install-recommends -y nodejs \
    && corepack enable && corepack prepare pnpm@latest --activate \
    && git clone --depth 1 https://github.com/lnzx/strnx.git \
    && rm -rf /var/lib/apt/lists/* \
    && rm $GO_NAME

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


