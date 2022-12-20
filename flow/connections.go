package flow

import (
	"github.com/NubeDev/flow-eng/db"
	"github.com/NubeDev/flow-eng/helpers/names"
)

func (inst *Flow) AddConnection(body *db.Connection) (*db.Connection, error) {
	return storage.AddConnection(body)
}

func (inst *Flow) DeleteConnection(uuid string) error {
	return storage.DeleteConnection(uuid)
}

func (inst *Flow) UpdateConnection(uuid string, body *db.Connection) (*db.Connection, error) {
	return storage.UpdateConnection(uuid, body)
}

func (inst *Flow) GetConnections() ([]db.Connection, error) {
	return storage.GetConnections()
}

func (inst *Flow) GetConnection(uuid string) (*db.Connection, error) {
	return storage.GetConnection(uuid)
}

func (inst *Flow) addDefaultConnection() error {
	c, err := inst.GetConnections()
	if err != nil {
		return err
	}
	var flowNetworkConnection = names.FlowFramework
	var found bool
	for _, connection := range c {
		if connection.Application == flowNetworkConnection {
			found = true
		}
	}
	if !found {
		_, err := inst.AddConnection(&db.Connection{
			Application: names.FlowFramework,
			Name:        "flow framework integration over MQTT (dont not edit/delete)",
			Host:        "127.0.0.1",
			Port:        1883,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
