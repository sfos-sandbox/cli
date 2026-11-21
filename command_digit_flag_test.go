package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSingleDashDigitFlags(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		shortOpts bool
		want4     bool
		want6     bool
		wantN     int
		wantArgs  []string
	}{
		{name: "single dash digit bool", args: []string{"app", "-4", "operand"}, want4: true, wantArgs: []string{"operand"}},
		{name: "double dash digit bool", args: []string{"app", "--4", "operand"}, want4: true, wantArgs: []string{"operand"}},
		{name: "later flags still parsed", args: []string{"app", "-4", "-6", "operand"}, want4: true, want6: true, wantArgs: []string{"operand"}},
		{name: "digit flag with value", args: []string{"app", "-1", "7", "x"}, wantN: 7, wantArgs: []string{"x"}},
		{name: "digit flag with equals", args: []string{"app", "-1=9", "x"}, wantN: 9, wantArgs: []string{"x"}},
		{name: "short option group of digits", args: []string{"app", "-46", "x"}, shortOpts: true, want4: true, want6: true, wantArgs: []string{"x"}},
		{name: "undefined digit stays positional", args: []string{"app", "-5", "-4"}, wantArgs: []string{"-5", "-4"}},
		{name: "partially defined group stays positional", args: []string{"app", "-42", "x"}, shortOpts: true, wantArgs: []string{"-42", "x"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotArgs []string
			cmd := &Command{
				Name:                   "app",
				UseShortOptionHandling: tt.shortOpts,
				Flags: []Flag{
					&BoolFlag{Name: "4"},
					&BoolFlag{Name: "6"},
					&IntFlag{Name: "1"},
				},
				Action: func(_ context.Context, cmd *Command) error {
					gotArgs = cmd.Args().Slice()
					assert.Equal(t, tt.want4, cmd.Bool("4"))
					assert.Equal(t, tt.want6, cmd.Bool("6"))
					assert.Equal(t, tt.wantN, cmd.Int("1"))
					return nil
				},
			}
			require.NoError(t, cmd.Run(buildTestContext(t), tt.args))
			assert.Equal(t, tt.wantArgs, gotArgs)
		})
	}
}

func TestSingleDashNegativeNumberDoesNotTriggerInvalidFlagAccessHandler(t *testing.T) {
	called := false
	cmd := &Command{
		Name:  "app",
		Flags: []Flag{&BoolFlag{Name: "v"}},
		InvalidFlagAccessHandler: func(context.Context, *Command, string) {
			called = true
		},
		Action: func(_ context.Context, cmd *Command) error {
			assert.Equal(t, []string{"-5"}, cmd.Args().Slice())
			return nil
		},
	}
	require.NoError(t, cmd.Run(buildTestContext(t), []string{"app", "-5"}))
	assert.False(t, called)
}
