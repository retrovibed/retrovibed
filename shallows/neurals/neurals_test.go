package neurals

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTextDefaults(t *testing.T) {
	tx := NewText("/tmp/model.onnx")
	require.Equal(t, "/tmp/model.onnx", tx.model)
	require.Equal(t, 256, tx.seqLen)
	require.Equal(t, int64(4096), tx.numTokens)
	require.Equal(t, int64(0), tx.pad)
	require.Equal(t, int64(1), tx.bos)
	require.Equal(t, int64(2), tx.eos)
	require.Equal(t, 4096, tx.outputLen)
}

func TestNewTextOptions(t *testing.T) {
	tx := NewText("/tmp/model.onnx", TextOptions().
		SeqLen(512).
		NumTokens(2048).
		PAD(10).
		BOS(11).
		EOS(12).
		OutputLen(8192)...,
	)
	require.Equal(t, "/tmp/model.onnx", tx.model)
	require.Equal(t, 512, tx.seqLen)
	require.Equal(t, int64(2048), tx.numTokens)
	require.Equal(t, int64(10), tx.pad)
	require.Equal(t, int64(11), tx.bos)
	require.Equal(t, int64(12), tx.eos)
	require.Equal(t, 8192, tx.outputLen)
}

func TestNewTextOptionsOverrideDefaults(t *testing.T) {
	tx := NewText("/tmp/model.onnx", TextOptions().SeqLen(128)...)
	require.Equal(t, 128, tx.seqLen)
	require.Equal(t, int64(4096), tx.numTokens)
	require.Equal(t, 4096, tx.outputLen)
}
