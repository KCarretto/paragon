package transport_test

import (
	"testing"

	"github.com/kcarretto/paragon/agent/transport"
	"github.com/stretchr/testify/require"
)

func TestBufferSync(t *testing.T) {
	buffer := transport.NewBuffer(1)
	require.NoError(t, buffer.Sync())
}

func TestBufferWrite(t *testing.T) {
	panic("Not implemented")
}

func TestBufferWriteTo(t *testing.T) {
	panic("Not implemented")
}

func TestBufferWriteTimestamp(t *testing.T) {
	panic("Not implemented")
}

func TestBufferNoLossOnError(t *testing.T) {
	panic("Not implemented")
}
