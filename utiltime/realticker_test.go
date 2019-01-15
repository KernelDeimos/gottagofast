package utiltime

import (
	"fmt"
	"testing"
	"time"
)

func ManualTestRealTicker(t *testing.T) {
	rt := NewRealTicker(time.Second)

	go func() {
		<-time.After(2 * time.Second)
		for i := 0; i < 6; i++ {
			<-time.After(500 * time.Millisecond)
			rt.Reset()
		}
	}()

	for {
		select {
		case <-rt.C:
			fmt.Println(time.Now().Unix())
		}
	}
}
