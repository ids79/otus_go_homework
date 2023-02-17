package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("map of environments", func(t *testing.T) {
		result, err := ReadDir("./testdata/env/")
		require.Nil(t, err)
		require.Equal(t, result["BAR"].Value, "bar")
		require.False(t, result["BAR"].NeedRemove)
		require.Equal(t, result["EMPTY"].Value, "")
		require.False(t, result["EMPTY"].NeedRemove)
		require.Equal(t, result["FOO"].Value, "   foo\nwith new line")
		require.False(t, result["FOO"].NeedRemove)
		require.Equal(t, result["HELLO"].Value, `"hello"`)
		require.False(t, result["HELLO"].NeedRemove)
		require.True(t, result["UNSET"].NeedRemove)
	})
}
