package flow

import (
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/db"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/mqttclient"
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
	mqttClient, err := mqttclient.NewClient(mqttclient.ClientOptions{
		Servers: []string{"tcp://0.0.0.0:1883"},
	})
	err = mqttClient.Connect()
	if err != nil {
		log.Error(err)
	}

	opts := &bacnet.Bacnet{
		Store:       points.New(names.Edge, nil, 0, 200, 200),
		MqttClient:  mqttClient,
		Application: names.RubixIO,
	}
	for _, n := range latestFlow {
		var node_ node.Node
		if n.Info.Category == "bacnet" {
			node_, err = nodes.Builder(n, storage, opts)
		} else {
			node_, err = nodes.Builder(n, storage)
		}

		if err != nil {
		}
		nodesList = append(nodesList, node_)
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
