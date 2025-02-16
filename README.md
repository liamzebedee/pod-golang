pod-go
======

WIP implementation of the [Pod consensus protocol](https://x.com/liamzebedee/status/1890324486181843111).

## Overview.

Pod consensus is a partially-ordered, single roundtrip, byzantine fault tolerant consensus protocol.

It is implemented in Go, using gRPC/Protobufs for the networking.

## Progress.

 - [x] Read paper
 - [x] Flesh out basic types - replica, client, signatures.
 - [x] Flesh out networking - gRPC, streaming responses.
 - [x] Write out the timestamp logic.
 - [x] Write out client.write
 - [x] Write out pod.Read
 - [ ] Develop PoC integration test.
 - [ ] Develop small demonstration of something useful.
 - [ ] Figure out what partial ordering is useful for / if this can have global total ordering.

Things not really on the roadmap but could be:

 - [ ] Slashing.

## Setup.

Requirements:

 * Go v1.22 or greater.
 * Protocol buffers compiler (`protoc` is found in `vendor/` for mac m2).
 * Protocol buffers plugins for Go and gRPC.
```sh
# https://protobuf.dev/downloads/
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

## Usage.

```sh
# Generate protobufs and gRPC types.
make proto
```

