package rules

import (
	"github.com/NubeDev/flow-eng/services/clients/ffclient/nresty"
	"github.com/NubeIO/rubix-rules/flow"
)

func (inst *FlowClient) FlowStart() (*flow.Message, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&flow.Message{}).
		Post("/api/flow/start"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*flow.Message), nil
}

func (inst *FlowClient) FlowStop() (*flow.Message, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&flow.Message{}).
		Post("/api/flow/stop"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*flow.Message), nil
}
