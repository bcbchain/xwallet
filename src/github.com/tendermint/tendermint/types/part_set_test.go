package types

import (
	"bytes"
	"io/ioutil"
	"testing"

	cmn "github.com/tendermint/tmlibs/common"
)

const (
	testPartSize = 65536
)

func TestBasicPartSet(t *testing.T) {

	data := cmn.RandBytes(testPartSize * 100)

	partSet := NewPartSetFromData(data, testPartSize)
	if len(partSet.Hash()) == 0 {
		t.Error("Expected to get hash")
	}
	if partSet.Total() != 100 {
		t.Errorf("Expected to get 100 parts, but got %v", partSet.Total())
	}
	if !partSet.IsComplete() {
		t.Errorf("PartSet should be complete")
	}

	partSet2 := NewPartSetFromHeader(partSet.Header())

	for i := 0; i < partSet.Total(); i++ {
		part := partSet.GetPart(i)

		added, err := partSet2.AddPart(part, true)
		if !added || err != nil {
			t.Errorf("Failed to add part %v, error: %v", i, err)
		}
	}

	if !bytes.Equal(partSet.Hash(), partSet2.Hash()) {
		t.Error("Expected to get same hash")
	}
	if partSet2.Total() != 100 {
		t.Errorf("Expected to get 100 parts, but got %v", partSet2.Total())
	}
	if !partSet2.IsComplete() {
		t.Errorf("Reconstructed PartSet should be complete")
	}

	data2Reader := partSet2.GetReader()
	data2, err := ioutil.ReadAll(data2Reader)
	if err != nil {
		t.Errorf("Error reading data2Reader: %v", err)
	}
	if !bytes.Equal(data, data2) {
		t.Errorf("Got wrong data.")
	}

}

func TestWrongProof(t *testing.T) {

	data := cmn.RandBytes(testPartSize * 100)
	partSet := NewPartSetFromData(data, testPartSize)

	partSet2 := NewPartSetFromHeader(partSet.Header())

	part := partSet.GetPart(0)
	part.Proof.Aunts[0][0] += byte(0x01)
	added, err := partSet2.AddPart(part, true)
	if added || err == nil {
		t.Errorf("Expected to fail adding a part with bad trail.")
	}

	part = partSet.GetPart(1)
	part.Bytes[0] += byte(0x01)
	added, err = partSet2.AddPart(part, true)
	if added || err == nil {
		t.Errorf("Expected to fail adding a part with bad bytes.")
	}

}
