package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSingleRuneUnicodeFlagParsing(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		shortOpts bool
		wantStr   string
		wantBool  bool
		wantArgs  []string
	}{
		{name: "single dash string", args: []string{"-א", "v"}, wantStr: "v", wantArgs: []string{}},
		{name: "single dash string with equals", args: []string{"-א=v"}, wantStr: "v", wantArgs: []string{}},
		{name: "double dash string", args: []string{"--א", "v"}, wantStr: "v", wantArgs: []string{}},
		{name: "single dash bool alias", args: []string{"-ש", "pos"}, wantBool: true, wantArgs: []string{"pos"}},
		{name: "combined short options bool then string", args: []string{"-שא", "v"}, shortOpts: true, wantStr: "v", wantBool: true, wantArgs: []string{}},
		{name: "combined short options mixed ascii", args: []string{"-bא", "v"}, shortOpts: true, wantStr: "v", wantBool: true, wantArgs: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotStr string
			var gotBool bool
			var gotArgs []string
			cmd := &Command{
				Name:                   "x",
				UseShortOptionHandling: tt.shortOpts,
				Flags: []Flag{
					&StringFlag{Name: "א"},
					&BoolFlag{Name: "b", Aliases: []string{"ש"}},
				},
				Action: func(_ context.Context, c *Command) error {
					gotStr = c.String("א")
					gotBool = c.Bool("b")
					gotArgs = c.Args().Slice()
					return nil
				},
			}
			require.NoError(t, cmd.Run(context.Background(), append([]string{"x"}, tt.args...)))
			assert.Equal(t, tt.wantStr, gotStr)
			assert.Equal(t, tt.wantBool, gotBool)
			assert.Equal(t, tt.wantArgs, gotArgs)
		})
	}
}

func TestSingleRuneUnicodeFlagMissingValue(t *testing.T) {
	cmd := &Command{
		Name:      "x",
		Writer:    &bytes.Buffer{},
		ErrWriter: &bytes.Buffer{},
		Flags:     []Flag{&StringFlag{Name: "א"}},
	}
	err := cmd.Run(context.Background(), []string{"x", "-א"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "flag needs an argument")
}

func TestSingleRuneUnicodeFlagHelpPrefix(t *testing.T) {
	var buf bytes.Buffer
	cmd := &Command{
		Name:   "x",
		Writer: &buf,
		Flags: []Flag{
			&StringFlag{Name: "א"},
			&BoolFlag{Name: "bee", Aliases: []string{"ש"}},
			&BoolWithInverseFlag{Name: "ת", Aliases: []string{"é"}},
		},
	}
	require.NoError(t, cmd.Run(context.Background(), []string{"x", "-h"}))
	out := buf.String()
	assert.Contains(t, out, "-א string")
	assert.NotContains(t, out, "--א")
	assert.Contains(t, out, "--bee, -ש")
	assert.NotContains(t, out, "--ש")
	assert.Contains(t, out, "-[no-]ת")
	assert.Contains(t, out, "-é")
	assert.NotContains(t, out, "--é")
}

func TestPrefixForSingleRune(t *testing.T) {
	assert.Equal(t, "-", prefixFor("a"))
	assert.Equal(t, "-", prefixFor("א"))
	assert.Equal(t, "-", prefixFor("界"))
	assert.Equal(t, "--", prefixFor("ab"))
	assert.Equal(t, "--", prefixFor("אב"))
}

func TestSuggestFlagSingleRune(t *testing.T) {
	cmd := &Command{Flags: []Flag{&StringFlag{Name: "א"}}}
	assert.Equal(t, "-א", suggestFlag(cmd.Flags, "א", true))
}
