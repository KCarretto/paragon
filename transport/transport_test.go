package transport_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/kcarretto/paragon/transport"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	name := "MyTransport"
	interval := time.Hour * 5
	jitter := time.Nanosecond * 10
	priority := 1337

	transport := transport.New(
		name,
		ioutil.Discard,
		transport.SetInterval(interval),
		transport.SetJitter(jitter),
		transport.SetPriority(priority),
	)

	require.Equal(t, name, transport.Name)
	require.Equal(t, interval, transport.Interval)
	require.Equal(t, jitter, transport.Jitter)
	require.Equal(t, priority, transport.Priority)
}
