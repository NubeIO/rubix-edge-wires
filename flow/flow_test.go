package flow

import (
	"github.com/NubeDev/flow-eng/db"
	pprint "github.com/NubeIO/rubix-edge-wires/helpers/print"
	"testing"
)

func TestFlow_Start(t *testing.T) {

	d := db.New("../data/flow.db")
	connection, err := d.AddConnection(&db.Connection{
		Application: "flow-framework",
		Name:        "flow-framework",
		Host:        "0.0.0.0",
		Port:        1660,
	})
	pprint.Print(connection)
	if err != nil {
		return
	}
}
