package common

import (
	"encoding/hex"
	"fmt"
	"strings"
)

var Fmt = func(format string, a ...interface{}) string {
	if len(a) == 0 {
		return format
	}
	return fmt.Sprintf(format, a...)
}

func IsHex(s string) bool {
	if len(s) > 2 && strings.EqualFold(s[:2], "0x") {
		_, err := hex.DecodeString(s[2:])
		return err == nil
	}
	return false
}

func StripHex(s string) string {
	if IsHex(s) {
		return s[2:]
	}
	return s
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func SplitAndTrim(s, sep, cutset string) []string {
	if s == "" {
		return []string{}
	}

	spl := strings.Split(s, sep)
	for i := 0; i < len(spl); i++ {
		spl[i] = strings.Trim(spl[i], cutset)
	}
	return spl
}

func IsASCIIText(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, b := range []byte(s) {
		if 32 <= b && b <= 126 {

		} else {
			return false
		}
	}
	return true
}

func ASCIITrim(s string) string {
	r := make([]byte, 0, len(s))
	for _, b := range []byte(s) {
		if b == 32 {
			continue
		} else if 32 < b && b <= 126 {
			r = append(r, b)
		} else {
			panic(fmt.Sprintf("non-ASCII (non-tab) char 0x%X", b))
		}
	}
	return string(r)
}
