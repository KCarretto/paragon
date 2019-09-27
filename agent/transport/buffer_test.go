package transport_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kcarretto/paragon/agent/transport"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

//go:generate mockgen -imports=io -destination=buffer_mocks_test.go -package=transport_test io Writer,WriteCloser

func TestBufferSync(t *testing.T) {
	buffer := transport.Buffer{}
	require.NoError(t, buffer.Sync())
}

func TestBufferWrite(t *testing.T) {
	expected := []byte("Some test data")

	buffer := transport.Buffer{}

	n, err := buffer.Write(expected)
	require.NoError(t, err)
	require.Equal(t, len(expected), n)
}

func TestBufferWriteTo(t *testing.T) {
	expected := []byte("Some test data")

	// Prepare mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	writer := NewMockWriter(ctrl)
	writer.EXPECT().Write(expected).DoAndReturn(func(p []byte) (int, error) { return len(p), nil })

	// Initialize buffer and write expected data
	buffer := transport.NewBuffer(expected)

	// Write buffer to mock writer
	num, err := buffer.WriteTo(writer)
	require.NoError(t, err)
	require.Equal(t, int64(len(expected)), num)

	// Ensure buffer timestamp was updated
	require.NotZero(t, buffer.Timestamp(), "Buffer should update timestamp after successful write")
}

func TestBufferWriteToError(t *testing.T) {
	expected := []byte("Some test data")

	// Prepare mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	writer := NewMockWriter(ctrl)
	gomock.InOrder(
		writer.EXPECT().Write(expected).Return(4, errors.New("uh oh")),
		writer.EXPECT().Write(expected).DoAndReturn(func(p []byte) (int, error) { return len(p), nil }),
	)

	// Initialize buffer and write expected data
	buffer := transport.NewBuffer(expected)

	// Write buffer to mock writer (will error with "uh oh")
	_, err := buffer.WriteTo(writer)
	require.Error(t, err)

	// Ensure buffer timestamp was not updated
	require.Zero(t, buffer.Timestamp(), "Buffer should not update timestamp after failed write")

	// Write buffer to mock writer again, ensure no data was lost
	num, err := buffer.WriteTo(writer)
	require.NoError(t, err)
	require.Equal(t, int64(len(expected)), num)

	// Ensure buffer timestamp was updated
	require.NotZero(t, buffer.Timestamp(), "Buffer should update timestamp after successful write")
}
