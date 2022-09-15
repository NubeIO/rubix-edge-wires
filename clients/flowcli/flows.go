package flowcli

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/services/clients/ffclient/nresty"
	"github.com/NubeIO/rubix-rules/flow"
)

func (inst *FlowClient) GetFlow() ([]node.Spec, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&[]node.Spec{}).
		Get("/api/flows"))
	if err != nil {
		return nil, err
	}
	var out []node.Spec
	out = *resp.Result().(*[]node.Spec)
	return out, nil
}

func (inst *FlowClient) FlowStart() (*flow.Message, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&flow.Message{}).
		Post("/api/flows/start"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*flow.Message), nil
}

func (inst *FlowClient) FlowStop() (*flow.Message, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&flow.Message{}).
		Post("/api/flows/stop"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*flow.Message), nil
}
