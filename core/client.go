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
	"google.golang.org/grpc/credentials/insecure"
)

type ReplicaInfo struct {
	DialAddress          string
	PK                   PublicKey
	ReplicaServiceClient pb.ReplicaServiceClient
}

type VoteInBacklog struct {
	vote    *pb.Vote
	replica *ReplicaInfo
}

type ReplicaID = *ReplicaInfo

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
	voteBacklog      *list.List
	voteBacklogMutex sync.Mutex
}

func NewClient() Client {
	return Client{
		replicas:            []*ReplicaInfo{},
		timestamps:          make(map[pb.TXID]map[*ReplicaInfo]timestamp),
		mostRecentTimestamp: make(map[*ReplicaInfo]timestamp),
		nextSeqNum:          make(map[*ReplicaInfo]int64),
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
		// TODO.
		// defer conn.Close()

		// 2. Create RPC client.
		replicaServiceClient := pb.NewReplicaServiceClient(conn)

		// 3. Save details.
		cl.replicas = append(cl.replicas, &ReplicaInfo{
			DialAddress:          replicaInfo.DialAddress,
			PK:                   replicaInfo.PK,
			ReplicaServiceClient: replicaServiceClient,
		})
	}

	// Setup replica context.
	for _, r := range cl.replicas {
		cl.mostRecentTimestamp[r] = 0
		cl.nextSeqNum[r] = 0 // TODO is this right?
	}

	// 4. Start replica routines.
	for _, rep := range cl.replicas {
		go func(rep *ReplicaInfo) {
			// Stream the votes from replica to client.
			stream, err := rep.ReplicaServiceClient.StreamVotes(context.Background(), &pb.Empty{})
			if err != nil {
				panic(err)
			}

			// For each vote.
			for {
				vote, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Printf("error streaming votes: %v\n", err)
				}

				// Handle vote
				fmt.Println(vote)

				// 1. Process backlog for any missing sequence numbers.
				cl.processBacklog(rep)

				// 2. Process vote.
				err = cl.receiveVote(vote, rep, true)
				if err != nil {
					fmt.Printf("Error receiving vote: %v\n", err)
				}
			}
		}(rep)
	}
}

// Process the vote backlog.
// The vote backlog is implemented as a doubly linked list.
// This algorithm is O(N) where N is the number of votes in the backlog.
// Just a simple linear scan - will be improved later.
func (cl *Client) processBacklog(forReplica *ReplicaInfo) {
	cl.voteBacklogMutex.Lock()
	defer cl.voteBacklogMutex.Unlock()

	if cl.voteBacklog.Len() == 0 {
		return
	}

	for e := cl.voteBacklog.Front(); e != nil; e = e.Next() {
		backlogItem := e.Value.(VoteInBacklog)

		// Only process votes for this replica.
		if forReplica != nil && backlogItem.replica != forReplica {
			continue
		}

		if backlogItem.vote.Sn == cl.nextSeqNum[backlogItem.replica] {
			// Remove from backlog.
			cl.voteBacklog.Remove(e)

			// Process backlog item.
			err := cl.receiveVote(backlogItem.vote, backlogItem.replica, false)
			if err != nil {
				fmt.Printf("Error processing vote in backlog: %v\n", err)
				// discard the item.
			}
		}
	}
}

func (cl *Client) receiveVote(vote *pb.Vote, replica *ReplicaInfo, verifySequence bool) error {
	// Verify signature.
	isSigValid := true
	if !isSigValid {
		return fmt.Errorf("Invalid signature")
	}

	// Verify sequence number.
	// If verifySequence is false, we are processing a vote from the backlog.
	if verifySequence && vote.Sn != cl.nextSeqNum[replica] {
		// Backlog vote.
		cl.voteBacklogMutex.Lock()
		cl.voteBacklog.PushBack(VoteInBacklog{
			vote:    vote,
			replica: replica,
		})
		cl.voteBacklogMutex.Unlock()
		fmt.Printf("adding vote to backlog sn=%d backlog.len=%d\n", vote.Sn, cl.voteBacklog.Len())
		return fmt.Errorf("Invalid sequence number: expected=%d, got=%d\n", cl.nextSeqNum[replica], vote.Sn)
	}

	// Update next sequence number.
	cl.nextSeqNum[replica]++

	// Verify timestamp.
	if vote.Ts <= cl.mostRecentTimestamp[replica] {
		// Rj sent old timestamp
		return fmt.Errorf("Invalid timestamp")
	}

	// Update most recent timestamp.
	cl.mostRecentTimestamp[replica] = vote.Ts

	if vote.GetIsHeartbeat() {
		// Heartbeat vote.
		return nil
	}

	// Update transaction timestamp.
	txid := vote.Tx.ID()

	// Upsert (txid, map(replica->timestamp)) if not already exist.
	if _, ok := cl.timestamps[txid]; !ok {
		cl.timestamps[vote.Tx.ID()] = make(map[*ReplicaInfo]timestamp)
	}

	// Check map(replica->timestamp) for pre-existing key.
	// If timestamps already contains this transaction and timestamp, discard it as duplicate
	if ts2, ok := cl.timestamps[txid][replica]; ok && ts2 != vote.Ts {
		// Duplicate vote.
		// TODO is paper correct here? This seems wrong.
		return fmt.Errorf("Duplicate vote")
	}

	// Store timestamp.
	cl.timestamps[txid][replica] = vote.Ts

	return nil
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
	votesCh := make(chan *pb.Vote, len(cl.replicas))
	votes := make([]*pb.Vote, 0)

	// Send transaction to all replicas.
	for _, replica := range cl.replicas {
		go func(replica *ReplicaInfo) {
			vote, err := replica.ReplicaServiceClient.Write(context.Background(), tx)
			if err != nil {
				// skip.
				fmt.Printf("Error writing to replica: %v", err)
				votesCh <- nil
				return
			}

			votesCh <- vote
		}(replica)
	}

	// Collect votes.
	for len(votes) < len(cl.replicas) {
		vote := <-votesCh
		votes = append(votes, vote)
	}

	// Log all votes for write.
	for _, vote := range votes {
		if vote == nil {
			continue
		}
		fmt.Printf("Vote: %v\n", vote)
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
	// CProp []*pb.Vote
}

func (cl *Client) Read() ReadResponse {
	var T []*pb.Transaction
	// var Cpp []*pb.Vote

	// 1. Compute rmin, rmax, rconf for each transaction.
	for tx, replicas := range cl.timestamps {
		// 1. rmin
		rMin := MinPossibleTimestamp(
			tx,
			cl.timestamps,
			cl.replicas,
			pAlpha, pBeta,
			cl.mostRecentTimestamp,
		)

		// 2. rmax
		rMax := MaxPossibleTimestamp(
			tx,
			cl.timestamps,
			cl.replicas,
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

	return ReadResponse{T, rPerf}
}

// For a timestamp ts, notation ts.getVoteMsg() denotes the vote message from some replica through which a client obtained timestamp ts. We abstract away the logic of how getVoteMsg() is implemented
func (cl *Client) getVoteMessage() {}
