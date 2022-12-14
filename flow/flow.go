package flow

import (
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/db"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/helpers/store"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes"
	bacnetio "github.com/NubeDev/flow-eng/nodes/protocols/bacnet"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/nodes/protocols/driver"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	"github.com/NubeIO/rubix-edge-wires/config"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type Flow struct {
	DbPath           string `json:"dbPath"`
	AutoStartDisable bool
}

var runner *flowctrl.SerialRunner
var latestFlow []*node.Spec
var flowInst *flowctrl.Flow
var storage db.DB
var cacheStore *store.Store
var bacnetStore *bacnetio.Bacnet
var networksPool driver.Driver
var mqttClient *mqttclient.Client

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
	err = f.addDefaultConnection()
	if err != nil {
		log.Error(err)
	}
	return f
}

var quit chan struct{}

type Message struct {
	Message string
}

func makeStore() *store.Store {
	return store.Init()
}

func makeBacnetStore(application string, deviceCount int) *bacnetio.Bacnet {
	ip := "0.0.0.0"
	var err error
	mqttClient, err = mqttclient.NewClient(mqttclient.ClientOptions{
		Servers: []string{fmt.Sprintf("tcp://%s:1883", ip)},
	})
	fmt.Println("mqttClient>>>>", mqttClient)
	if err != nil {
		log.Error(err)
	}
	err = mqttClient.Connect()
	if err != nil {
		log.Error(err)
	}
	app := names.ApplicationName(application)
	bacnet := &bacnetio.Bacnet{
		Store:       points.New(names.ApplicationName(application), nil, deviceCount, 200, 200),
		MqttClient:  mqttClient,
		Application: app,
		Ip:          ip,
	}
	return bacnet
}

func beforeStart() {
	cacheStore = makeStore()
	networksPool = driver.New(&driver.Networks{}) // flow-framework networks instance

	app := names.Modbus
	var deviceCount *string
	for _, n := range latestFlow {
		if n.GetName() == "bacnet-server" {
			schema, _ := bacnetio.GetBacnetSchema(n.Settings)
			if schema != nil {
				deviceCount = &schema.DeviceCount
			}
			break
		}
	}
	if deviceCount != nil {
		i, _ := strconv.Atoi(*deviceCount)
		bacnetStore = makeBacnetStore(string(app), i)
	}
}

func beforeStop() {
	if runner != nil {
		runner.Stop()
	}

	if cacheStore != nil {
		cacheStore.Store.Flush()
	}

	if networksPool != nil {
		networksPool = nil
	}

	if mqttClient != nil {
		if mqttClient.IsConnected() {
			mqttClient.Close()
		}
	}
	if bacnetStore != nil {
		bacnetStore = nil
	}
}

func start() {
	var err error
	var nodesList []node.Node
	var parentList = nodes.FilterNodes(latestFlow, nodes.FilterIsParent, "")
	var parentChildList = nodes.FilterNodes(latestFlow, nodes.FilterIsParentChild, "")
	var childList = nodes.FilterNodes(latestFlow, nodes.FilterIsChild, "")
	var nonChildNodes = nodes.FilterNodes(latestFlow, nodes.FilterNonContainer, "")

	// add the container nodes first, then the children and so on
	for _, n := range parentList {
		var node_ node.Node
		if n.Info.Category == "bacnet" {
			node_, err = nodes.Builder(n, storage, cacheStore, bacnetStore)
		} else if n.Info.Category == "flow" {
			node_, err = nodes.Builder(n, storage, cacheStore, networksPool)
		} else {
			node_, err = nodes.Builder(n, storage, cacheStore)
		}
		nodesList = append(nodesList, node_)
	}

	for _, n := range parentChildList {
		var node_ node.Node
		if n.Info.Category == "flow" {
			node_, err = nodes.Builder(n, storage, cacheStore, networksPool)
		} else {
			node_, err = nodes.Builder(n, storage, cacheStore)
		}
		nodesList = append(nodesList, node_)
	}

	for _, n := range childList {
		var node_ node.Node
		if n.Info.Category == "bacnet" {
			node_, err = nodes.Builder(n, storage, cacheStore, bacnetStore)
		} else if n.Info.Category == "flow" {
			node_, err = nodes.Builder(n, storage, cacheStore, networksPool)
		} else {
			node_, err = nodes.Builder(n, storage, cacheStore)
		}
		nodesList = append(nodesList, node_)
	}
	for _, n := range nonChildNodes {
		var node_ node.Node
		node_, err = nodes.Builder(n, storage, cacheStore)
		nodesList = append(nodesList, node_)
	}

	if err != nil {
		log.Error(err)
	}
	flowInst.AddNodes(nodesList...)
	flowInst.MakeNodeConnections(true)
	flowInst.MakeGraph()
	for _, n := range flowInst.Get().GetNodes() { // add all nodes to each node so data can be passed between nodes easy
		n.AddNodes(flowInst.Get().GetNodes())
	}
	log.Infof("graphs count: %d nodes count: %d", len(flowInst.Graphs), len(flowInst.GetNodes()))
	runner = flowctrl.NewSerialRunner(flowInst)
	runner.Start()
	for {
		select {
		case <-quit:
			return
		default:
			runner.Process()
			time.Sleep(100 * time.Millisecond)
		}
	}
}
