package storage

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"io/ioutil"
	"os"
	"testing"
)

func TestInitializeBuntDB(t *testing.T) {

	var nodesParsed []*node.Spec
	jsonFile, err := os.Open("./test.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &nodesParsed)

	//home, _ := fileutils.HomeDir()

	db := New(fmt.Sprintf("../data/flow.db"))
	_, err = db.AddBackup(&Backup{
		Data: nodesParsed,
	})

	backup, err := db.GetLatestBackup()
	if err != nil {
		return
	}
	fmt.Println(backup)

	backups, err := db.GetBackups()
	if err != nil {
		return
	}
	for i, b := range backups {
		fmt.Println(i, b.Time)
	}

}
