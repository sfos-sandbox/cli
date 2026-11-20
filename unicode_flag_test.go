package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSingleRuneUnicodeFlagParsing(t *testing.T) {
	for _, tc := range []struct {
		name  string
		args  []string
		short bool
	}{
		{name: "single dash", args: []string{"app", "-א", "val"}},
		{name: "single dash equals", args: []string{"app", "-א=val"}},
		{name: "double dash", args: []string{"app", "--א", "val"}},
		{name: "short option handling", args: []string{"app", "-א", "val"}, short: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got string
			cmd := &Command{
				Name:                   "app",
				UseShortOptionHandling: tc.short,
				Flags:                  []Flag{&StringFlag{Name: "א"}},
				Action: func(_ context.Context, c *Command) error {
					got = c.String("א")
					return nil
				},
			}
			require.NoError(t, cmd.Run(buildTestContext(t), tc.args))
			assert.Equal(t, "val", got)
		})
	}
}

func TestSingleRuneUnicodeShortOptionCombination(t *testing.T) {
	var a, b bool
	var s string
	cmd := &Command{
		Name:                   "app",
		UseShortOptionHandling: true,
		Flags: []Flag{
			&BoolFlag{Name: "alpha", Aliases: []string{"א"}},
			&BoolFlag{Name: "bet", Aliases: []string{"ב"}},
			&StringFlag{Name: "gimel", Aliases: []string{"ג"}},
		},
		Action: func(_ context.Context, c *Command) error {
			a, b, s = c.Bool("alpha"), c.Bool("bet"), c.String("gimel")
			return nil
		},
	}
	require.NoError(t, cmd.Run(buildTestContext(t), []string{"app", "-אבג", "x"}))
	assert.True(t, a)
	assert.True(t, b)
	assert.Equal(t, "x", s)
}

func TestSingleRuneUnicodeFlagDisplay(t *testing.T) {
	assert.Equal(t, "-א", prefixFor("א")+"א")
	assert.Equal(t, "--אב", prefixFor("אב")+"אב")

	f := &StringFlag{Name: "א", Usage: "hebrew"}
	assert.True(t, strings.HasPrefix(f.String(), "-א "), f.String())

	bif := &BoolWithInverseFlag{Name: "א", Aliases: []string{"ב"}}
	assert.True(t, strings.HasPrefix(bif.String(), "-[no-]א, -ב"), bif.String())

	buf := &bytes.Buffer{}
	printFlagSuggestions("-", []Flag{f}, buf)
	assert.Contains(t, buf.String(), "-א:hebrew")
	assert.NotContains(t, buf.String(), "--א")
}
