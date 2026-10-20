package base64x

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
)

func TestJSONFocus(t *testing.T) {
	cases := []string{
		"aGVsbG8=",          // plain valid
		`aGVsbG8\=`,         // escaped =
		`aGVsbG8\u003d`,     // unicode =
		`FPu\rcA9l+\n`,      // escaped r n
		`\/`,                // just escaped slash
		`aGVsbG8\\`,         // escaped backslash
		`aGVsbG8\"`,         // escaped quote
		"aGVsbG8\"",         // raw quote (invalid)
		"aGVsbG8\\x",        // invalid escape x
		`aGVsbG8\u0041`,     // unicode A
		`\u0041`,            // just unicode A
		`aGVsbG8\u`,         // truncated u escape
		`aGVsbG8\u00`,       // truncated
		`aGVsbG8\uZZZZ`,     // bad hex
		`aGVsbG8\u00e9`,     // non-ascii unicode
	}
	for _, c := range cases {
		in := []byte(c)
		obuf := make([]byte, 128)
		n0, err0 := JSONStdEncoding.Decode(obuf, in)
		// reference: json.Unmarshal the quoted string, then stdlib base64
		var ref string
		jerr := json.Unmarshal(append(append([]byte{'"'}, in...), '"'), &ref)
		var n1 int
		var err1 error
		var sbuf []byte
		if jerr == nil {
			sbuf = make([]byte, 128)
			n1, err1 = base64.StdEncoding.Decode(sbuf, []byte(ref))
		} else {
			err1 = jerr
		}
		fmt.Printf("input=%-24q ours: n=%d out=%q err=%v | ref: n=%d out=%q err=%v\n",
			c, n0, obuf[:n0], err0, n1, sbuf[:n1], err1)
	}
}
