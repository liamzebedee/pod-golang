package core

import (
	"container/list"
	"context"
	"fmt"
	"io"
	"sort"
	"sync"

	"github.com/liamzebedee/pod-go/core/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// Downstream consumers use the Client to interact with replicas
type Client struct {
	// List of all replicas in the pod.
	replicas []*ReplicaInfo

	// A lookup for transaction timestamps, indexed by transaction ID and replica.
	// transaction -> replica -> timestamp
	// termed "tsps" in the paper.
	timestamps map[pb.TXID]map[ReplicaID]timestamp

	// The most recent timestamp returned by each replica
	// termed "mrt" in the paper.
	mostRecentTimestamp map[ReplicaID]timestamp

	// The next sequence number expected by each replica
	// termed "nextsn" in the paper.
	nextSeqNum map[ReplicaID]int64

	// Vote backlog.
	voteBacklog *list.List

	// Vote lock.
	// Serialises votes through ingestVote to be processed one-by-one.
	voteLock sync.Mutex
}

func NewClient() Client {
	return Client{
		replicas:            []*ReplicaInfo{},
		timestamps:          make(map[pb.TXID]map[ReplicaID]timestamp),
		mostRecentTimestamp: make(map[ReplicaID]timestamp),
		nextSeqNum:          make(map[ReplicaID]int64),
		voteBacklog:         list.New(),
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

		// NOTE: We never close the connection, as we run a streaming RPC.
		// defer conn.Close()

		// 2. Create RPC client.
		replicaServiceClient := pb.NewReplicaServiceClient(conn)

		// 3. Connect to replica to verify liveness.
		_, err = replicaServiceClient.Connect(context.Background(), &pb.Empty{})
		if err != nil {
			panic(err)
		}

		// 4. Save details.
		cl.replicas = append(cl.replicas, &ReplicaInfo{
			DialAddress:          replicaInfo.DialAddress,
			PK:                   replicaInfo.PK,
			ReplicaServiceClient: replicaServiceClient,
		})
	}

	// Setup replica context.
	for _, r := range cl.replicas {
		cl.mostRecentTimestamp[r.ID()] = 0
		cl.nextSeqNum[r.ID()] = 0
	}

	// 4. Start replica routines.
	for _, rep := range cl.replicas {
		// Stream the votes from replica to client.
		stream, err := rep.ReplicaServiceClient.StreamVotes(context.Background(), &pb.Empty{})
		got := status.Code(err)
		if got != codes.OK {
			fmt.Printf("client error connecting to replica: %v\n", err)
			continue
		}

		fmt.Printf("client connected to replica: %s\n", rep.ID())

		go func(rep *ReplicaInfo, stream grpc.ServerStreamingClient[pb.Vote]) {
			// For each vote.
			for {
				vote, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Printf("client error streaming votes: %v\n", err)
				}

				cl.ingestVote(vote, rep)
			}
		}(rep, stream)
	}
}

func (cl *Client) ingestVote(vote *pb.Vote, rep *ReplicaInfo) {
	// Lock
	cl.voteLock.Lock()
	defer cl.voteLock.Unlock()

	// Handle vote
	fmt.Printf("client got vote: seq=%d replica=%s heartbeat=%v\n", vote.Sn, rep.ID(), vote.IsHeartbeat)

	// 1. Process backlog for any missing sequence numbers.
	cl.processBacklog()

	// 2. Check vote sequence and add to backlog if necessary.
	backlogged := cl.voteCheckSequence(vote, rep)
	if backlogged {
		return
	}

	// 3. If vote is not backlogged, process it.
	err := cl.receiveVote(vote, rep)
	if err != nil {
		fmt.Printf("client error receiving vote: %v\n", err)
	}
}

// Process the vote backlog.
// The vote backlog is implemented as a doubly linked list.
// This algorithm is O(N) where N is the number of votes in the backlog.
// Just a simple linear scan - can be improved later.
func (cl *Client) processBacklog() {
	if cl.voteBacklog.Len() == 0 {
		return
	}

	for e := cl.voteBacklog.Front(); e != nil; e = e.Next() {
		backlogItem := e.Value.(VoteInBacklog)
		vote := backlogItem.vote
		replicaID := backlogItem.replica.ID()

		// Skip if not next in sequence.
		if vote.Sn != cl.nextSeqNum[replicaID] {
			continue
		}

		// Remove from backlog.
		cl.voteBacklog.Remove(e)

		// Process backlog item.
		err := cl.receiveVote(backlogItem.vote, backlogItem.replica)
		if err != nil {
			fmt.Printf("Error processing vote in backlog: %v\n", err)
			// Vote is discarded in the case of error.
			// Errors in processing the vote are non-recoverable, so we can discard safely here.
		}
	}
}

func (cl *Client) voteCheckSequence(vote *pb.Vote, replica *ReplicaInfo) (addedToBacklog bool) {
	// If vote is ready to be processed, don't add to backlog.
	if vote.Sn == cl.nextSeqNum[replica.ID()] {
		return false
	}

	// Backlog vote.
	cl.voteBacklog.PushBack(VoteInBacklog{
		vote:    vote,
		replica: replica,
	})
	fmt.Printf("adding vote to backlog expected=%d got=%d backlog.len=%d\n", cl.nextSeqNum[replica.ID()], vote.Sn, cl.voteBacklog.Len())
	return true
}

func (cl *Client) receiveVote(vote *pb.Vote, replica *ReplicaInfo) error {
	// Verify signature.
	isSigValid := VerifySignature(replica.PK, vote, vote.Sig)
	if !isSigValid {
		return fmt.Errorf("invalid signature")
	}

	// Verify sequence number.
	if vote.Sn != cl.nextSeqNum[replica.ID()] {
		return fmt.Errorf("vote out-of-sequence, expected=%d got=%d", cl.nextSeqNum[replica.ID()], vote.Sn)
	}

	// Verify timestamp.
	if vote.Ts <= cl.mostRecentTimestamp[replica.ID()] {
		// Rj sent old timestamp
		return fmt.Errorf("invalid timestamp")
	}

	cl.nextSeqNum[replica.ID()] += 1
	// Update most recent timestamp.
	cl.mostRecentTimestamp[replica.ID()] = vote.Ts

	if vote.GetIsHeartbeat() {
		// Heartbeat vote.
		return nil
	}

	// Update transaction timestamp.
	txid := vote.Tx.ID()

	// Upsert (txid, map(replica->timestamp)) if not already exist.
	if _, ok := cl.timestamps[txid]; !ok {
		cl.timestamps[vote.Tx.ID()] = make(map[ReplicaID]timestamp)
	}

	// Check map(replica->timestamp) for pre-existing key.
	// If timestamps already contains this transaction and timestamp, discard it as duplicate
	if ts2, ok := cl.timestamps[txid][replica.ID()]; ok && ts2 != vote.Ts {
		// Duplicate vote.
		// TODO is paper correct here? This seems wrong.
		return fmt.Errorf("duplicate vote for tx")
	}

	// Store timestamp.
	cl.timestamps[txid][replica.ID()] = vote.Ts

	return nil
}

func (cl *Client) Write(tx *pb.Transaction) {
	// Send transaction to all replicas.
	for _, replica := range cl.replicas {
		go func(replica *ReplicaInfo) {
			vote, err := replica.ReplicaServiceClient.Write(context.Background(), tx)
			if err != nil {
				// skip.
				fmt.Printf("Error writing to replica: %v", err)
				return
			}

			// Process vote.
			cl.ingestVote(vote, replica)
		}(replica)
	}
}

type ReadResponse struct {
	// T is a set of transactions with their associated timestamps
	Txs []*pb.Transaction

	// The past-perfect round
	// "The past-perfection safety property of pod guarantees that T contains all transactions that every other honest party will ever read with a confirmed round smaller than rperf"
	// TLDR: forall tx in T, tx.RConf <= rperf
	RPerf timestamp

	// Auxiliary data.
	// Cpp is a set of votes that the client has received from replicas
	// For now, commenting this out as we don't need it.
	// Cpp []*pb.Vote
}

func (cl *Client) Read() ReadResponse {
	var T []*pb.Transaction

	// 1. Compute rmin, rmax, rconf for each transaction.
	for tx, replicas := range cl.timestamps {
		// 1. rmin
		rMin := MinPossibleTimestamp(
			tx,
			cl.timestamps,
			getReplicaIds(cl.replicas),
			pAlpha, pBeta,
			cl.mostRecentTimestamp,
		)

		// 2. rmax
		rMax := MaxPossibleTimestamp(
			tx,
			cl.timestamps,
			getReplicaIds(cl.replicas),
			pAlpha, pBeta,
		)

		// 3. rconf
		//
		var rConf timestamp = -1
		var timestampsList []timestamp
		// var Ctx []*pb.Vote

		// If we have a quorum of votes, calculate rconf.
		nVotes := len(replicas)
		if pAlpha <= nVotes {
			// 1. Get all timestamps.
			for _, ts := range replicas {
				timestampsList = append(timestampsList, ts)
				// Ctx = append(Ctx, ts.getVoteMsg())
			}

			// 2. Sort timestamps.
			sort.Float64s(timestampsList)

			// 3. Calculate median.
			rConf = Median(timestampsList)
		}

		T = append(T, &pb.Transaction{
			// TODO.
			RMin: rMin, RMax: rMax, RConf: rConf, //, Ctx: Ctx
		})
	}

	rPerf := MinPossibleTimestampForNewTx(cl.mostRecentTimestamp, pAlpha, pBeta)
	// for _, ts := range cl.mostRecentTimestamp {
	// 	Cpp = append(Cpp, ts.getVoteMsg())
	// }

	fmt.Printf("Read rperf=%f\n", rPerf)
	for _, tx := range T {
		fmt.Printf("%s\n", tx.FormatString())
	}

	return ReadResponse{T, rPerf}
}

type ReplicaInfo struct {
	DialAddress          string
	PK                   PublicKey
	ReplicaServiceClient pb.ReplicaServiceClient
}

func (r *ReplicaInfo) ID() ReplicaID {
	return r.PK.String()
}

type VoteInBacklog struct {
	vote    *pb.Vote
	replica *ReplicaInfo
}

func getReplicaIds(replicas []*ReplicaInfo) []ReplicaID {
	ids := []ReplicaID{}
	for _, r := range replicas {
		ids = append(ids, r.ID())
	}
	return ids
}
