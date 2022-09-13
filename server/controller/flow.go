package controller

import (
	"github.com/gin-gonic/gin"
)

func (inst *Controller) StartFlow(c *gin.Context) {
	resp := inst.Flow.Start()
	reposeHandler(resp, nil, c)
}

func (inst *Controller) StopFlow(c *gin.Context) {
	resp := inst.Flow.Stop()
	reposeHandler(resp, nil, c)
}
