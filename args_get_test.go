package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringSliceArgsGetOutOfRange(t *testing.T) {
	empty := &stringSliceArgs{}
	require.NotPanics(t, func() { _ = empty.Get(empty.Len() - 1) })
	require.Equal(t, "", empty.Get(-1))
	require.Equal(t, "", empty.Get(0))
	require.Equal(t, "", empty.Get(1))

	args := &stringSliceArgs{v: []string{"a", "b"}}
	require.Equal(t, "a", args.Get(0))
	require.Equal(t, "b", args.Get(args.Len()-1))
	require.Equal(t, "", args.Get(2))
	require.Equal(t, "", args.Get(-1))
	require.Equal(t, "", args.Get(-100))
}
