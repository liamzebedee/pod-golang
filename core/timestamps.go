package core

import (
	"slices"
	"sort"

	"github.com/liamzebedee/pod-go/core/pb"
)

// The client uses the minPossibleTs() and maxPossibleTs() functions to compute rmin (line 9) and rmax (line 10)
// They return the minimum and maximum, respectively, round number that any client can ever confirm tx with.
//
// For rperf (line 22) the client uses minPossibleTsForNewTx(), which returns the smallest round number that any client can ever assign to a transaction not yet seen by the client.
//
// A transaction becomes confirmed (i.e., it gets assigned a confirmation round rconf ̸= ⊥) when the client receives α votes for tx from different replicas (line 12).
// Before it becomes confirmed, a transaction has rconf = ⊥ and Ctx = [ ] (line 11).
// When confirmed, rconf is the median of all received timestamps (line 18), and the auxiliary data Ctx contains all the received votes on tx (line 16).
//

// Returns the minimum timestamp bound (rmin) for a tx.
func MinPossibleTimestamp(tx pb.TXID, tsps map[pb.TXID]map[ReplicaID]timestamp, R []ReplicaID, alpha int, beta int, mrt map[ReplicaID]timestamp) timestamp {
	if _, ok := tsps[tx]; !ok {
		panic("tx not found in tsps")
	}

	// Aggregate all timestamps for a tx into a single list
	timestamps := []timestamp{}

	// 1. For each replica, get either their timestamp for the tx or the most recent timestamp from them
	for _, Rj := range R {
		if _, ok := tsps[tx][Rj]; ok {
			timestamps = append(timestamps, tsps[tx][Rj])
		} else {
			timestamps = append(timestamps, mrt[Rj])
		}
	}

	// 2. Sort the timestamps
	sort.Float64s(timestamps)

	// 3. Prepends β times the 0 value, pessimistically assuming that up to β replicas will try to bias tx by sending a timestamp 0 to other clients
	zeros := []timestamp{}
	for i := 0; i < beta; i++ {
		zeros = append(zeros, 0)
	}
	timestamps = slices.Concat(zeros, timestamps)

	// 4. Return the median of the first α timestamps
	return Median(timestamps[:alpha])
}

// Returns the maximum timestamp bound (rmax) for a tx.
func MaxPossibleTimestamp(tx pb.TXID, tsps map[pb.TXID]map[ReplicaID]timestamp, R []ReplicaID, alpha int, beta int) timestamp {
	if _, ok := tsps[tx]; !ok {
		panic("tx not found in tsps")
	}

	// Aggregate all timestamps for a tx into a single list
	timestamps := []timestamp{}

	// 1. For each replica, get either their timestamp for the tx or fill a missing vote with the ∞ value
	for _, Rj := range R {
		if _, ok := tsps[tx][Rj]; ok {
			timestamps = append(timestamps, tsps[tx][Rj])
		} else {
			timestamps = append(timestamps, 1<<31-1) // Using max int for ∞
		}
	}

	// 2. Sort the timestamps
	sort.Float64s(timestamps)

	// 3. Append β times the ∞ value (line 24), the worst-case timestamp that malicious replicas may send to other clients
	infs := []timestamp{}
	for i := 0; i < beta; i++ {
		infs = append(infs, 1<<31-1)
	}
	timestamps = slices.Concat(timestamps, infs)

	// 4. Return the median of the last α timestamps
	return Median(timestamps[len(timestamps)-alpha:])
}

func MinPossibleTimestampForNewTx(mrt []timestamp, alpha int, beta int) timestamp {
	// 1. Set timestamps to most recent timestamps
	timestamps := mrt

	// 2. Sort timestamps in ascending order
	sort.Float64s(timestamps)

	// 3. Prepend β times the 0 value
	zeros := []timestamp{}
	for i := 0; i < beta; i++ {
		zeros = append(zeros, 0)
	}
	timestamps = slices.Concat(zeros, timestamps)

	// 4. Return the median of the first α timestamps
	return Median(timestamps[:alpha])
}

func Median[T timestamp](Y []T) T {
	return Y[len(Y)/2]
}
