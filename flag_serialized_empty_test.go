package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSliceSetFromSerializedEmpty(t *testing.T) {
	var dest []string
	v := StringSlice{}.Create([]string{"x"}, &dest, StringConfig{})
	require.NoError(t, v.Set(NewStringSlice().Serialize()))
	assert.Nil(t, v.Get())
	assert.Nil(t, dest)
	require.NoError(t, v.Set("new"))
	assert.Equal(t, []string{"new"}, v.Get())
	assert.Equal(t, []string{"new"}, dest)

	// non-empty serialized values still overwrite
	require.NoError(t, v.Set(NewStringSlice("a", "b").Serialize()))
	assert.Equal(t, []string{"a", "b"}, dest)
}

func TestMapSetFromSerializedEmpty(t *testing.T) {
	var dest map[string]string
	v := StringMap{}.Create(map[string]string{"x": "y"}, &dest, StringConfig{})
	require.NoError(t, v.Set(NewStringMap(nil).Serialize()))
	assert.Nil(t, v.Get())
	assert.Nil(t, dest)
	require.NoError(t, v.Set("k=v"))
	assert.Equal(t, map[string]string{"k": "v"}, v.Get())
	assert.Equal(t, map[string]string{"k": "v"}, dest)
}

func TestSliceFlagSerializedEmptyFromArgs(t *testing.T) {
	var sdest []string
	var mdest map[string]string
	cmd := &Command{
		Flags: []Flag{
			&StringSliceFlag{Name: "s", Destination: &sdest},
			&StringMapFlag{Name: "m", Destination: &mdest},
		},
		Action: func(_ context.Context, _ *Command) error { return nil },
	}
	err := cmd.Run(buildTestContext(t), []string{"app",
		"--s", NewStringSlice().Serialize(), "--s", "new",
		"--m", NewStringMap(nil).Serialize(), "--m", "k=v"})
	require.NoError(t, err)
	assert.Equal(t, []string{"new"}, sdest)
	assert.Equal(t, map[string]string{"k": "v"}, mdest)
	assert.Equal(t, []string{"new"}, cmd.StringSlice("s"))
	assert.Equal(t, map[string]string{"k": "v"}, cmd.StringMap("m"))
}
