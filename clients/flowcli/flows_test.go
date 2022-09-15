package flowcli

import (
	"fmt"
	"testing"
)

var client = New(nil)

func TestFlowClient_GetBaseNodesList(t *testing.T) {
	start, err := client.GetBaseNodesList()
	fmt.Println(start, err)
	if err != nil {
		return
	}
}

func TestFlowClient_GetFlow(t *testing.T) {
	start, err := client.GetFlow()
	fmt.Println(start, err)
	if err != nil {
		return
	}
}

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
