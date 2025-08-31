package main

import (
	"fmt"
	"io"
	"math/rand/v2"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/sjiekak/logen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Logen(t *testing.T) {
	st, err := logen.NewSanitizer()
	require.NoError(t, err)
	assert.NotNil(t, st)
}

func Test_Logmatch(t *testing.T) {
	for _, tc := range []struct {
		name           string
		carrots        int
		sellers        int
		weird          int
		expectedGroups int
	}{
		{
			name: "Empty",
		},
		{
			name:           "Carrots",
			carrots:        100,
			expectedGroups: 1,
		},
		{
			name:           "Sellers",
			sellers:        100,
			expectedGroups: 1,
		},
		{
			name:           "Weirds",
			weird:          1,
			expectedGroups: 1,
		},
		{
			name:           "All",
			carrots:        100,
			sellers:        100,
			weird:          1,
			expectedGroups: 3,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := makeText(tc.carrots, tc.sellers, tc.weird)
			state, err := logmatch(r)
			require.NoError(t, err)

			assert.Len(t, state.classes, tc.expectedGroups)

		})
	}
}

func Benchmark_Logmatch(b *testing.B) {
	text := makeText(100, 100, 100)

	for b.Loop() {
		logmatch(text)
	}
}

func patternCarrotAndStick() string {
	return fmt.Sprintf("I have %d carrots and %d sticks", rand.Uint32(), rand.Uint32())
}

func patternSeller() string {
	return fmt.Sprintf("%s is a unique seller", uuid.NewString())
}

func patternWeirdString() string { // currently fail matching
	return fmt.Sprintf("This %s is a complex case %d [id/%s]", randomAlphanumeric(24), rand.Int32(), uuid.NewString())
}

func randomAlphanumeric(length int) string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	var b strings.Builder

	for i := 0; i < length; i++ {
		b.WriteByte(chars[rand.IntN(len(chars))])
	}
	return b.String()
}

func makeText(carrotCount, sellerCount, weirdCount int) io.Reader {
	lines := make([]string, carrotCount+sellerCount+weirdCount)

	for i := range lines {
		var text string
		switch {
		case i < carrotCount:
			text = patternCarrotAndStick()
		case i < carrotCount+sellerCount:
			text = patternSeller()
		case i < carrotCount+sellerCount+weirdCount:
			text = patternWeirdString()
		}
		lines[i] = text
	}

	rand.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })
	return strings.NewReader(strings.Join(lines, "\n"))
}
