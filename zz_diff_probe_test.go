package base64x

import (
	"encoding/base64"
	"testing"
)

var diffCases [][]byte

func init() {
	base := []string{
		"", "A", "AB", "ABC", "ABCD", "ABCDE", "ABCDEF", "ABCDEFG", "ABCDEFGH",
		"Zm8=", "Zm8", "Zg==", "Zg=", "Zg", "Zm9v",
		"=", "==", "===", "=A", "A=", "A===", "AA=", "AAA=", "AAAA=",
		"Zg==A", "Zg=A", "Z=g=",
		" ABC", "ABC ", " AB ", "A BC",
		"AB\r\nCD", "AB\nCD", "AB\rCD", "\r\nAB", "AB\r\n",
		"Z\ng\n8\n=", "Zg\r\n==",
		"!A", "A!", "A!B", "A@#$",
		"Zm8=", "Zm9vYg==", "Zm9vYmE=", "Zm9vYmFy",
		// raw mode without padding
		"Zm8", "Zg", "Zm9v",
		// trailing whitespace after padding
		"Zg==\r\n", "Zg== ", "Zg==\n\r",
		// embedded invalid after valid
		"AB=C", "AB==C", "AB=C=",
	}
	for _, s := range base {
		diffCases = append(diffCases, []byte(s))
	}
}

func diffPair(t *testing.T, ours Encoding, std *base64.Encoding, input []byte) {
	// ours
	obuf := make([]byte, ours.DecodedLen(len(input))+4)
	n0, err0 := ours.Decode(obuf, input)
	// stdlib
	sbuf := make([]byte, std.DecodedLen(len(input))+4)
	n1, err1 := std.Decode(sbuf, input)

	b0, b1 := obuf[:n0], sbuf[:n1]
	errMatch := (err0 != nil) == (err1 != nil)
	if string(b0) != string(b1) || !errMatch || n0 != n1 {
		t.Errorf("DIVERGE input=%q\n  ours: n=%d out=%q err=%v\n  std : n=%d out=%q err=%v",
			input, n0, b0, err0, n1, b1, err1)
	}
}

func TestDiffAll(t *testing.T) {
	pairs := []struct {
		ours Encoding
		std  *base64.Encoding
	}{
		{StdEncoding, base64.StdEncoding},
		{URLEncoding, base64.URLEncoding},
		{RawStdEncoding, base64.RawStdEncoding},
		{RawURLEncoding, base64.RawURLEncoding},
	}
	for _, p := range pairs {
		for _, c := range diffCases {
			diffPair(t, p.ours, p.std, c)
		}
	}
}

func TestDiffLengths(t *testing.T) {
	// feed valid base64 of every length 0..64, both std and raw
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	for n := 0; n <= 64; n++ {
		input := []byte(alphabet[:n])
		diffPair(t, StdEncoding, base64.StdEncoding, input)
		diffPair(t, RawStdEncoding, base64.RawStdEncoding, input)
	}
}

func TestDiffWhitespacePositions(t *testing.T) {
	// base valid 48-char string, insert \r \n \t space at every position
	valid := "Tm93IHRoYXQgdGhlIHdvcmRzIGFyZSB3cm9uZyBkb25l"
	for pos := 0; pos <= len(valid); pos++ {
		for _, ws := range []string{"\r", "\n", "\r\n", " ", "\t"} {
			input := []byte(valid[:pos] + ws + valid[pos:])
			diffPair(t, StdEncoding, base64.StdEncoding, input)
		}
	}
}

func TestDiffPrintableInvalid(t *testing.T) {
	// for each printable non-base64 char, inject at position 0, middle, end
	base := "Zm9vYmFy"
	for ch := 32; ch < 127; ch++ {
		c := byte(ch)
		// skip valid base64 alphabet
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '+' || c == '/' || c == '=' {
			continue
		}
		for _, pos := range []int{0, 3, len(base)} {
			input := []byte(base[:pos] + string(c) + base[pos:])
			diffPair(t, StdEncoding, base64.StdEncoding, input)
		}
	}
}
