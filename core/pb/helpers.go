package pb

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type TXID = string

// Returns a formatted string of the transaction, suitable for logging.
func (tx *Transaction) FormatString() string {
	rrange := tx.GetRMax() - tx.GetRMin()
	return fmt.Sprintf("Transaction{ID: %s, RMin: %f, RConf: %f, RMax: %f, Range: %f, Ctx: %v}", tx.ID(), tx.GetRMin(), tx.GetRConf(), tx.GetRMax(), rrange, tx.GetCtx())
}

func (tx *Transaction) ID() TXID {
	sha256 := sha256.New()
	sha256.Write(tx.GetCtx())
	buf := sha256.Sum(nil)
	// convert to hex
	return hex.EncodeToString(buf)
}

// A transaction with confirmed round rconf is called confirmed if rconf ̸= ⊥, and unconfirmed otherwise.
func (tx *Transaction) IsConfirmed() bool {
	return tx.GetRConf() != 0
}
