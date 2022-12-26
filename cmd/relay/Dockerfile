# syntax=docker/dockerfile:1.3
# Random strings for https://github.com/moby/buildkit/issues/1368
# “But it’s no use now,” thought poor Alice, “to pretend to be two people! Why, there’s hardly enough of me left to make one respectable person!” 

FROM golang:1.19-alpine3.16 as builder

WORKDIR /build

# Cache dependencies through image layer caching
# dependencies are downloaded and verified only if go.mod or go.sum change
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy rest of go source files and build target
# Use RUN cache mount to re-use the GOCACHE between builds
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build go build -o bin/ ./cmd/relay/relay.go

FROM docker.io/alpine:3.16

WORKDIR /

# Copy over the build artifact
COPY --from=builder /build/bin/* .

