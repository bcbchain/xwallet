package main

import (
	"fmt"
	"bcbchain.io/bcdb"
	"testing"
)

func TestLevelDB(t *testing.T) {

	db, err := bcdb.OpenDB("testdb", "127.0.0.1", "8888")
	if err != nil {
		fmt.Println(err)
	}

	key1 := []byte{0x01}

	db.Set([]byte("/genesis/token"), []byte("/genesis/token"))
	fmt.Println(key1, ":", db.Get([]byte("/genesis/token")))

}
