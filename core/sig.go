package core

// Mock signatures.
type Signature struct{}
type PublicKey struct{}
type SecretKey struct{}
type Keypair struct {
	P PublicKey
	S SecretKey
}

func Keygen() Keypair {
	return Keypair{
		P: PublicKey{},
		S: SecretKey{},
	}
}
func Sign(sk SecretKey, msg interface{}) Signature             { return Signature{} }
func Verify(pk PublicKey, msg interface{}, sig Signature) bool { return true }
