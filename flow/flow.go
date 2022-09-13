package flow

import (
	"encoding/json"
	"flag"
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes"
	"github.com/NubeIO/rubix-rules/storage"
	"io/ioutil"
	"os"
	"time"
)

type Flow struct {
	DbPath  string `json:"dbPath"`
	storage storage.Storage
}

func New(f *Flow) *Flow {
	f.storage = storage.New(f.DbPath)
	return f
}

var quit chan struct{}

func (inst *Flow) Start() {
	quit = make(chan struct{})
	go loop()
}

func (inst *Flow) Stop() {
	quit <- struct{}{}
}

func (inst *Flow) getFirstFlow() {

}

func (inst *Flow) getDB() storage.Storage {
	return inst.storage
}

func loop() {
	filePath := flag.String("f", fmt.Sprintf("./test.json"), "flow file")
	flag.Parse()
	fmt.Println("file:", *filePath)

	var nodesParsed []*node.BaseNode
	jsonFile, err := os.Open(*filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &nodesParsed)

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
