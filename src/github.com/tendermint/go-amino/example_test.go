package amino_test

import (
	"fmt"

	"github.com/tendermint/go-amino"
)

func Example() {

	defer func() {
		if e := recover(); e != nil {
			fmt.Println("Recovered:", e)
		}
	}()

	type Message interface{}

	type bcMessage struct {
		Message	string
		Height	int
	}

	type bcResponse struct {
		Status	int
		Message	string
	}

	type bcStatus struct {
		Peers int
	}

	var cdc = amino.NewCodec()
	cdc.RegisterInterface((*Message)(nil), nil)
	cdc.RegisterConcrete(&bcMessage{}, "bcMessage", nil)
	cdc.RegisterConcrete(&bcResponse{}, "bcResponse", nil)
	cdc.RegisterConcrete(&bcStatus{}, "bcStatus", nil)

	var bm = &bcMessage{Message: "ABC", Height: 100}
	var msg = bm

	var bz []byte
	var err error
	bz, err = cdc.MarshalBinary(msg)
	fmt.Printf("Encoded: %X (err: %v)\n", bz, err)

	var msg2 Message
	err = cdc.UnmarshalBinary(bz, &msg2)
	fmt.Printf("Decoded: %v (err: %v)\n", msg2, err)
	var bm2 = msg2.(*bcMessage)
	fmt.Printf("Decoded successfully: %v\n", *bm == *bm2)

}
