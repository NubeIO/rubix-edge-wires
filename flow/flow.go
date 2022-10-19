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
	"github.com/NubeDev/flow-eng/nodes/protocols/driver"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	"github.com/NubeIO/rubix-edge-wires/config"
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
var bacnetStore *bacnet.Bacnet
var storage db.DB

func New(f *Flow) *Flow {
	storage = db.New(config.Config.GetAbsDatabaseFile())
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

func makeBacnetStore() *bacnet.Bacnet {
	mqttClient, err := mqttclient.NewClient(mqttclient.ClientOptions{
		Servers: []string{"tcp://0.0.0.0:1883"},
	})
	err = mqttClient.Connect()
	if err != nil {
		log.Error(err)
	}
	opts := &bacnet.Bacnet{
		Store:       points.New(names.RubixIOAndModbus, nil, 2, 200, 200),
		MqttClient:  mqttClient,
		Application: names.RubixIOAndModbus,
	}
	return opts
}

func loop() {
	var err error
	var nodesList []node.Node
	var parentList = nodes.FilterNodes(latestFlow, nodes.FilterIsParent, "")
	var parentChildList = nodes.FilterNodes(latestFlow, nodes.FilterIsParentChild, "")
	var childList = nodes.FilterNodes(latestFlow, nodes.FilterIsChild, "")
	var nonChildNodes = nodes.FilterNodes(latestFlow, nodes.FilterNonContainer, "")

	mqttClient, err := mqttclient.NewClient(mqttclient.ClientOptions{
		Servers: []string{"tcp://0.0.0.0:1883"},
	})
	err = mqttClient.Connect()
	if err != nil {
		log.Error(err)
	}
	opts := &bacnet.Bacnet{
		Store:       points.New(names.RubixIOAndModbus, nil, 2, 200, 200),
		MqttClient:  mqttClient,
		Application: names.RubixIOAndModbus,
	}

	var networksPool driver.Driver // flow networks inst
	if networksPool == nil {
		networksPool = driver.New(&driver.Networks{})
	}

	// add the container nodes first, then the children and so on
	for _, n := range parentList {
		var node_ node.Node
		if n.Info.Category == "bacnet" {
			node_, err = nodes.Builder(n, storage, opts)
		} else if n.Info.Category == "flow" {
			node_, err = nodes.Builder(n, storage, networksPool)
		} else {
			node_, err = nodes.Builder(n, storage)
		}
		nodesList = append(nodesList, node_)
	}

	for _, n := range parentChildList {
		var node_ node.Node
		if n.Info.Category == "flow" {
			node_, err = nodes.Builder(n, storage, networksPool)
		} else {
			node_, err = nodes.Builder(n, storage)
		}
		nodesList = append(nodesList, node_)
	}

	for _, n := range childList {
		var node_ node.Node
		if n.Info.Category == "flow" {
			node_, err = nodes.Builder(n, storage, networksPool)
		} else {
			node_, err = nodes.Builder(n, storage)
		}
		nodesList = append(nodesList, node_)
	}
	for _, n := range nonChildNodes {
		var node_ node.Node
		node_, err = nodes.Builder(n, storage)
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
			time.Sleep(500 * time.Millisecond)
		}
	}
}
