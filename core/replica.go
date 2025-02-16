package core

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net"
	"sync"
	"time"

	"github.com/liamzebedee/pod-go/core/pb"
	"google.golang.org/grpc"
)

type Replica struct {
	pb.ReplicaServiceServer

	keypair Keypair

	// The transaction log of the replica.
	Log []pb.Vote

	// The next sequence number to assign to votes
	nextSeqNum int64

	// The set of all connected clients.
	clients []grpc.ServerStreamingServer[pb.Vote]

	writeMutex sync.Mutex
}

func getTimestamp() float64 {
	return float64(time.Now().UnixNano()) / 1e6
}

func signVote(vote *pb.Vote, keypair Keypair) []byte {
	// 1. Construct envelope to be signed.
	buf := []byte{}

	// buf = vote.tx
	buf = append(buf, vote.Tx.GetCtx()...)

	// buf = vote.tx ++ vote.ts
	buf_ts := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf_ts[:], math.Float64bits(vote.Ts))
	buf = append(buf, buf_ts...)

	// buf = vote.tx ++ vote.ts ++ vote.sn
	buf_sn := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf_sn[:], uint64(vote.Sn))
	buf = append(buf, buf_sn...)

	// 2. Sign the envelope.
	sig := Sign(keypair.S, buf)
	return sig.ToBytes()
}

func NewReplica() Replica {
	keypair := Keygen()
	// votesCh := make(chan pb.Vote)
	return Replica{
		keypair:    keypair,
		Log:        []pb.Vote{},
		nextSeqNum: 0,
		clients:    []grpc.ServerStreamingServer[pb.Vote]{},
	}
}

func (r *Replica) PublicKey() PublicKey {
	return r.keypair.P
}

func (s *Replica) makeHeartbeatVote() (*pb.Vote, error) {
	s.writeMutex.Lock()
	defer s.writeMutex.Unlock()

	// 1. Assign a timestamp to the transaction.
	ts := getTimestamp()

	// 2. Assign a sequence number to the transaction.
	sn := s.nextSeqNum
	s.nextSeqNum++

	// 3. Sign the transaction.
	heartbeatTx := &pb.Transaction{
		Ctx:   []byte{},
		RMin:  0,
		RMax:  0,
		RConf: 0,
	}

	vote := &pb.Vote{
		Tx:          heartbeatTx,
		IsHeartbeat: true,
		Ts:          ts,
		Sn:          sn,
		Sig:         []byte{},
	}
	vote.Sig = signVote(vote, s.keypair)

	return vote, nil
}

func (s *Replica) Write(ctx context.Context, tx *pb.Transaction) (*pb.Vote, error) {
	// When it receives ⟨WRITE tx⟩. a replica first checks whether it has already seen tx, in which case the message is ignored. Otherwise, it assigns tx a timestamp ts equal its local round number and the next available sequence number sn, and signs the message (tx, ts, sn) (line 19). Honest replicas use incremental sequence numbers for each transaction, so if a vote has a larger sequence number than another vote, it will have a larger or equal timestamp than the other. The replica appends (tx,ts,sn,σ) to replicaLog, and sends it via a ⟨VOTE (tx,ts,sn,σ,R)⟩ message to all connected clients (line 22
	s.writeMutex.Lock()
	defer s.writeMutex.Unlock()

	// 1. Assign a timestamp to the transaction.
	ts := getTimestamp()

	// 2. Assign a sequence number to the transaction.
	sn := s.nextSeqNum
	s.nextSeqNum++

	// 3. Sign the transaction.
	vote := &pb.Vote{
		Tx:          tx,
		IsHeartbeat: false,
		Ts:          ts,
		Sn:          sn,
		Sig:         []byte{},
	}
	vote.Sig = signVote(vote, s.keypair)

	// 4. Append the transaction to the log.
	s.Log = append(s.Log, *vote)

	fmt.Println("replica write tx=", tx, "ts=", ts, "sn=", sn)

	// 5. Send vote to all connected clients in background.
	go func() {
		for _, stream := range s.clients {
			go func() {
				err := stream.Send(vote)
				if err != nil {
					fmt.Printf("failed to send vote: %v\n", err)
				}
			}()
		}
	}()

	// 6. Return signed timestamped transaction (vote).
	return vote, nil
}

func (s *Replica) Read(ctx context.Context, _ *pb.Empty) (*pb.ReadResponse, error) {
	return nil, nil
}

func (s *Replica) StreamVotes(_ *pb.Empty, stream grpc.ServerStreamingServer[pb.Vote]) error {
	// Add client to list of connected clients.
	fmt.Printf("appending client client=%v\n", stream)
	s.clients = append(s.clients, stream)

	// Send regular heartbeats.
	// In the paper, this is done at wall clock speed - a round is defined as the granularity of the wall clock.
	// The heartbeat is only meant to be sent if there are no other transactions timestamped in that round (millisecond).
	// I've taken a design decision here which seems to follow the spirit of the paper, but not the letter.
	// The heartbeat is sent every 150ms. Every 1ms would be very fast, and incur a lot of cryptographic overhead in signing.
	// The point of the heartbeat is to keep the bounds of the timestamps in check, ie. literally rMin is the minimum timestamp seen by the client.
	// If we send a heartbeat every 150ms, we can be sure that the rMin is at most 150ms behind the current time.
	// This seems to be a reasonable tradeoff.
	for {
		<-time.After(150 * time.Millisecond)
		heartbeat, err := s.makeHeartbeatVote()
		if err != nil {
			panic(err)
		}
		stream.Send(heartbeat)
	}
}

func (rep *Replica) ListenAndServe(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterReplicaServiceServer(grpcServer, rep)
	grpcServer.Serve(lis)
}

// Each replica maintains a sequence number, which it increments and includes every time it assigns a timestamp to a transaction

// We assume that all messages between clients and replicas are concatenated with a session identifier (sid), which is unique for each concurrent execution of the protocol. Moreover, the sid is implicitly included in all messages signed by the replicas
