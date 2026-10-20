package base64x

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestFocus(t *testing.T) {
	cases := []string{
		"Zg==",      // valid "f"
		"Zg=",       // len 3
		"Zm9vYg==",  // valid "foob"
		"ABCDE",     // valid ABCD + extra E
		"Zg==A",     // valid + extra
		"123456",    // len 6 missing pad
		"AA=",       // len 3
		"Zm8=",      // valid "fo"
	}
	for _, c := range cases {
		in := []byte(c)
		// native (this binary)
		obuf := make([]byte, 64)
		n0, err0 := StdEncoding.Decode(obuf, in)
		// stdlib
		sbuf := make([]byte, 64)
		n1, err1 := base64.StdEncoding.Decode(sbuf, in)
		fmt.Printf("input=%-12q ours: n=%d out=%q err=%v | std: n=%d out=%q err=%v\n",
			c, n0, obuf[:n0], err0, n1, sbuf[:n1], err1)
	}
}
