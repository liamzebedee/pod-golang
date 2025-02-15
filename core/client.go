package core

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/liamzebedee/pod-go/core/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ReplicaInfo struct {
	DialAddress          string
	PK                   PublicKey
	ReplicaServiceClient pb.ReplicaServiceClient
}

// Downstream consumers use the Client to interact with replicas
type Client struct {
	replicas []ReplicaInfo

	// transaction -> replica -> timestamp
	// termed "tsps" in the paper.
	timestamps map[Transaction]map[*ReplicaInfo]timestamp

	// termed "mrt" in the paper.
	mostRecentTimestamp map[*ReplicaInfo]timestamp

	// termed "nextsn" in the paper.
	nextSeqNum map[*ReplicaInfo]int
}

func NewClient() Client {
	return Client{
		replicas:            []ReplicaInfo{},
		timestamps:          make(map[Transaction]map[*ReplicaInfo]timestamp),
		mostRecentTimestamp: make(map[*ReplicaInfo]timestamp),
		nextSeqNum:          make(map[*ReplicaInfo]int),
	}
}

func (cl *Client) Start(infos []ReplicaInfo) {
	// Connect to all replicas.
	for _, replicaInfo := range infos {
		// 1. Connect to RPC.
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		conn, err := grpc.NewClient(replicaInfo.DialAddress, opts...)
		if err != nil {
			panic(err)
		}
		// TODO.
		// defer conn.Close()

		// 2. Create RPC client.
		replicaServiceClient := pb.NewReplicaServiceClient(conn)

		// 3. Save details.
		cl.replicas = append(cl.replicas, ReplicaInfo{
			DialAddress:          replicaInfo.DialAddress,
			PK:                   replicaInfo.PK,
			ReplicaServiceClient: replicaServiceClient,
		})
	}

	// Setup replica context.
	for _, r := range cl.replicas {
		cl.mostRecentTimestamp[&r] = 0
		cl.nextSeqNum[&r] = -1
	}

	// 4. Start replica routines.
	for _, rep := range cl.replicas {
		go func(rep ReplicaInfo) {
			stream, err := rep.ReplicaServiceClient.StreamVotes(context.Background(), &pb.Empty{})
			if err != nil {
				panic(err)
			}

			for {
				vote, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatalf("%v.ListFeatures(_) = _, %v", rep.ReplicaServiceClient, err)
				}

				// Handle vote
				fmt.Println(vote)
			}
		}(rep)
	}
}

func makeTx(data byte) *pb.Transaction {
	return &pb.Transaction{
		Ctx:   []byte{data},
		RMin:  0,
		RMax:  0,
		RConf: 0,
	}
}

func (cl *Client) Write(tx *pb.Transaction) {
	// Send transaction to all replicas.
	for _, replica := range cl.replicas {
		_, err := replica.ReplicaServiceClient.Write(context.Background(), tx)
		if err != nil {
			// skip.
			fmt.Printf("Error writing to replica: %v", err)
			continue
		}

		// Collect votes.
	}
}

func (cl *Client) Read() {

}

func (cl *Client) startup() {
	// 1. Send CONNECT to all Replicas.
	// At initialization the client also sends a ⟨CONNECT⟩ message to each replica, which initiates a streaming connection from the replica to the client.

	// A client maintains a connection to each replica and receives votes through ⟨VOTE (tx, ts, sn, σ, Rj )⟩ messages (lines 15–24).

	// When a vote is received from replica Rj , the client first verifies the signature σ under Rj ’s public key (line 16). If invalid, the vote is ignored. Then the client verifies that the vote contains the next sequence number it expects to receive from replica Rj (line 17). If this is not the case, the vote is backlogged and given again to the client at a later point (the backlogging functionality is not shown in the pseudocode)

	// Heartbeat messages. Clients update their most-recent timestamp mrt[Rj] every time they receive a vote from replica Rj (line 20 in Algorithm 1). However, in case Rj has not written any transaction in round r (simply because no client called write() in round r), Rj will advance its round number, but clients will not advance mrt[Rj]. We solve this by having replicas send a vote on a dummy heartBeat transaction the end of each round (lines 26–28). An obvious practi- cal implementation is to send heartBeat only for rounds when no other transactions were sent. When received by a client, a heartBeat is handled as a vote (i.e., it triggers line 15 in Algorithm 1), except that we do not check for duplicate heartBeat votes and do not update any associated values for the heartbeat transaction (see line 21 in Algorithm 1).

	// When clients read the pod, they obtain a pod data structure D = (T, rperf, Cpp), where T is set of transactions with their associated timestamps, rperf is a past- perfect round and Cpp is auxiliary data.
}
