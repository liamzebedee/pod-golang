package core

// Downstream consumers use the Client to interact with replicas
type Client struct{}

func (cl *Client) startup() {
	// 1. Send CONNECT to all Replicas.
	// At initialization the client also sends a ⟨CONNECT⟩ message to each replica, which initiates a streaming connection from the replica to the client.

	// A client maintains a connection to each replica and receives votes through ⟨VOTE (tx, ts, sn, σ, Rj )⟩ messages (lines 15–24).

	// When a vote is received from replica Rj , the client first verifies the signature σ under Rj ’s public key (line 16). If invalid, the vote is ignored. Then the client verifies that the vote contains the next sequence number it expects to receive from replica Rj (line 17). If this is not the case, the vote is backlogged and given again to the client at a later point (the backlogging functionality is not shown in the pseudocode)

	// Heartbeat messages. Clients update their most-recent timestamp mrt[Rj] every time they receive a vote from replica Rj (line 20 in Algorithm 1). However, in case Rj has not written any transaction in round r (simply because no client called write() in round r), Rj will advance its round number, but clients will not advance mrt[Rj]. We solve this by having replicas send a vote on a dummy heartBeat transaction the end of each round (lines 26–28). An obvious practi- cal implementation is to send heartBeat only for rounds when no other transactions were sent. When received by a client, a heartBeat is handled as a vote (i.e., it triggers line 15 in Algorithm 1), except that we do not check for duplicate heartBeat votes and do not update any associated values for the heartbeat transaction (see line 21 in Algorithm 1).

	// When clients read the pod, they obtain a pod data structure D = (T, rperf, Cpp), where T is set of transactions with their associated timestamps, rperf is a past- perfect round and Cpp is auxiliary data.
}
