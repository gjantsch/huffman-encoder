package htree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConversionCodeToBytes(t *testing.T) {

	tests := []struct {
		name     string
		code     string
		expected []byte
	}{
		{
			name:     "all zeroes",
			code:     "000",
			expected: []byte{0},
		},
		{
			name:     "one",
			code:     "00000001",
			expected: []byte{1},
		},
		{
			name:     "two",
			code:     "00000010",
			expected: []byte{2},
		}, {
			name:     "sixty four",
			code:     "01000000",
			expected: []byte{0x40},
		}, {
			name:     "hundred twenty height",
			code:     "10000000",
			expected: []byte{0x80},
		},
		{
			name:     "padding two",
			code:     "0000001",
			expected: []byte{0x02},
		},
		{
			name:     "padding four",
			code:     "000001",
			expected: []byte{0x04},
		},
		{
			name:     "padding 128",
			code:     "1",
			expected: []byte{0x80},
		},

		{
			name:     "padding 0,128",
			code:     "000000001",
			expected: []byte{0, 0x80},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			code := Code{Symbol: 0, Code: tt.code}
			r.Equal(tt.expected, code.ToBytes())
		})
	}
}
