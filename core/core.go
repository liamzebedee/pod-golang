package core

// When clients read the pod, they obtain a pod data structure D = (T, rperf, Cpp), where T is set of transactions with their associated timestamps, rperf is a past-perfect round and Cpp is auxiliary data
type ReadReponse struct {
	Txs   []Transaction
	RPerf uint
	Cpp   interface{}
}

// the reader obtains associated timestamps rmin, rmax,rconf and auxiliary data Ctx, which may evolve
type Transaction struct {
	// The minimum round.
	RMin float64
	// The maximum round.
	RMax float64
	// Undefined confirmed round.
	RConf float64
	Ctx   interface{}
}

func (tx *Transaction) IsConfirmed() bool {
	// A transaction with confirmed round rconf is called confirmed if rconf ̸= ⊥, and unconfirmed otherwise.
	return tx.RConf != 0
}

// A vote is a tuple (tx,ts,sn,σ,R), where tx is a trans- action, ts is a timestamp, sn is a sequence number, σ is a signature, and R is a replica. A vote is valid if σ is a valid signature on message m = (tx, ts, sn) with respect to the public key pkR of replica R.
type Vote struct {
	Tx  *Transaction
	Ts  float64
	Sn  int64
	Sig Signature // σ
	R   Replica
}
