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
	tests := []struct {
		name     string
		args     []string
		wantStr  string
		wantBool bool
		wantArgs []string
	}{
		{name: "string flag with separate value", args: []string{"app", "-ש", "shalom"}, wantStr: "shalom"},
		{name: "string flag with equals", args: []string{"app", "-ש=shalom"}, wantStr: "shalom"},
		{name: "bool flag", args: []string{"app", "-ב"}, wantBool: true},
		{name: "ascii alias still works", args: []string{"app", "-s", "x", "-b"}, wantStr: "x", wantBool: true},
		{name: "positional after flag", args: []string{"app", "-ש", "v", "pos"}, wantStr: "v", wantArgs: []string{"pos"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gotStr string
			var gotBool bool
			var gotArgs []string
			cmd := &Command{
				Name: "app",
				Flags: []Flag{
					&StringFlag{Name: "ש", Aliases: []string{"s"}},
					&BoolFlag{Name: "ב", Aliases: []string{"b"}},
				},
				Action: func(_ context.Context, c *Command) error {
					gotStr = c.String("ש")
					gotBool = c.Bool("ב")
					gotArgs = c.Args().Slice()
					return nil
				},
			}
			require.NoError(t, cmd.Run(buildTestContext(t), tc.args))
			assert.Equal(t, tc.wantStr, gotStr)
			assert.Equal(t, tc.wantBool, gotBool)
			if tc.wantArgs == nil {
				assert.Empty(t, gotArgs)
			} else {
				assert.Equal(t, tc.wantArgs, gotArgs)
			}
		})
	}
}

func TestSingleRuneUnicodeFlagShortOptionHandling(t *testing.T) {
	var a, b bool
	var s string
	cmd := &Command{
		Name:                   "app",
		UseShortOptionHandling: true,
		Flags: []Flag{
			&BoolFlag{Name: "א"},
			&BoolFlag{Name: "ב"},
			&StringFlag{Name: "ש"},
		},
		Action: func(_ context.Context, c *Command) error {
			a, b, s = c.Bool("א"), c.Bool("ב"), c.String("ש")
			return nil
		},
	}
	require.NoError(t, cmd.Run(buildTestContext(t), []string{"app", "-אבש", "val"}))
	assert.True(t, a)
	assert.True(t, b)
	assert.Equal(t, "val", s)
}

func TestSingleRuneUnicodeFlagHelpAndSuggestion(t *testing.T) {
	assert.Equal(t, "-", prefixFor("ש"))
	assert.Equal(t, "-", prefixFor("s"))
	assert.Equal(t, "--", prefixFor("שש"))

	f := &StringFlag{Name: "ש", Usage: "hebrew"}
	assert.True(t, strings.HasPrefix(f.String(), "-ש "), f.String())

	bif := &BoolWithInverseFlag{Name: "ש", Aliases: []string{"ב"}}
	assert.True(t, strings.HasPrefix(bif.String(), "-[no-]ש, -ב"), bif.String())

	assert.Equal(t, "-ש", suggestFlag([]Flag{&StringFlag{Name: "ש"}}, "ש", true))

	out := &bytes.Buffer{}
	printFlagSuggestions("-", []Flag{&StringFlag{Name: "ש"}}, out)
	assert.Contains(t, out.String(), "-ש\n")
	assert.NotContains(t, out.String(), "--ש")
}
