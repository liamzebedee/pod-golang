pod-go
======

Implementation of the [Pod consensus protocol](https://x.com/liamzebedee/status/1890324486181843111).

## Overview.

Pod consensus is a partially-ordered, single roundtrip, byzantine fault tolerant consensus protocol.

It is implemented in Go, using gRPC/Protobufs for the networking.

## Design.

**Overview**:

 - this is a real implementation, not a simulation of the concepts.
 - nodes run a real RPC protocol (gRPC) and share no internal state.
 - signatures are mocked for testing; in production, you would use ECDSA (trivial to implement).
 - votes are streamed from replicas to client.
 - system parameterisation (number of replicas + failure tolerance) is in `params.go`.

**Differences from paper**:

 - heartbeats are sent at regular intervals (150ms) rather than every millisecond.
 - read auxiliary data (`Cpp`) is removed - no need for this.
 - `getVoteMessage` helper is removed - only needed for auxiliary data.
 - `forall replicas, client.nextSeqNum[replica] = 0`. In the paper, replicas set their initial sequence number to 0 whereas clients expect -1. This doesn't compute lol. Simple off-by-one error.
 - vote processing inside the client is slightly changed. In the paper, the client updates state like the sequence number BEFORE doing some validations. This means errors can leave the client in an unclean state, which doesn't seem right.
 - heartbeats are represented by a simple boolean flag on the tx.

**Missing from paper**:

 - Unique transaction ID's. The paper makes mention of "session ID's", which are not implemented here (didn't appear to be needed?). Transactions however did need a unique identifier - so in line with Bitcoin, we've chosen to use SHA256 of the transaction payload.

**Missing from implementation**:

 - clients should re-sync with replicas if they experience churn.
 - replicas/clients have no persistent state; in production, they would save their state to disk and load it on start.
 - there is no discrete timestep simulation. so you cannot test eg. network conditions

## Progress.

 - [x] Read paper
 - [x] Flesh out basic types - replica, client, signatures.
 - [x] Flesh out networking - gRPC, streaming responses.
 - [x] Write out the timestamp logic.
 - [x] Write out client.write
 - [x] Write out pod.Read
 - [x] Develop PoC integration test.
 - [ ] Write tests for the read/write/process vote routines.
 - [ ] Test `pod.Read()`
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

