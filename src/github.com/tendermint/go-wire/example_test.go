package wire_test

import (
	"bytes"
	"fmt"
	"log"

	"github.com/tendermint/go-wire"
)

func Example_RegisterInterface() {
	type Receiver interface{}
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

	var _ = wire.RegisterInterface(
		struct{ Receiver }{},
		wire.ConcreteType{&bcMessage{}, 0x01},
		wire.ConcreteType{&bcResponse{}, 0x02},
		wire.ConcreteType{&bcStatus{}, 0x03},
	)
}

func Example_EndToEnd_ReadWriteBinary() {
	type Receiver interface{}
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

	var _ = wire.RegisterInterface(
		struct{ Receiver }{},
		wire.ConcreteType{&bcMessage{}, 0x01},
		wire.ConcreteType{&bcResponse{}, 0x02},
		wire.ConcreteType{&bcStatus{}, 0x03},
	)

	var n int
	var err error
	buf := new(bytes.Buffer)
	bm := &bcMessage{Message: "Tendermint", Height: 100}
	wire.WriteBinary(bm, buf, &n, &err)
	if err != nil {
		log.Fatalf("writeBinary: %v", err)
	}
	fmt.Printf("Encoded: %x\n", buf.Bytes())

	recv := wire.ReadBinary(struct{ Receiver }{}, buf, 0, &n, &err).(struct{ Receiver }).Receiver
	if err != nil {
		log.Fatalf("readBinary: %v", err)
	}
	decoded := recv.(*bcMessage)
	fmt.Printf("Decoded: %#v\n", decoded)

}
