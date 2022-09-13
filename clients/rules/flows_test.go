package rules

import (
	"fmt"
	"testing"
)

var client = New(nil)

func TestFlowClient_FlowStop(t *testing.T) {
	start, err := client.FlowStop()
	fmt.Println(start, err)
	if err != nil {
		return
	}
}

func TestFlowClient_FlowStart(t *testing.T) {
	start, err := client.FlowStart()
	fmt.Println(start, err)
	if err != nil {
		return
	}
}
