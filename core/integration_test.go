package core

import (
	"fmt"
	"net"
	"strconv"
	"testing"
)

// // State:
// // All replicas and their public keys
// var R []Replica
// var pk map[Replica]PublicKey

// // The most recent timestamp returned by each replica
// var mrt map[Replica]int

// // The next sequence number expected by each replica
// var nextsn map[Replica]int

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

// // Types
// type Replica string
// type PublicKey string
// type Transaction string
// type Signature string
// type Pod struct {
// 	T      int
// 	r_perf int
// 	C_pp   []int
// }

// // Constants
// const HEARTBEAT Transaction = "HEARTBEAT"

// func TestFlow(t *testing.T) {
// 	// Example usage
// 	replicas := []Replica{"R1", "R2", "R3"}
// 	keys := []PublicKey{"Key1", "Key2", "Key3"}

// 	// Initialize event
// 	initEvent(replicas, keys)

// 	// Example vote reception
// 	receiveVote("tx1", 10, 0, "sig1", "R1")

// 	// Active validator set.
// 	// Whenever a validator receives a new transaction from a client, it appends it to its local log, together with the current timestamp based on its local clock. It then signs the transaction with the timestamp and hands it back to the client. The client receives the signed transaction and timestamp and validates the validator’s signature using its known public key. As soon as the client has collected a certain number of signatures from the validators (e.g., α = 2/3 of the validators), the client considers the transaction confirmed. The client associates the transaction with a timestamp, too: the median among the timestamps signed by the validators.
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

func makeReplica() (r Replica, address string) {
	replica1 := Replica{}
	replicaAddress := fmt.Sprintf("%s:%s", "localhost", getRandomPort())
	return replica1, replicaAddress
}

func TestFlow(t *testing.T) {
	replicas := []Replica{}
	replicaAddrs := []string{}

	for i := 0; i < 3; i++ {
		replica, replicaAddress := makeReplica()
		replicas = append(replicas, replica)
		replicaAddrs = append(replicaAddrs, replicaAddress)
		go replica.ListenAndServe(replicaAddress)
	}

	client1 := NewClient()
	for i := 0; i < 3; i++ {
		client1.connectToReplica(&replicaAddrs[i])
	}
	// client1.connectToReplica(&replicaAddress)
}
