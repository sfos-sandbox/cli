package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommand_DigitFlags(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		shortOpts     bool
		expectFour    bool
		expectSix     bool
		expectVerbose bool
		expectNum     int
		expectArgs    []string
	}{
		{name: "single dash digit", args: []string{"app", "-4", "operand"}, expectFour: true, expectArgs: []string{"operand"}},
		{name: "double dash digit", args: []string{"app", "--4", "operand"}, expectFour: true, expectArgs: []string{"operand"}},
		{name: "parsing continues", args: []string{"app", "-4", "-v", "operand"}, expectFour: true, expectVerbose: true, expectArgs: []string{"operand"}},
		{name: "operand between", args: []string{"app", "-4", "operand", "-6"}, expectFour: true, expectSix: true, expectArgs: []string{"operand"}},
		{name: "digit with equals", args: []string{"app", "-4=false", "-6=true"}, expectSix: true, expectArgs: []string{}},
		{name: "digit non-bool flag", args: []string{"app", "-1", "7", "x"}, expectNum: 7, expectArgs: []string{"x"}},
		{name: "short option group", args: []string{"app", "-46v", "x"}, shortOpts: true, expectFour: true, expectSix: true, expectVerbose: true, expectArgs: []string{"x"}},
		{name: "undefined negative number stays positional", args: []string{"app", "-5", "-v"}, expectArgs: []string{"-5", "-v"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gotArgs []string
			cmd := &Command{
				Name:                   "app",
				UseShortOptionHandling: tc.shortOpts,
				Flags: []Flag{
					&BoolFlag{Name: "4"},
					&BoolFlag{Name: "6"},
					&BoolFlag{Name: "verbose", Aliases: []string{"v"}},
					&IntFlag{Name: "1"},
				},
				Action: func(_ context.Context, cmd *Command) error {
					gotArgs = cmd.Args().Slice()
					return nil
				},
			}

			require.NoError(t, cmd.Run(buildTestContext(t), tc.args))
			assert.Equal(t, tc.expectFour, cmd.Bool("4"))
			assert.Equal(t, tc.expectSix, cmd.Bool("6"))
			assert.Equal(t, tc.expectVerbose, cmd.Bool("verbose"))
			assert.Equal(t, tc.expectNum, cmd.Int("1"))
			assert.Equal(t, tc.expectArgs, gotArgs)
		})
	}
}
