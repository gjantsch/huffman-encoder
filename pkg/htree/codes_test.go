package htree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCanonicalCodes(t *testing.T) {
	r := require.New(t)

	// Input uses arbitrary tree-traversal codes; only the lengths matter.
	// Sorted by (len, symbol): A(len=1), D(len=2), B(len=3), C(len=3)
	// Expected canonical assignment: A=0, D=10, B=110, C=111
	input := map[byte]Code{
		'A': {Symbol: 'A', Code: "0"},
		'D': {Symbol: 'D', Code: "11"},
		'B': {Symbol: 'B', Code: "100"},
		'C': {Symbol: 'C', Code: "101"},
	}

	got := CanonicalCodes(input)

	// Lengths must be preserved
	r.Equal(len(input['A'].Code), len(got['A'].Code))
	r.Equal(len(input['D'].Code), len(got['D'].Code))
	r.Equal(len(input['B'].Code), len(got['B'].Code))
	r.Equal(len(input['C'].Code), len(got['C'].Code))

	// Exact canonical values
	r.Equal("0", got['A'].Code)
	r.Equal("10", got['D'].Code)
	r.Equal("110", got['B'].Code)
	r.Equal("111", got['C'].Code)

	// Codes must be prefix-free
	codes := make([]string, 0, len(got))
	for _, c := range got {
		codes = append(codes, c.Code)
	}
	for i, a := range codes {
		for j, b := range codes {
			if i != j {
				r.NotEqual(a, b[:min(len(a), len(b))], "code %q is a prefix of %q", a, b)
			}
		}
	}

	// Output is deterministic regardless of map iteration order
	got2 := CanonicalCodes(input)
	for sym, c := range got {
		r.Equal(c.Code, got2[sym].Code)
	}
}
