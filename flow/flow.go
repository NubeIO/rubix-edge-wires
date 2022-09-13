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
	DbPath  string `json:"dbPath"`
	storage storage.Storage
}

var flowFile []*node.BaseNode

func New(f *Flow) *Flow {
	f.storage = storage.New("./data/flow.db")
	f.getLatestFlow()
	return f
}

var quit chan struct{}

type Message struct {
	Message string
}

func (inst *Flow) Start() *Message {
	quit = make(chan struct{})
	go loop()
	return &Message{"started ok"}
}

func (inst *Flow) Stop() *Message {
	quit <- struct{}{}
	return &Message{"stop ok"}
}

func (inst *Flow) getLatestFlow() {
	backup, err := inst.getDB().GetLatestBackup()
	if err != nil {
		return
	}
	flowFile = backup.Data
}

func (inst *Flow) getDB() storage.Storage {
	return inst.storage
}

func loop() {

	var nodesParsed = flowFile
	graph := flowctrl.New()
	for _, n := range nodesParsed {
		node_, err := nodes.Builder(n)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("ADD:", node_.GetName(), node_.GetNodeName(), "ERR", err)
		graph.AddNode(node_)
	}

	graph.ReBuildFlow(true)

	runner := flowctrl.NewSerialRunner(graph)
	for {
		select {
		case <-quit:
			return
		default:
			err := runner.Process()
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			time.Sleep(1000 * time.Millisecond)

		}
	}
}
