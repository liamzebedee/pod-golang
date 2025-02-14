package core

type Replica struct {
	PublicKey PublicKey
	TxLog     []Transaction
}

// Each replica maintains a sequence number, which it increments and includes every time it assigns a timestamp to a transaction

// We assume that all messages between clients and replicas are concatenated with a session identifier (sid), which is unique for each concurrent execution of the protocol. Moreover, the sid is im- plicitly included in all messages signed by the replicas
// For a timestamp ts, notation ts.getVoteMsg() denotes the vote message from some replica through which a client obtained timestamp ts. We abstract away the logic of how getVoteMsg() is implemented

func (rep *Replica) run() {
	// When it receives ⟨WRITE tx⟩. a replica first checks whether it has already seen tx, in which case the message is ignored. Otherwise, it assigns tx a timestamp ts equal its local round number and the next available sequence number sn, and signs the message (tx, ts, sn) (line 19). Honest replicas use incremental sequence numbers for each transaction, so if a vote has a larger sequence number than another vote, it will have a larger or equal timestamp than the other. The replica appends (tx,ts,sn,σ) to replicaLog, and sends it via a ⟨VOTE (tx,ts,sn,σ,R)⟩ message to all connected clients (line 22
}
