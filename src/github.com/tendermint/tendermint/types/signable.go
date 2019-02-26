package types

type Signable interface {
	SignBytes(chainID string) []byte
}
