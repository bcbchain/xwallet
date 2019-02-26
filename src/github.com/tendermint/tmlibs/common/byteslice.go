package common

import (
	"bytes"
)

func Fingerprint(slice []byte) []byte {
	fingerprint := make([]byte, 6)
	copy(fingerprint, slice)
	return fingerprint
}

func IsZeros(slice []byte) bool {
	for _, byt := range slice {
		if byt != byte(0) {
			return false
		}
	}
	return true
}

func RightPadBytes(slice []byte, l int) []byte {
	if l < len(slice) {
		return slice
	}
	padded := make([]byte, l)
	copy(padded[0:len(slice)], slice)
	return padded
}

func LeftPadBytes(slice []byte, l int) []byte {
	if l < len(slice) {
		return slice
	}
	padded := make([]byte, l)
	copy(padded[l-len(slice):], slice)
	return padded
}

func TrimmedString(b []byte) string {
	trimSet := string([]byte{0})
	return string(bytes.TrimLeft(b, trimSet))

}

func PrefixEndBytes(prefix []byte) []byte {
	if prefix == nil {
		return nil
	}

	end := make([]byte, len(prefix))
	copy(end, prefix)
	finished := false

	for !finished {
		if end[len(end)-1] != byte(255) {
			end[len(end)-1]++
			finished = true
		} else {
			end = end[:len(end)-1]
			if len(end) == 0 {
				end = nil
				finished = true
			}
		}
	}
	return end
}
