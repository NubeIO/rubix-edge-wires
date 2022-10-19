package controller

import (
	"github.com/NubeDev/flow-eng/db"
	"github.com/gin-gonic/gin"
)

func (inst *Controller) AddConnection(c *gin.Context) {
	var body *db.Connection
	err := c.ShouldBindJSON(&body)
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	resp, err := inst.Flow.AddConnection(body)
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(resp, err, c)
}

func (inst *Controller) UpdateConnection(c *gin.Context) {
	var body *db.Connection
	err := c.ShouldBindJSON(&body)
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	resp, err := inst.Flow.UpdateConnection(c.Param("uuid"), body)
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(resp, err, c)
}

func (inst *Controller) GetConnections(c *gin.Context) {
	resp, err := inst.Flow.GetConnections()
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(resp, err, c)
}

func (inst *Controller) DeleteConnection(c *gin.Context) {
	err := inst.Flow.DeleteConnection(c.Param("uuid"))
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(Message{Message: "ok"}, err, c)
}

func (inst *Controller) GetConnection(c *gin.Context) {
	resp, err := inst.Flow.GetConnection(c.Param("uuid"))
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(resp, err, c)
}
