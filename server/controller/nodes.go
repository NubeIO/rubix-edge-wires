package controller

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/gin-gonic/gin"
)

func (inst *Controller) GetBaseNodesList(c *gin.Context) {
	resp := inst.Flow.GetBaseNodesList()
	reposeHandler(resp, nil, c)
}

func (inst *Controller) NodeSchema(c *gin.Context) {
	resp, err := inst.Flow.NodeSchema(c.Param("node"))
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(resp, err, c)
}

func (inst *Controller) NodesValue(c *gin.Context) {
	resp, err := inst.Flow.NodesValue(c.Param("uuid"))
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(resp, err, c)
}

func (inst *Controller) SetNodePayload(c *gin.Context) {
	var body *node.Payload
	err := c.ShouldBindJSON(&body)
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	resp, err := inst.Flow.SetNodePayload(c.Param("uuid"), body)
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(resp, err, c)
}

func (inst *Controller) NodesValues(c *gin.Context) {
	resp := inst.Flow.NodesValues()
	reposeHandler(resp, nil, c)
}

func (inst *Controller) NodePallet(c *gin.Context) {
	resp, err := inst.Flow.NodePallet()
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(resp, err, c)
}

func (inst *Controller) NodesHelp(c *gin.Context) {
	resp := inst.Flow.NodesHelp()
	reposeHandler(resp, nil, c)
}

func (inst *Controller) NodeHelpByName(c *gin.Context) {
	resp, err := inst.Flow.NodeHelpByName(c.Param("node"))
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(resp, nil, c)
}
