package core

import (
	"fmt"
	"net"
	"strconv"
	"testing"
)

// // Timestamp received for each transaction from each replica
// var tsps map[Transaction]map[Replica]int

// // The pod observed by the client so far
// var D Pod

// // Event: Initialize replicas and their public keys
// func initEvent(replicas []Replica, keys []PublicKey) {
// 	R = replicas
// 	pk = make(map[Replica]PublicKey)

// 	// Assign public keys to replicas
// 	for i, r := range R {
// 		pk[r] = keys[i]
// 	}

// 	// Initialize state variables
// 	tsps = make(map[Transaction]map[Replica]int)
// 	D = Pod{0, 0, []int{}}

// 	for _, r := range R {
// 		mrt[r] = 0
// 		nextsn[r] = -1
// 		sendConnect(r)
// 	}
// }

// // Function: Send CONNECT message to a replica
// func sendConnect(r Replica) {
// 	fmt.Println("Send <CONNECT> to", r)
// }

// // Event: Receive a vote from a replica
// func receiveVote(tx Transaction, ts, sn int, sigma Signature, r Replica) {
// 	// Received a vote from replica Rj

// 	// If the vote signature is invalid, discard it
// 	if !verify(pk[r], tx, ts, sn, sigma) {
// 		return
// 	}

// 	// If the sequence number is not as expected, discard it
// 	if sn != nextsn[r] {
// 		return
// 	}

// 	// Increment the expected sequence number for Rj
// 	nextsn[r]++

// 	// If the timestamp is outdated, discard it
// 	if ts < mrt[r] {
// 		return
// 	}

// 	// Update the most recent timestamp for Rj
// 	mrt[r] = ts

// 	// If the transaction is a heartbeat, do not update tsps
// 	if tx == HEARTBEAT {
// 		return
// 	}

// 	// If tsps already contains this transaction and timestamp, discard it as duplicate
// 	if tsps[tx][r] != 0 && tsps[tx][r] != ts {
// 		return
// 	}

// 	// Store the timestamp for this transaction from Rj
// 	if tsps[tx] == nil {
// 		tsps[tx] = make(map[Replica]int)
// 	}
// 	tsps[tx][r] = ts
// }

// // Function: Verify a vote signature
// func verify(key PublicKey, tx Transaction, ts, sn int, sigma Signature) bool {
// 	// Placeholder verification logic
// 	return true
// }

func getRandomPort() string {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	portStr := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	fmt.Printf("got random port: %s\n", portStr)
	return portStr
}

func TestFlow(t *testing.T) {
	// Setup replicas config.
	replicas := make(map[PublicKey]*Replica)
	replicaConfigs := []ReplicaInfo{}
	N_REPLICAS := 5

	// Create replicas.
	for i := 0; i < N_REPLICAS; i++ {
		replica := NewReplica()
		addr := fmt.Sprintf("%s:%s", "localhost", getRandomPort())

		conf := ReplicaInfo{
			DialAddress: addr,
			PK:          replica.PublicKey(),
		}
		replicaConfigs = append(replicaConfigs, conf)

		replicas[replica.PublicKey()] = &replica

		go replica.ListenAndServe(conf.DialAddress)
	}

	// Create client.
	client := NewClient()

	// Start client with replica config.
	client.Start(replicaConfigs)

	// State:

	// The most recent timestamp returned by each replica
	// mostRecentTimestamp := make(map[*Replica]int)

	// // The next sequence number expected by each replica
	// nextSeqNum := make(map[*Replica]int)

	// // Timestamp -> Replica -> Sequence number
	// tsps := make(map[Transaction]map[Replica]int)
	// 	D = Pod{0, 0, []int{}}

	client.Write(makeTx(1))

	// Wait forever.
	ch := make(chan bool)
	<-ch
}
