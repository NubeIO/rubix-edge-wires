package flow

import (
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/db"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes"
	log "github.com/sirupsen/logrus"
	"time"
)

type Flow struct {
	DbPath           string `json:"dbPath"`
	storage          db.DB
	AutoStartDisable bool
}

var latestFlow []*node.Spec
var flowInst *flowctrl.Flow
var storage db.DB

func New(f *Flow) *Flow {
	storage = db.New("./data/flow.db")
	err := f.getLatestFlow()
	if err != nil {
		log.Error(err)
	}
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

	var nodesList []node.Node

	for _, n := range latestFlow {
		newNode, err := nodes.Builder(n, storage)
		if err != nil {
		}
		nodesList = append(nodesList, newNode)
	}
	flowInst.AddNodes(nodesList...)
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
			time.Sleep(100 * time.Millisecond)

		}
	}
}
