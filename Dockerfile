FROM --platform=$BUILDPLATFORM golang:1.22-alpine3.20 AS build

RUN apk add --no-cache make

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ARG COMMIT
ARG BUILD_TIME

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    make build

FROM alpine:3.20

WORKDIR /app

COPY --from=build /app/bin/golink .
COPY frontend/assets ./frontend/assets

EXPOSE 8080

ENTRYPOINT ["/app/golink"]
