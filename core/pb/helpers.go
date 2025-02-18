package pb

import (
	"crypto/sha256"
	"encoding/hex"
)

type TXID = string

func (tx *Transaction) ID() TXID {
	sha256 := sha256.New()
	sha256.Write(tx.GetCtx())
	buf := sha256.Sum(nil)
	// convert to hex
	return hex.EncodeToString(buf)
}

func (tx *Transaction) IsConfirmed() bool {
	// A transaction with confirmed round rconf is called confirmed if rconf ̸= ⊥, and unconfirmed otherwise.
	return tx.GetRConf() != 0
}
