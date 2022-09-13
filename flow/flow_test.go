package flow

import (
	"testing"
	"time"
)

func TestFlow_Start(t *testing.T) {

	f := New(&Flow{})
	f.Start()
	time.Sleep(10 * time.Second)
	f.Stop()

}
