package pb

import (
	"crypto/sha256"
	"encoding/hex"
)

type TXID = string

func (tx *Transaction) ID() TXID {
	buf := sha256.Sum256(tx.GetData())
	return hex.EncodeToString(buf[:])
}
