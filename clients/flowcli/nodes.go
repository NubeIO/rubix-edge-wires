package flowcli

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeIO/rubix-rules/clients/nresty"
)

func (inst *FlowClient) GetBaseNodesList() ([]node.Spec, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&[]node.Spec{}).
		Get("/api/nodes"))
	if err != nil {
		return nil, err
	}
	var out []node.Spec
	out = *resp.Result().(*[]node.Spec)
	return out, nil
}

type Nodes struct {
	Inputs []struct {
		Name       string `json:"name"`
		Type       string `json:"type"`
		Connection struct {
			OverrideValue int    `json:"overrideValue,omitempty"`
			NodeID        string `json:"nodeID,omitempty"`
			NodePortName  string `json:"nodePortName,omitempty"`
		} `json:"connection"`
	} `json:"inputs"`
	Outputs []struct {
		Name       string `json:"name"`
		Type       string `json:"type"`
		Connection []struct {
			NodeID       string `json:"nodeID"`
			NodePortName string `json:"nodePortName"`
		} `json:"connection"`
	} `json:"outputs"`
	Info struct {
		NodeID   string `json:"nodeID"`
		Name     string `json:"name"`
		NodeName string `json:"nodeName"`
		Category string `json:"category"`
	} `json:"info"`
	Settings []struct {
		Type       string `json:"type"`
		Title      string `json:"title"`
		Properties struct {
			Type      string      `json:"type"`
			Title     string      `json:"title"`
			MinLength int         `json:"minLength"`
			MaxLength int         `json:"maxLength"`
			ReadOnly  interface{} `json:"readOnly"`
			Value     interface{} `json:"value"`
		} `json:"properties"`
	} `json:"settings,omitempty"`
}
