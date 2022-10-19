package flow

import "github.com/NubeDev/flow-eng/db"

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
