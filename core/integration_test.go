package core

import "testing"

func TestFlow(t *testing.T) {
	// Active validator set.
	// Whenever a validator receives a new transaction from a client, it appends it to its local log, together with the current timestamp based on its local clock. It then signs the transaction with the timestamp and hands it back to the client. The client receives the signed transaction and timestamp and validates the validator’s signature using its known public key. As soon as the client has collected a certain number of signatures from the validators (e.g., α = 2/3 of the validators), the client considers the transaction confirmed. The client associates the transaction with a timestamp, too: the median among the timestamps signed by the validators.
}
