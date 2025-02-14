package core

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/liamzebedee/pod-go/core/pb"
	"google.golang.org/grpc"
)

type Replica struct {
	pb.ReplicaServiceServer

	PublicKey PublicKey
	TxLog     []Transaction
}

func (s *Replica) Write(ctx context.Context, tx *pb.Transaction) (*pb.Vote, error) {
	return nil, nil
}

func (s *Replica) Read(ctx context.Context, _ *pb.Empty) (*pb.ReadResponse, error) {
	return nil, nil
}

func (s *Replica) StreamVotes(_ *pb.Empty, stream grpc.ServerStreamingServer[pb.Vote]) error {
	// every 250ms send a vote

	for {
		select {
		case <-time.After(250 * time.Millisecond):
			// send vote
			stream.Send(&pb.Vote{
				Sn: time.Now().Unix(),
			})
		}
	}

	// for _, feature := range s.savedFeatures {
	// 	if inRange(feature.Location, rect) {
	// 	  if err := stream.Send(feature); err != nil {
	// 		return err
	// 	  }
	// 	}
	//   }
	//   return nil
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

// We assume that all messages between clients and replicas are concatenated with a session identifier (sid), which is unique for each concurrent execution of the protocol. Moreover, the sid is im- plicitly included in all messages signed by the replicas
// For a timestamp ts, notation ts.getVoteMsg() denotes the vote message from some replica through which a client obtained timestamp ts. We abstract away the logic of how getVoteMsg() is implemented

func (rep *Replica) run() {
	// When it receives ⟨WRITE tx⟩. a replica first checks whether it has already seen tx, in which case the message is ignored. Otherwise, it assigns tx a timestamp ts equal its local round number and the next available sequence number sn, and signs the message (tx, ts, sn) (line 19). Honest replicas use incremental sequence numbers for each transaction, so if a vote has a larger sequence number than another vote, it will have a larger or equal timestamp than the other. The replica appends (tx,ts,sn,σ) to replicaLog, and sends it via a ⟨VOTE (tx,ts,sn,σ,R)⟩ message to all connected clients (line 22
}
