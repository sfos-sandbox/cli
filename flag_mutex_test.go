package cli

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newCommand() *Command {
	return &Command{
		MutuallyExclusiveFlags: []MutuallyExclusiveFlags{
			{
				Flags: [][]Flag{
					{
						&Int64Flag{
							Name: "i",
						},
						&StringFlag{
							Name:    "s",
							Sources: EnvVars("S_VAR"),
						},
						&BoolWithInverseFlag{
							Name: "b",
						},
					},
					{
						&Int64Flag{
							Name:    "t",
							Aliases: []string{"ai"},
							Sources: EnvVars("T_VAR"),
						},
						&BoolFlag{
							Name: "q",
						},
					},
				},
			},
		},
	}
}

func TestFlagMutuallyExclusiveFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		errStr   string
		required bool
		envs     map[string]string
	}{
		{
			name: "simple",
		},
		{
			name: "set one flag",
			args: []string{"--i", "10"},
		},
		{
			name:   "set both flags",
			args:   []string{"--i", "11", "--ai", "12"},
			errStr: "option i cannot be set along with option t",
		},
		{
			name:   "set both flags second member",
			args:   []string{"--i", "11", "--q"},
			errStr: "option i cannot be set along with option q",
		},
		{
			name:     "required none set",
			required: true,
			errStr:   "one of these flags needs to be provided",
		},
		{
			name:     "required one set",
			args:     []string{"--i", "10"},
			required: true,
		},
		{
			name:     "required both set",
			args:     []string{"--i", "11", "--ai", "12"},
			errStr:   "option i cannot be set along with option t",
			required: true,
		},
		{
			name:     "required both set second member",
			args:     []string{"--i", "11", "--q"},
			errStr:   "option i cannot be set along with option q",
			required: true,
		},
		{
			name:     "required both set second member both groups",
			args:     []string{"--s", "value", "--q"},
			errStr:   "option s cannot be set along with option q",
			required: true,
		},
		{
			name:     "set env var",
			required: true,
			envs: map[string]string{
				"S_VAR": "some",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.envs != nil {
				for k, v := range test.envs {
					t.Setenv(k, v)
				}
			}
			cmd := newCommand()
			cmd.MutuallyExclusiveFlags[0].Required = test.required

			err := cmd.Run(buildTestContext(t), append([]string{"foo"}, test.args...))
			if test.errStr == "" {
				assert.NoError(t, err)
				return
			}
			if err == nil {
				t.Error("Expected mutual exclusion error")
				return
			}

			switch err.(type) {
			case (*mutuallyExclusiveGroup), (*mutuallyExclusiveGroupRequiredFlag):
				assert.Contains(t, err.Error(), test.errStr)
			default:
				t.Errorf("got invalid error type %T", err)
			}
		})
	}
}
func TestFlagMutuallyExclusiveFlagsInheritedBySubcommand(t *testing.T) {
	newRoot := func(actionRan *string) *Command {
		return &Command{
			Name: "root",
			MutuallyExclusiveFlags: []MutuallyExclusiveFlags{
				{
					Required: true,
					Flags: [][]Flag{
						{&StringFlag{Name: "alpha"}},
						{&StringFlag{Name: "beta"}},
					},
				},
			},
			Commands: []*Command{
				{
					Name: "sub",
					Action: func(_ context.Context, cmd *Command) error {
						*actionRan = cmd.String("alpha") + cmd.String("beta")
						return nil
					},
					Commands: []*Command{
						{
							Name: "leaf",
							Action: func(_ context.Context, cmd *Command) error {
								*actionRan = cmd.String("alpha") + cmd.String("beta")
								return nil
							},
						},
					},
				},
			},
		}
	}

	tests := []struct {
		name   string
		args   []string
		errStr string
		want   string
	}{
		{name: "flag before subcommand", args: []string{"root", "--alpha", "hello", "sub"}, want: "hello"},
		{name: "flag after subcommand", args: []string{"root", "sub", "--alpha", "hello"}, want: "hello"},
		{name: "flag after nested subcommand", args: []string{"root", "sub", "leaf", "--beta", "hi"}, want: "hi"},
		{name: "flag between subcommands", args: []string{"root", "sub", "--alpha", "x", "leaf"}, want: "x"},
		{name: "missing", args: []string{"root", "sub"}, errStr: "one of these flags needs to be provided: alpha"},
		{name: "missing nested", args: []string{"root", "sub", "leaf"}, errStr: "one of these flags needs to be provided"},
		{name: "help command at root", args: []string{"root", "help"}},
		{name: "help command for subcommand", args: []string{"root", "help", "sub"}},
		{name: "help command on subcommand", args: []string{"root", "sub", "help"}},
		{name: "help flag on subcommand", args: []string{"root", "sub", "--help"}},
		{name: "both split across levels", args: []string{"root", "--alpha", "a", "sub", "--beta", "b"}, errStr: "option alpha cannot be set along with option beta"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got string
			cmd := newRoot(&got)
			cmd.Writer = io.Discard
			cmd.ErrWriter = io.Discard
			err := cmd.Run(buildTestContext(t), test.args)
			if test.errStr == "" {
				require.NoError(t, err)
				assert.Equal(t, test.want, got)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), test.errStr)
			assert.Empty(t, got)
		})
	}
}

func TestMutuallyExclusiveFlags_PropagateStringer(t *testing.T) {
	customStringer := func(f Flag) string {
		return "custom:" + f.Names()[0]
	}

	grp := MutuallyExclusiveFlags{
		Stringer: customStringer,
		Flags: [][]Flag{
			{
				&StringFlag{Name: "foo"},
				&BoolWithInverseFlag{Name: "bar"},
			},
			{
				&Int64Flag{Name: "baz"},
			},
		},
	}

	grp.propagateStringer()

	assert.Equal(t, "custom:foo", grp.Flags[0][0].String())
	assert.Contains(t, grp.Flags[0][1].String(), "custom:bar")
	assert.Equal(t, "custom:baz", grp.Flags[1][0].String())
}

func TestMutuallyExclusiveFlags_PropagateStringerNil(t *testing.T) {
	grp := MutuallyExclusiveFlags{
		Flags: [][]Flag{
			{
				&StringFlag{Name: "foo"},
			},
		},
	}

	// should not panic and should leave flags using the default FlagStringer
	grp.propagateStringer()

	assert.NotEqual(t, "", grp.Flags[0][0].String())
}

func TestMutuallyExclusiveFlags_PropagateStringerNilResetsToDefault(t *testing.T) {
	f := &StringFlag{Name: "foo"}
	defaultOut := f.String()

	customStringer := func(f Flag) string {
		return "custom:" + f.Names()[0]
	}

	grp := MutuallyExclusiveFlags{
		Stringer: customStringer,
		Flags:    [][]Flag{{f}},
	}
	grp.propagateStringer()
	assert.Equal(t, "custom:foo", f.String())

	// Resetting the stringer to nil via SetStringer must restore the
	// default FlagStringer-based behavior.
	grp.Stringer = nil
	f.SetStringer(nil)
	assert.Equal(t, defaultOut, f.String())
}

// nonStringerSettingFlag is a minimal Flag implementation that deliberately
// does not implement StringerSetter, so that propagateStringer must skip it
// via its type assertion rather than panicking or otherwise misbehaving.
type nonStringerSettingFlag struct {
	name string
}

func (f *nonStringerSettingFlag) String() string           { return "plain:" + f.name }
func (f *nonStringerSettingFlag) Get() any                 { return nil }
func (f *nonStringerSettingFlag) PreParse() error          { return nil }
func (f *nonStringerSettingFlag) PostParse() error         { return nil }
func (f *nonStringerSettingFlag) Set(string, string) error { return nil }
func (f *nonStringerSettingFlag) Names() []string          { return []string{f.name} }
func (f *nonStringerSettingFlag) IsSet() bool              { return false }

func TestMutuallyExclusiveFlags_PropagateStringerSkipsNonImplementor(t *testing.T) {
	customStringer := func(f Flag) string {
		return "custom:" + f.Names()[0]
	}

	plain := &nonStringerSettingFlag{name: "plain"}

	grp := MutuallyExclusiveFlags{
		Stringer: customStringer,
		Flags: [][]Flag{
			{plain},
		},
	}

	// should not panic
	grp.propagateStringer()

	// plain does not implement StringerSetter, so it must keep its own
	// String() implementation rather than picking up the custom stringer.
	assert.Equal(t, "plain:plain", grp.Flags[0][0].String())
}
