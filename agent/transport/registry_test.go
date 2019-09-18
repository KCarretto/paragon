package transport_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kcarretto/paragon/agent/transport"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=transport_mocks_test.go -package=transport_test -source=transport.go

func TestRegistryGet(t *testing.T) {
	// Prepare mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	writer := NewMockWriteCloser(ctrl)
	tasker := NewMockTasker(ctrl)
	factory := NewMockFactory(ctrl)
	factory.EXPECT().New(gomock.Not(nil), gomock.Not(nil)).Return(writer, nil)

	reg := transport.Registry{}
	reg.Add("TestAdd", factory)

	meta, err := reg.Get("TestAdd", zap.NewNop(), tasker)
	require.NoError(t, err)
	require.NotNil(t, meta)
}
func TestRegistryRemove(t *testing.T) {
	panic("Not implemented")
}

func TestRegistryUpdate(t *testing.T) {
	panic("Not implemented")
}
func TestRegistryList(t *testing.T) {
	panic("Not implemented")
}

func TestRegistryGetNotRegistered(t *testing.T) {
	panic("Not implemented")
}

func TestRegistryGetInvalid(t *testing.T) {
	panic("Not implemented")
}

func TestRegistryGetAlreadyActive(t *testing.T) {
	panic("Not implemented")
}
