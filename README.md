pod-go
======

Implementation of the [Pod consensus protocol](https://x.com/liamzebedee/status/1890324486181843111).

## Overview.

Pod consensus is a partially-ordered, single roundtrip, byzantine fault tolerant consensus protocol.

It is implemented in Go, using gRPC/Protobufs for the networking.

## Setup.

Requirements:

 * Go v1.22 or greater.

For contributing, you will need the protocol buffers toolchain:

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

# Run integration test.
cd core/
go test
```

## Demo.

```sh
➜  core git:(main) ✗ go test
got random port: 50056
got random port: 50057
got random port: 50058
got random port: 50059
got random port: 50060
client connected to replica: a2091fc44c23154e7e443ad6ca97dac2621f8d5e89d623a50842eee23172013a
client connected to replica: 71094f0a812b21340e36a3745b5b57408ba62c2339c474835db60d9171899b24
client connected to replica: 77c2d430fcceec03de9243ef19f7fb608aa91ae3269c43a1b29ccb2e65c2688c
client connected to replica: c2bf8ca85058cafc024f1dd1a671c794a9ccc1abadd3b6a045cc86e34d0d9c32
client connected to replica: 933bf06f44926cb230872ca0c08b3786c264de6e1fd63a87781b50a75270f197
appending client client=&{0x140003642a0}
appending client client=&{0x1400047c1c0}
appending client client=&{0x140000fe7e0}
replica write tx=ctx:"\x01" ts=1739874225550.512939 sn=0
replica write tx=ctx:"\x01" ts=1739874225550.559082 sn=0
appending client client=&{0x140001ea2a0}
replica write tx=ctx:"\x01" ts=1739874225550.594971 sn=0
replica write tx=ctx:"\x01" ts=1739874225550.591064 sn=0
appending client client=&{0x1400022e0e0}
replica write tx=ctx:"\x01" ts=1739874225550.606201 sn=0
client got vote: seq=0 replica=a2091fc44c23154e7e443ad6ca97dac2621f8d5e89d623a50842eee23172013a heartbeat=false
client got vote: seq=0 replica=933bf06f44926cb230872ca0c08b3786c264de6e1fd63a87781b50a75270f197 heartbeat=false
client got vote: seq=0 replica=71094f0a812b21340e36a3745b5b57408ba62c2339c474835db60d9171899b24 heartbeat=false
client got vote: seq=0 replica=c2bf8ca85058cafc024f1dd1a671c794a9ccc1abadd3b6a045cc86e34d0d9c32 heartbeat=false
client got vote: seq=0 replica=77c2d430fcceec03de9243ef19f7fb608aa91ae3269c43a1b29ccb2e65c2688c heartbeat=false
client got vote: seq=0 replica=933bf06f44926cb230872ca0c08b3786c264de6e1fd63a87781b50a75270f197 heartbeat=false
adding vote to backlog expected=1 got=0 backlog.len=1
client got vote: seq=0 replica=77c2d430fcceec03de9243ef19f7fb608aa91ae3269c43a1b29ccb2e65c2688c heartbeat=false
adding vote to backlog expected=1 got=0 backlog.len=2
client got vote: seq=0 replica=c2bf8ca85058cafc024f1dd1a671c794a9ccc1abadd3b6a045cc86e34d0d9c32 heartbeat=false
adding vote to backlog expected=1 got=0 backlog.len=3
client got vote: seq=0 replica=71094f0a812b21340e36a3745b5b57408ba62c2339c474835db60d9171899b24 heartbeat=false
adding vote to backlog expected=1 got=0 backlog.len=4
replica write tx=ctx:"\x02" ts=1739874225651.753906 sn=1
replica write tx=ctx:"\x02" ts=1739874225651.750000 sn=1
replica write tx=ctx:"\x02" ts=1739874225651.747070 sn=1
replica write tx=ctx:"\x02" ts=1739874225651.747070 sn=1
replica write tx=ctx:"\x02" ts=1739874225651.793945 sn=1
client got vote: seq=1 replica=c2bf8ca85058cafc024f1dd1a671c794a9ccc1abadd3b6a045cc86e34d0d9c32 heartbeat=false
client got vote: seq=1 replica=c2bf8ca85058cafc024f1dd1a671c794a9ccc1abadd3b6a045cc86e34d0d9c32 heartbeat=false
adding vote to backlog expected=2 got=1 backlog.len=5
client got vote: seq=1 replica=71094f0a812b21340e36a3745b5b57408ba62c2339c474835db60d9171899b24 heartbeat=false
client got vote: seq=1 replica=71094f0a812b21340e36a3745b5b57408ba62c2339c474835db60d9171899b24 heartbeat=false
adding vote to backlog expected=2 got=1 backlog.len=6
client got vote: seq=1 replica=a2091fc44c23154e7e443ad6ca97dac2621f8d5e89d623a50842eee23172013a heartbeat=false
client got vote: seq=1 replica=933bf06f44926cb230872ca0c08b3786c264de6e1fd63a87781b50a75270f197 heartbeat=false
client got vote: seq=1 replica=a2091fc44c23154e7e443ad6ca97dac2621f8d5e89d623a50842eee23172013a heartbeat=false
adding vote to backlog expected=2 got=1 backlog.len=7
client got vote: seq=1 replica=933bf06f44926cb230872ca0c08b3786c264de6e1fd63a87781b50a75270f197 heartbeat=false
adding vote to backlog expected=2 got=1 backlog.len=8
client got vote: seq=1 replica=77c2d430fcceec03de9243ef19f7fb608aa91ae3269c43a1b29ccb2e65c2688c heartbeat=false
client got vote: seq=1 replica=77c2d430fcceec03de9243ef19f7fb608aa91ae3269c43a1b29ccb2e65c2688c heartbeat=false
adding vote to backlog expected=2 got=1 backlog.len=9
client got vote: seq=2 replica=c2bf8ca85058cafc024f1dd1a671c794a9ccc1abadd3b6a045cc86e34d0d9c32 heartbeat=true
client got vote: seq=2 replica=77c2d430fcceec03de9243ef19f7fb608aa91ae3269c43a1b29ccb2e65c2688c heartbeat=true
client got vote: seq=2 replica=a2091fc44c23154e7e443ad6ca97dac2621f8d5e89d623a50842eee23172013a heartbeat=true
client got vote: seq=2 replica=71094f0a812b21340e36a3745b5b57408ba62c2339c474835db60d9171899b24 heartbeat=true
client got vote: seq=2 replica=933bf06f44926cb230872ca0c08b3786c264de6e1fd63a87781b50a75270f197 heartbeat=true
replica write tx=ctx:"\x03" ts=1739874225752.715088 sn=3
replica write tx=ctx:"\x03" ts=1739874225752.729004 sn=3
replica write tx=ctx:"\x03" ts=1739874225752.743896 sn=3
replica write tx=ctx:"\x03" ts=1739874225752.718018 sn=3
replica write tx=ctx:"\x03" ts=1739874225752.726074 sn=3
client got vote: seq=3 replica=933bf06f44926cb230872ca0c08b3786c264de6e1fd63a87781b50a75270f197 heartbeat=false
client got vote: seq=3 replica=c2bf8ca85058cafc024f1dd1a671c794a9ccc1abadd3b6a045cc86e34d0d9c32 heartbeat=false
client got vote: seq=3 replica=c2bf8ca85058cafc024f1dd1a671c794a9ccc1abadd3b6a045cc86e34d0d9c32 heartbeat=false
adding vote to backlog expected=4 got=3 backlog.len=10
client got vote: seq=3 replica=933bf06f44926cb230872ca0c08b3786c264de6e1fd63a87781b50a75270f197 heartbeat=false
adding vote to backlog expected=4 got=3 backlog.len=11
client got vote: seq=3 replica=71094f0a812b21340e36a3745b5b57408ba62c2339c474835db60d9171899b24 heartbeat=false
client got vote: seq=3 replica=71094f0a812b21340e36a3745b5b57408ba62c2339c474835db60d9171899b24 heartbeat=false
adding vote to backlog expected=4 got=3 backlog.len=12
client got vote: seq=3 replica=a2091fc44c23154e7e443ad6ca97dac2621f8d5e89d623a50842eee23172013a heartbeat=false
client got vote: seq=3 replica=77c2d430fcceec03de9243ef19f7fb608aa91ae3269c43a1b29ccb2e65c2688c heartbeat=false
client got vote: seq=3 replica=a2091fc44c23154e7e443ad6ca97dac2621f8d5e89d623a50842eee23172013a heartbeat=false
adding vote to backlog expected=4 got=3 backlog.len=13
client got vote: seq=3 replica=77c2d430fcceec03de9243ef19f7fb608aa91ae3269c43a1b29ccb2e65c2688c heartbeat=false
adding vote to backlog expected=4 got=3 backlog.len=14
Read rperf=1739874225752.718018
Transaction{ID: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855, RMin: 1739874225651.747070, RConf: 1739874225651.750000, RMax: 1739874225651.753906, Range: 0.006836, Ctx: []}
Transaction{ID: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855, RMin: 1739874225752.718018, RConf: 1739874225752.726074, RMax: 1739874225752.729004, Range: 0.010986, Ctx: []}
Transaction{ID: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855, RMin: 1739874225550.559082, RConf: 1739874225550.591064, RMax: 1739874225550.594971, Range: 0.035889, Ctx: []}
client got vote: seq=4 replica=933bf06f44926cb230872ca0c08b3786c264de6e1fd63a87781b50a75270f197 heartbeat=true
client got vote: seq=4 replica=77c2d430fcceec03de9243ef19f7fb608aa91ae3269c43a1b29ccb2e65c2688c heartbeat=true
```

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

This implementation is a working demo of the Pod protocol. I now am focused on implementing tests to assert functionality, and building a simple demo on top.

 - [x] Read paper
 - [x] Flesh out basic types - replica, client, signatures.
 - [x] Flesh out networking - gRPC, streaming responses.
 - [x] Write out the timestamp logic.
 - [x] Write out client.write
 - [x] Write out pod.Read
 - [x] Develop PoC integration test.
 - [x] Test `pod.Read()`
 - [ ] Write unit tests for the read/write/process vote routines.
 - [ ] Figure out what partial ordering is useful for / if this can have global total ordering.

Things not really on the roadmap but could be:

 - [ ] Slashing.
