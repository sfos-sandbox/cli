package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Flags whose names start with a digit (like ssh -4/-6 or ls -1) must be
// parsed when given with a single dash, while undefined digit arguments such
// as negative numbers keep being treated as positional arguments.
func TestCommand_DigitNamedSingleDashFlags(t *testing.T) {
	tests := []struct {
		name         string
		flags        []Flag
		useShortOpts bool
		args         []string
		expectSet    map[string]bool
		expectInt    map[string]int
		expectArgs   []string
	}{
		{
			name:       "single dash bool digit flag",
			flags:      []Flag{&BoolFlag{Name: "4"}},
			args:       []string{"app", "-4", "operand"},
			expectSet:  map[string]bool{"4": true},
			expectArgs: []string{"operand"},
		},
		{
			name:       "double dash bool digit flag",
			flags:      []Flag{&BoolFlag{Name: "4"}},
			args:       []string{"app", "--4", "operand"},
			expectSet:  map[string]bool{"4": true},
			expectArgs: []string{"operand"},
		},
		{
			name:       "digit flag followed by later flags",
			flags:      []Flag{&BoolFlag{Name: "6"}, &BoolFlag{Name: "v"}},
			args:       []string{"app", "-6", "-v", "operand"},
			expectSet:  map[string]bool{"6": true, "v": true},
			expectArgs: []string{"operand"},
		},
		{
			name:       "digit alias",
			flags:      []Flag{&BoolFlag{Name: "one", Aliases: []string{"1"}}},
			args:       []string{"app", "-1", "operand"},
			expectSet:  map[string]bool{"one": true},
			expectArgs: []string{"operand"},
		},
		{
			name:       "digit flag with value",
			flags:      []Flag{&IntFlag{Name: "4"}},
			args:       []string{"app", "-4", "10", "operand"},
			expectInt:  map[string]int{"4": 10},
			expectArgs: []string{"operand"},
		},
		{
			name:       "digit flag with equals value",
			flags:      []Flag{&IntFlag{Name: "4"}},
			args:       []string{"app", "-4=10", "operand"},
			expectInt:  map[string]int{"4": 10},
			expectArgs: []string{"operand"},
		},
		{
			name:         "short option group with digit flags",
			flags:        []Flag{&BoolFlag{Name: "1"}, &BoolFlag{Name: "a"}},
			useShortOpts: true,
			args:         []string{"app", "-1a", "operand"},
			expectSet:    map[string]bool{"1": true, "a": true},
			expectArgs:   []string{"operand"},
		},
		{
			name:       "undefined digit flag is positional and stops parsing",
			flags:      []Flag{&BoolFlag{Name: "v"}},
			args:       []string{"app", "-5", "-v", "operand"},
			expectSet:  map[string]bool{"v": false},
			expectArgs: []string{"-5", "-v", "operand"},
		},
		{
			name:       "negative number not matching a defined digit flag",
			flags:      []Flag{&BoolFlag{Name: "4"}, &BoolFlag{Name: "v"}},
			args:       []string{"app", "-42", "-v"},
			expectSet:  map[string]bool{"4": false, "v": false},
			expectArgs: []string{"-42", "-v"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gotArgs []string
			cmd := &Command{
				Name:                   "app",
				Flags:                  tc.flags,
				UseShortOptionHandling: tc.useShortOpts,
				Action: func(_ context.Context, c *Command) error {
					gotArgs = c.Args().Slice()
					for k, v := range tc.expectSet {
						assert.Equal(t, v, c.Bool(k), "flag %q", k)
					}
					for k, v := range tc.expectInt {
						assert.Equal(t, v, c.Int(k), "flag %q", k)
					}
					return nil
				},
			}
			require.NoError(t, cmd.Run(buildTestContext(t), tc.args))
			assert.Equal(t, tc.expectArgs, gotArgs)
		})
	}
}
