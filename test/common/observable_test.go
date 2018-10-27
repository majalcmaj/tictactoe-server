package test

import (
	"fmt"
	"testing"

	"github.com/majalcmaj/tictactoe-server/common"
	"github.com/stretchr/testify/require"
)

type testSubscriber struct {
	wasEventFired bool
	firesCounter  int
}

func (ts *testSubscriber) EventFired(value int) {
	ts.wasEventFired = true
	ts.firesCounter++
}

func TestOneSubscriber(t *testing.T) {
	emitter := common.NewIntEventEmitter()
	sub := testSubscriber{}
	emitter.Subscribe(&sub)
	emitter.FireEvent(10)
	require.True(t, sub.wasEventFired, "Subscriber got no event")
}

func TestMultipleSubscribers(t *testing.T) {
	require := require.New(t)
	emitter := common.NewIntEventEmitter()
	const subscribersCount = 10
	var subs [10]*testSubscriber
	for i := 0; i < subscribersCount; i++ {
		sub := &testSubscriber{}
		emitter.Subscribe(sub)
		subs[i] = sub
	}
	emitter.FireEvent(10)
	for i := 0; i < subscribersCount; i++ {
		require.True(subs[i].wasEventFired, fmt.Sprint("No  event for subscriber with index ", i))
	}
}

func TestMultipleEvents(t *testing.T) {
	emitter := common.NewIntEventEmitter()
	sub := testSubscriber{}
	emitter.Subscribe(&sub)
	emitter.FireEvent(10)
	emitter.FireEvent(11)
	emitter.FireEvent(12)
	require.True(t, sub.wasEventFired, "Subscriber got no event")
	require.Equal(t, sub.firesCounter, 3, "Subscriber was notified wrong count of times")
}
