package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDigitFlagSingleDash(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		short     bool
		want4     bool
		want6     bool
		wantCount int64
		wantArgs  []string
	}{
		{name: "single dash", args: []string{"app", "-4", "operand"}, want4: true, wantArgs: []string{"operand"}},
		{name: "double dash", args: []string{"app", "--4", "operand"}, want4: true, wantArgs: []string{"operand"}},
		{name: "later flags parsed", args: []string{"app", "-4", "-6", "operand"}, want4: true, want6: true, wantArgs: []string{"operand"}},
		{name: "with value", args: []string{"app", "-1", "3", "-4", "x"}, want4: true, wantCount: 3, wantArgs: []string{"x"}},
		{name: "with equals", args: []string{"app", "-1=5", "x"}, wantCount: 5, wantArgs: []string{"x"}},
		{name: "negative number stays positional", args: []string{"app", "-5", "-4"}, wantArgs: []string{"-5", "-4"}},
		{name: "short group", args: []string{"app", "-46", "x"}, short: true, want4: true, want6: true, wantArgs: []string{"x"}},
		{name: "unknown short group stays positional", args: []string{"app", "-45", "x"}, short: true, wantArgs: []string{"-45", "x"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotArgs []string
			cmd := &Command{
				Name:                   "app",
				UseShortOptionHandling: tt.short,
				Flags: []Flag{
					&BoolFlag{Name: "4"},
					&BoolFlag{Name: "6"},
					&Int64Flag{Name: "1"},
				},
				Action: func(_ context.Context, cmd *Command) error {
					gotArgs = cmd.Args().Slice()
					return nil
				},
			}
			require.NoError(t, cmd.Run(buildTestContext(t), tt.args))
			assert.Equal(t, tt.want4, cmd.Bool("4"))
			assert.Equal(t, tt.want6, cmd.Bool("6"))
			assert.Equal(t, tt.wantCount, cmd.Int64("1"))
			assert.Equal(t, tt.wantArgs, gotArgs)
		})
	}
}

func TestDigitFlagSingleDash_NoInvalidFlagAccess(t *testing.T) {
	for _, shortOptions := range []bool{false, true} {
		var invalid []string
		cmd := &Command{
			UseShortOptionHandling: shortOptions,
			Flags:                  []Flag{&BoolFlag{Name: "4"}},
			InvalidFlagAccessHandler: func(_ context.Context, _ *Command, name string) {
				invalid = append(invalid, name)
			},
			Action: func(context.Context, *Command) error { return nil },
		}
		require.NoError(t, cmd.Run(buildTestContext(t), []string{"app", "-4", "-5", "-45"}))
		assert.Empty(t, invalid, "shortOptions=%v", shortOptions)
	}
}
