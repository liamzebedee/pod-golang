package core

import (
	"crypto/rand"
	"encoding/hex"
)

// Mock signatures.
type Signature struct{}

func (s Signature) ToBytes() []byte { return []byte{} }

type PublicKey struct {
	string
}

func (pk PublicKey) String() string { return pk.string }

type SecretKey struct {
	string
}
type Keypair struct {
	P PublicKey
	S SecretKey
}

func genRandomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func Keygen() Keypair {
	kp := Keypair{
		P: PublicKey{
			genRandomString(32),
		},
		S: SecretKey{
			genRandomString(32),
		},
	}
	return kp
}

func Sign(sk SecretKey, msg interface{}) Signature             { return Signature{} }
func Verify(pk PublicKey, msg interface{}, sig Signature) bool { return true }
