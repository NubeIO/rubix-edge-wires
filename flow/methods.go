package flow

import (
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes"
	"github.com/NubeIO/rubix-rules/storage"
	"github.com/mitchellh/mapstructure"
)

func (inst *Flow) NodeSchema(nodeName string) (interface{}, error) {
	return nodes.GetSchema(nodeName)
}

func (inst *Flow) NodePallet() ([]*nodes.PalletNode, error) {
	return nodes.EncodePallet()
}

// DownloadFlow to the flow-eng
func (inst *Flow) DownloadFlow(encodedNodes *nodes.NodesList, restartFlow, saveFlowToDB bool) (*Message, error) {
	nodeList := &nodes.NodesList{}
	err := mapstructure.Decode(encodedNodes, &nodeList)
	if err != nil {
		return nil, err
	}
	decode, err := inst.decode(encodedNodes)
	if err != nil || decode == nil {
		return nil, err
	}
	err = inst.setLatestFlow(decode, saveFlowToDB)
	if err != nil {
		return nil, err
	}
	if restartFlow {
		inst.Restart()
	}
	return &Message{"downloaded new flow ok"}, err
}

//decode flow data from the UI
func (inst *Flow) decode(encodedNodes *nodes.NodesList) ([]*node.Spec, error) {
	return nodes.Decode(encodedNodes)
}

//encode the flow to send to UI
func (inst *Flow) encode() (*nodes.NodesList, error) {
	return nodes.Encode(inst.getFlowInst())
}

func (inst *Flow) GetFlow() []*node.Spec {
	if inst.getFlowInst() != nil {
		return inst.getFlowInst().GetNodesSpec()
	}
	return nil
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
	quit = make(chan struct{})
	go loop()
	return &Message{"started ok"}
}

func (inst *Flow) Stop() *Message {
	inst.WipeFlow()
	quit <- struct{}{}
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

func (inst *Flow) getLatestFlow() {
	backup, err := inst.getDB().GetLatestBackup()
	if err != nil {
		return
	}
	latestFlow = backup.Data
}

func (inst *Flow) saveFlowDB(flow []*node.Spec) (*storage.Backup, error) {
	back := &storage.Backup{Data: flow}
	return inst.getDB().AddBackup(back)
}

func (inst *Flow) getDB() storage.Storage {
	return inst.storage
}
