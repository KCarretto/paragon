package transport_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kcarretto/paragon/transport"
	"github.com/kcarretto/paragon/transport/mocks"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	// Prepare mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Prepare mocks
	writer1 := mocks.NewMockWriteCloser(ctrl)
	writer2 := mocks.NewMockWriteCloser(ctrl)
	gomock.InOrder(
		writer1.EXPECT().Close(),
	)

	// Add transports to the registry
	reg := transport.Registry{}
	reg.Add(transport.New("RegTest1", writer1, transport.SetInterval(time.Second*100)))
	reg.Add(transport.New("RegTest2", writer2))

	// Get a transport
	tp1, err := reg.Get("RegTest1")
	require.NoError(t, err)
	require.Equal(t, time.Second*100, tp1.Interval)

	// Update the transport
	err = reg.Update("RegTest2", transport.SetInterval(time.Second*23), transport.SetPriority(-15))
	require.NoError(t, err)
	tp2, err := reg.Get("RegTest2")
	require.NoError(t, err)
	require.Equal(t, time.Second*23, tp2.Interval)
	require.Equal(t, -15, tp2.Priority)

	// List transports
	lst := reg.List()
	require.Equal(t, 2, len(lst))
	require.Equal(t, "RegTest2", lst[0].Name, "Transport list in unexpected order")
	require.Equal(t, "RegTest1", lst[1].Name, "Transport list in unexpected order")

	// Close the first transport
	err = reg.Close(tp1.Name)
	require.NoError(t, err)
}
