package utiltime

import (
	"time"
)

// RealTicker implements a more useful ticker using the builtin ticker behind
// the scenes.
type RealTicker struct {
	Period      time.Duration
	Ticker      *time.Ticker
	C           <-chan time.Time
	c           chan time.Time
	resetSignal chan struct{}
}

func NewRealTicker(period time.Duration) *RealTicker {
	artificialChannel := make(chan time.Time)
	ticker := &RealTicker{
		Period:      period,
		Ticker:      time.NewTicker(period),
		C:           artificialChannel,
		c:           artificialChannel,
		resetSignal: make(chan struct{})}
	go ticker.monitor()
	return ticker
}

func (t *RealTicker) monitor() {
	var currentChannel <-chan time.Time

GetNewChannel:
	currentChannel = t.Ticker.C

	for {
		select {
		case value := <-currentChannel:
			// Relay ticker value to artificial ticker channel
			t.c <- value
		case <-t.resetSignal:
			// Restart the monitor with the channel of the new ticker
			goto GetNewChannel
		}
	}
}

func (t *RealTicker) Reset() {
	t.Ticker = time.NewTicker(t.Period)
	t.resetSignal <- struct{}{}
}
