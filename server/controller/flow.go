package controller

import (
	"github.com/NubeDev/flow-eng/nodes"
	"github.com/gin-gonic/gin"
)

func (inst *Controller) NodeSchema(c *gin.Context) {
	resp, err := inst.Flow.NodeSchema(c.Param("node"))
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

func (inst *Controller) DownloadFlow(c *gin.Context) {
	var body *nodes.NodesList
	err := c.ShouldBindJSON(&body)
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	resp, err := inst.Flow.DownloadFlow(body, true, true)
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(resp, nil, c)
}

func (inst *Controller) GetFlow(c *gin.Context) {
	resp := inst.Flow.GetFlow()
	reposeHandler(resp, nil, c)
}

func (inst *Controller) GetBaseNodesList(c *gin.Context) {
	resp := inst.Flow.GetBaseNodesList()
	reposeHandler(resp, nil, c)
}

func (inst *Controller) RestartFlow(c *gin.Context) {
	resp := inst.Flow.Restart()
	reposeHandler(resp, nil, c)
}

func (inst *Controller) StartFlow(c *gin.Context) {
	resp := inst.Flow.Start()
	reposeHandler(resp, nil, c)
}

func (inst *Controller) StopFlow(c *gin.Context) {
	resp := inst.Flow.Stop()
	reposeHandler(resp, nil, c)
}
