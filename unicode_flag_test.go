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
		expected string
	}{
		{name: "single dash", args: []string{"app", "-ש", "value"}, expected: "value"},
		{name: "single dash with equals", args: []string{"app", "-ש=value"}, expected: "value"},
		{name: "double dash", args: []string{"app", "--ש", "value"}, expected: "value"},
		{name: "long name", args: []string{"app", "--שם", "value"}, expected: "value"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got string
			cmd := &Command{
				Name: "app",
				Flags: []Flag{
					&StringFlag{Name: "שם", Aliases: []string{"ש"}},
				},
				Action: func(_ context.Context, cmd *Command) error {
					got = cmd.String("שם")
					return nil
				},
			}
			require.NoError(t, cmd.Run(buildTestContext(t), tc.args))
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestSingleRuneUnicodeBoolFlagShortOptionHandling(t *testing.T) {
	var a, b bool
	var s string
	cmd := &Command{
		Name:                   "app",
		UseShortOptionHandling: true,
		Flags: []Flag{
			&BoolFlag{Name: "א"},
			&BoolFlag{Name: "ב"},
			&StringFlag{Name: "ג"},
		},
		Action: func(_ context.Context, cmd *Command) error {
			a, b, s = cmd.Bool("א"), cmd.Bool("ב"), cmd.String("ג")
			return nil
		},
	}
	require.NoError(t, cmd.Run(buildTestContext(t), []string{"app", "-אבג", "value"}))
	assert.True(t, a)
	assert.True(t, b)
	assert.Equal(t, "value", s)
}

func TestSingleRuneUnicodeFlagDisplay(t *testing.T) {
	assert.Equal(t, "-", prefixFor("ש"))
	assert.Equal(t, "--", prefixFor("שם"))

	out := &bytes.Buffer{}
	cmd := &Command{
		Name:   "app",
		Writer: out,
		Flags: []Flag{
			&StringFlag{Name: "ש", Usage: "hebrew"},
			&BoolWithInverseFlag{Name: "ב", Aliases: []string{"ג"}},
		},
	}
	require.NoError(t, cmd.Run(buildTestContext(t), []string{"app", "--help"}))
	help := out.String()
	assert.Contains(t, help, "-ש string")
	assert.False(t, strings.Contains(help, "--ש"), help)
	assert.Contains(t, help, "-[no-]ב, -ג")
}

func TestSingleRuneUnicodeFlagSuggestion(t *testing.T) {
	assert.Equal(t, "-ש", suggestFlag([]Flag{&StringFlag{Name: "ש"}}, "ש", false))
}
