package transport_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kcarretto/paragon/agent/transport"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=transport_mocks_test.go -package=transport_test -source=transport.go

func TestRegistry(t *testing.T) {
	// Prepare mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Prepare mocks
	writer := NewMockWriteCloser(ctrl)
	tasker := NewMockTasker(ctrl)
	factory := NewMockFactory(ctrl)
	gomock.InOrder(
		factory.EXPECT().New(gomock.Not(nil), gomock.Not(nil)).Return(writer, nil),
		writer.EXPECT().Close(),
		factory.EXPECT().New(gomock.Not(nil), gomock.Not(nil)).Return(writer, nil),
		writer.EXPECT().Close(),
	)

	// New registry with test transport
	reg := transport.Registry{}
	reg.Add("RegTest1", factory)
	reg.Add("RegTest2", factory, transport.SetPriority(25))

	// Update the transport
	err := reg.Update("RegTest1", transport.SetInterval(time.Second*23), transport.SetPriority(-15))
	require.NoError(t, err)

	// List transports
	for i, meta := range reg.List() {
		if i == 0 {
			require.Equal(t, "RegTest1", meta.Name, "Transport list in unexpected order")
			require.Equal(t, -15, meta.Priority, "Expected priority -15 for RegTest1")
		} else {
			require.Equal(t, "RegTest2", meta.Name, "Transport list in unexpected order")
			require.Equal(t, 25, meta.Priority, "Expected priority 25 for RegTest2")
		}

		// Get the transport
		tp, err := reg.Get(meta.Name, zap.NewNop(), tasker)
		require.NoError(t, err)
		require.Equal(t, writer, tp)

		// Test cached Get (Mock will error if it calls the factory)
		tp2, err := reg.Get(meta.Name, zap.NewNop(), tasker)
		require.NoError(t, err)
		require.Equal(t, writer, tp2)

		// Close the transport
		err = reg.CloseTransport(meta.Name)
		require.NoError(t, err)
	}

}
