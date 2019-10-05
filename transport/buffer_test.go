package transport_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kcarretto/paragon/transport"
	"github.com/kcarretto/paragon/transport/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

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
	// Prepare test data
	expectedOutput := []byte("Some test output data")
	expectedResult := transport.NewResult(transport.Task{
		ID: "SomeTestTask",
	})
	_, err := expectedResult.Write([]byte("some test result data"))
	require.NoError(t, err)

	response := transport.Response{
		transport.Metadata{},
		[]*transport.Result{
			expectedResult,
		},
		expectedOutput,
	}
	encoder := transport.NewDefaultEncoder()
	expected, err := encoder.Encode(response)
	require.NoError(t, err)

	// Prepare mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dst := mocks.NewMockWriter(ctrl)
	dst.EXPECT().Write(gomock.Any()).DoAndReturn(func(p []byte) (int, error) {
		require.Equal(t, string(expected), string(p))
		return len(p), nil
	})

	// Initialize buffer and write expected data
	buffer := transport.NewBuffer(expectedOutput)

	// Write result to buffer
	buffer.WriteResult(expectedResult)

	// Write buffer to mock writer
	num, err := buffer.WriteTo(dst)
	require.NoError(t, err)
	require.Equal(t, int64(len(expected)), num)

	// Ensure buffer timestamp was updated
	require.NotZero(t, buffer.Timestamp(), "Buffer should update timestamp after successful write")
}

func TestBufferWriteToError(t *testing.T) {
	// Prepare test data
	expectedOutput := []byte("Some test output data")
	expectedResult := transport.NewResult(transport.Task{
		ID: "SomeTestTask",
	})
	_, err := expectedResult.Write([]byte("some test result data"))
	require.NoError(t, err)

	response := transport.Response{
		transport.Metadata{},
		[]*transport.Result{
			expectedResult,
		},
		expectedOutput,
	}
	encoder := transport.NewDefaultEncoder()
	expected, err := encoder.Encode(response)
	require.NoError(t, err)

	// Prepare mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dst := mocks.NewMockWriter(ctrl)
	gomock.InOrder(
		dst.EXPECT().Write(gomock.Any()).DoAndReturn(func(p []byte) (int, error) {
			require.Equal(t, string(expected), string(p))
			return 4, errors.New("uh oh")
		}),
		dst.EXPECT().Write(gomock.Any()).DoAndReturn(func(p []byte) (int, error) {
			require.Equal(t, string(expected), string(p))
			return len(p), nil
		}),
	)

	// Initialize buffer and write expected data
	buffer := transport.NewBuffer(expectedOutput)

	// Write result to buffer
	buffer.WriteResult(expectedResult)

	// Write buffer to mock writer (will error with "uh oh")
	_, err = buffer.WriteTo(dst)
	require.Error(t, err)

	// Ensure buffer timestamp was not updated
	require.Zero(t, buffer.Timestamp(), "Buffer should not update timestamp after failed write")

	// Write buffer to mock writer again, ensure no data was lost
	num, err := buffer.WriteTo(dst)
	require.NoError(t, err)
	require.Equal(t, int64(len(expected)), num)

	// Ensure buffer timestamp was updated
	require.NotZero(t, buffer.Timestamp(), "Buffer should update timestamp after successful write")
}
