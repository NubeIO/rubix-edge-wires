package flow

import (
	"fmt"
	"github.com/NubeDev/flow-eng/db"
	pprint "github.com/NubeIO/rubix-edge-wires/helpers/print"
	"testing"
)

func TestFlow_Start(t *testing.T) {

	d := db.New("../data/data.db")
	connection, err := d.AddConnection(&db.Connection{
		Application: "flow",
		Name:        "flow-3",
		Host:        "0.0.0.0",
		Port:        1660,
	})
	fmt.Println(err)
	pprint.Print(connection)
	if err != nil {
		return
	}
}
