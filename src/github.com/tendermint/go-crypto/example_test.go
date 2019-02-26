package crypto_test

import (
	"fmt"

	"github.com/tendermint/go-crypto"
)

func ExampleSha256() {
	sum := crypto.Sha256([]byte("This is Tendermint"))
	fmt.Printf("%x\n", sum)

}

func ExampleRipemd160() {
	sum := crypto.Ripemd160([]byte("This is Tendermint"))
	fmt.Printf("%x\n", sum)

}
