package db

import (
	"bytes"
)

func cp(bz []byte) (ret []byte) {
	ret = make([]byte, len(bz))
	copy(ret, bz)
	return ret
}

func cpIncr(bz []byte) (ret []byte) {
	if len(bz) == 0 {
		panic("cpIncr expects non-zero bz length")
	}
	ret = cp(bz)
	for i := len(bz) - 1; i >= 0; i-- {
		if ret[i] < byte(0xFF) {
			ret[i]++
			return
		}
		ret[i] = byte(0x00)
		if i == 0 {

			return nil
		}
	}
	return nil
}

func cpDecr(bz []byte) (ret []byte) {
	if len(bz) == 0 {
		panic("cpDecr expects non-zero bz length")
	}
	ret = cp(bz)
	for i := len(bz) - 1; i >= 0; i-- {
		if ret[i] > byte(0x00) {
			ret[i]--
			return
		}
		ret[i] = byte(0xFF)
		if i == 0 {

			return nil
		}
	}
	return nil
}

func IsKeyInDomain(key, start, end []byte, isReverse bool) bool {
	if !isReverse {
		if bytes.Compare(key, start) < 0 {
			return false
		}
		if end != nil && bytes.Compare(end, key) <= 0 {
			return false
		}
		return true
	} else {
		if start != nil && bytes.Compare(start, key) < 0 {
			return false
		}
		if end != nil && bytes.Compare(key, end) <= 0 {
			return false
		}
		return true
	}
}
