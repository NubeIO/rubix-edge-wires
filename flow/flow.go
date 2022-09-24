package flow

import (
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes"
	"github.com/NubeIO/rubix-rules/storage"
	"time"
)

type Flow struct {
	DbPath           string `json:"dbPath"`
	storage          storage.Storage
	AutoStartDisable bool
}

var latestFlow []*node.Spec
var flowInst *flowctrl.Flow

func New(f *Flow) *Flow {
	f.storage = storage.New("./data/flow.db")
	f.getLatestFlow()
	flowInst = flowctrl.New()
	if !f.AutoStartDisable {
		f.Start()
	}
	return f
}

var quit chan struct{}

type Message struct {
	Message string
}

func loop() {

	for _, n := range latestFlow {
		node_, err := nodes.Builder(n)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("ADD:", node_.GetName(), node_.GetNodeName(), "ERR", err)
		flowInst.AddNode(node_)
	}

	flowInst.ReBuildFlow(true)

	runner := flowctrl.NewSerialRunner(flowInst)
	for {
		select {
		case <-quit:
			return
		default:
			err := runner.Process()
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(200 * time.Millisecond)

		}
	}
}
