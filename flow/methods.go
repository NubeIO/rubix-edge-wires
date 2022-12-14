package flow

import (
	"encoding/json"
	"errors"
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/db"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes"
	pprint "github.com/NubeIO/rubix-edge-wires/helpers/print"
	"github.com/mitchellh/mapstructure"
)

func (inst *Flow) NodesValue(uuid string) (*node.Values, error) {
	return inst.getFlowInst().NodesValue(uuid)
}

func (inst *Flow) NodeSchema(nodeName string) (interface{}, error) {
	for _, n := range inst.getFlowInst().GetNodes() {
		if nodeName == n.GetName() {
			if n.GetSchema() != nil {
				return n.GetSchema(), nil
			}
		}
	}
	schema, err := nodes.GetSchema(nodeName)
	if err != nil {
		return nil, err
	}
	return schema, nil
}

// SetNodePayload write value to a node from an api
func (inst *Flow) SetNodePayload(uuid string, payload *node.Payload) (*flowctrl.Message, error) {
	return inst.getFlowInst().SetNodePayload(uuid, payload)
}

// NodesValues get all the node current values from the runtime
func (inst *Flow) NodesValues() []*node.Values {
	return inst.getFlowInst().NodesValues()
}

// NodesValuesInsideParent get all the node current values from the runtime for one parent
func (inst *Flow) NodesValuesInsideParent(parentID string) []*node.Values {
	return inst.getFlowInst().NodesValuesInsideParent(parentID)
}

// NodesValuesSubFlow get all the node current values from the runtime for a parent node with sub-flow inputs and outputs values
func (inst *Flow) NodesValuesSubFlow(parentID string) []*node.Values {
	return inst.getFlowInst().NodesValuesSubFlow(parentID)
}

func (inst *Flow) NodePallet() ([]*nodes.PalletNode, error) {
	return nodes.EncodePallet()
}

// DownloadFlow to the flow-eng
func (inst *Flow) DownloadFlow(encodedNodes *nodes.NodesList, restartFlow, saveFlowToDB bool) (*Message, error) {
	nodeList := &nodes.NodesList{}
	fmt.Println("From UI")
	pprint.PrintJOSN(encodedNodes)
	err := mapstructure.Decode(encodedNodes, &nodeList)
	if err != nil {
		return nil, err
	}
	decode, err := inst.decode(nodeList)
	if err != nil {
		return nil, err
	}
	err = inst.setLatestFlow(decode, saveFlowToDB)
	if err != nil {
		return nil, err
	}
	if restartFlow {
		inst.Restart()
	}
	return &Message{fmt.Sprintf("downloaded ok node count %d", len(decode))}, err
}

// decode flow data from the UI
func (inst *Flow) decode(encodedNodes *nodes.NodesList) ([]*node.Spec, error) {
	return nodes.Decode(encodedNodes)
}

func (inst *Flow) GetFlow() (*nodes.NodesList, error) {
	if inst.getFlowInst() != nil {
		return nodes.Encode(inst.getFlowInst().Get())
	}
	return nil, errors.New("failed to get flow instance")
}

func (inst *Flow) WipeFlow() []*node.Spec {
	if inst.getFlowInst() != nil {
		inst.getFlowInst().WipeFlow()
	}
	return nil
}

// GetBaseNodesList the current list of supported nodes from the base flow-eng lib
func (inst *Flow) GetBaseNodesList() []*node.Spec {
	if inst.getFlowInst() != nil {
		return nodes.All()
	}
	return nil
}

func (inst *Flow) getFlowInst() *flowctrl.Flow {
	return flowInst
}

func (inst *Flow) Restart() *Message {
	inst.Stop()
	inst.Start()
	return &Message{"restarted started ok"}
}

func (inst *Flow) Start() *Message {
	beforeStart()
	go start()
	quit = make(chan struct{})
	return &Message{"started ok"}
}

func (inst *Flow) Stop() *Message {
	beforeStop()
	quit <- struct{}{}
	inst.WipeFlow()
	return &Message{"stop ok"}
}

func (inst *Flow) setLatestFlow(flow []*node.Spec, saveFlowToDB bool) error {
	if saveFlowToDB {
		_, err := inst.saveFlowDB(flow)
		if err != nil {
			return err
		}
	}
	latestFlow = flow
	return nil
}

func (inst *Flow) getLatestFlow() error {
	backup, err := storage.GetLatestBackup()
	if err != nil {
		return err
	}
	var nodeList []*node.Spec
	b, err := json.Marshal(backup.Data)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &nodeList); err != nil {
		return err
	}
	latestFlow = nodeList
	return nil
}

func (inst *Flow) saveFlowDB(flow []*node.Spec) (*db.Backup, error) {
	back := &db.Backup{Data: flow}
	return storage.AddBackup(back)
}

func (inst *Flow) NodesHelp() []*node.Help {
	return nodes.NodeHelp()
}

func (inst *Flow) NodeHelpByName(nodeName string) (*node.Help, error) {
	return nodes.NodeHelpByName(nodeName)
}
