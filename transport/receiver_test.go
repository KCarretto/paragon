package transport_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kcarretto/paragon/transport"
	"github.com/kcarretto/paragon/transport/mocks"
	"github.com/stretchr/testify/require"
)

func TestReceiver(t *testing.T) {
	expectedData := []byte("here is a whole bunch of payload data")
	expectedPayload := transport.Payload{}

	// Prepare mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Prepare mocks
	decoder := mocks.NewMockDecoder(ctrl)
	decoder.EXPECT().Decode(expectedData).Return(expectedPayload, nil)

	recv := transport.Receiver{
		Decoder: transport.DecoderFn(func(data []byte) (transport.Payload, error) {
			return decoder.Decode(data)
		}),
		Handler: func(payload transport.Payload, err error) {
			require.Equal(t, expectedPayload, payload)
			require.NoError(t, err)
		},
	}

	recv.WritePayload(expectedData)
}
