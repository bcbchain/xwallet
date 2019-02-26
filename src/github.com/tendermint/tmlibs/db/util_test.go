package db

import (
	"fmt"
	"testing"
)

func TestPrefixIteratorNoMatchNil(t *testing.T) {
	for backend := range backends {
		t.Run(fmt.Sprintf("Prefix w/ backend %s", backend), func(t *testing.T) {
			db := newTempDB(t, backend)
			itr := IteratePrefix(db, []byte("2"))

			checkInvalid(t, itr)
		})
	}
}

func TestPrefixIteratorNoMatch1(t *testing.T) {
	for backend := range backends {
		t.Run(fmt.Sprintf("Prefix w/ backend %s", backend), func(t *testing.T) {
			db := newTempDB(t, backend)
			itr := IteratePrefix(db, []byte("2"))
			db.SetSync(bz("1"), bz("value_1"))

			checkInvalid(t, itr)
		})
	}
}

func TestPrefixIteratorNoMatch2(t *testing.T) {
	for backend := range backends {
		t.Run(fmt.Sprintf("Prefix w/ backend %s", backend), func(t *testing.T) {
			db := newTempDB(t, backend)
			db.SetSync(bz("3"), bz("value_3"))
			itr := IteratePrefix(db, []byte("4"))

			checkInvalid(t, itr)
		})
	}
}

func TestPrefixIteratorMatch1(t *testing.T) {
	for backend := range backends {
		t.Run(fmt.Sprintf("Prefix w/ backend %s", backend), func(t *testing.T) {
			db := newTempDB(t, backend)
			db.SetSync(bz("2"), bz("value_2"))
			itr := IteratePrefix(db, bz("2"))

			checkValid(t, itr, true)
			checkItem(t, itr, bz("2"), bz("value_2"))
			checkNext(t, itr, false)

			checkInvalid(t, itr)
		})
	}
}

func TestPrefixIteratorMatches1N(t *testing.T) {
	for backend := range backends {
		t.Run(fmt.Sprintf("Prefix w/ backend %s", backend), func(t *testing.T) {
			db := newTempDB(t, backend)

			db.SetSync(bz("a/1"), bz("value_1"))
			db.SetSync(bz("a/3"), bz("value_3"))

			db.SetSync(bz("b/3"), bz("value_3"))
			db.SetSync(bz("a-3"), bz("value_3"))
			db.SetSync(bz("a.3"), bz("value_3"))
			db.SetSync(bz("abcdefg"), bz("value_3"))
			itr := IteratePrefix(db, bz("a/"))

			checkValid(t, itr, true)
			checkItem(t, itr, bz("a/1"), bz("value_1"))
			checkNext(t, itr, true)
			checkItem(t, itr, bz("a/3"), bz("value_3"))

			checkNext(t, itr, false)

			checkInvalid(t, itr)
		})
	}
}
