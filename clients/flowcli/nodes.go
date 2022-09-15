package flowcli

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes"
	"github.com/NubeIO/rubix-rules/clients/nresty"
)

func (inst *FlowClient) NodePallet() ([]nodes.PalletNode, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&[]nodes.PalletNode{}).
		Get("/api/nodes/pallet"))
	if err != nil {
		return nil, err
	}
	var out []nodes.PalletNode
	out = *resp.Result().(*[]nodes.PalletNode)
	return out, nil
}

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
