package router

import (
	"fmt"
	"github.com/NubeIO/rubix-edge-wires/config"
	"github.com/NubeIO/rubix-edge-wires/flow"
	"github.com/NubeIO/rubix-edge-wires/logger"
	"github.com/NubeIO/rubix-edge-wires/server/constants"
	"github.com/NubeIO/rubix-edge-wires/server/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"time"
)

func NotFound() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		message := fmt.Sprintf("%s %s [%d]: %s", ctx.Request.Method, ctx.Request.URL, 404, "rubix-edge-bios: api not found")
		ctx.JSON(http.StatusNotFound, controller.Message{Message: message})
	}
}

func Setup() *gin.Engine {
	engine := gin.New()
	// Set gin access logs
	if viper.GetBool("gin.log.store") {
		fileLocation := fmt.Sprintf("%s/rubix-edge-wires.access.log", config.Config.GetAbsDataDir())
		f, err := os.OpenFile(fileLocation, os.O_CREATE|os.O_WRONLY|os.O_APPEND, constants.Permission)
		if err != nil {
			logger.Logger.Errorf("Failed to create access log file: %v", err)
		} else {
			gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
		}
	}
	gin.SetMode(viper.GetString("gin.log.level"))
	engine.NoRoute(NotFound())
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "DELETE", "OPTIONS", "PUT", "PATCH"},
		AllowHeaders: []string{
			"X-FLOW-Key", "Authorization", "Content-Type", "Upgrade", "Origin",
			"Connection", "Accept-Encoding", "Accept-Language", "Host",
		},
		ExposeHeaders:          []string{"Content-Length"},
		AllowCredentials:       true,
		AllowAllOrigins:        true,
		AllowBrowserExtensions: true,
		MaxAge:                 12 * time.Hour,
	}))

	f := flow.New(&flow.Flow{})
	api := controller.Controller{Flow: f}

	apiRoutes := engine.Group("/api")

	flowEng := apiRoutes.Group("/flows")
	{
		flowEng.GET("", api.GetFlow)
		flowEng.POST("/download", api.DownloadFlow)
		flowEng.POST("/start", api.StartFlow)
		flowEng.POST("/restart", api.RestartFlow)
		flowEng.POST("/stop", api.StopFlow)
	}

	flowEngNodes := apiRoutes.Group("/nodes")
	{
		flowEngNodes.GET("/schema/:node", api.NodeSchema)
		flowEngNodes.GET("/values", api.NodesValues)
		flowEngNodes.GET("/values/parent/:uuid", api.NodesValuesInsideParent)
		flowEngNodes.GET("/values/sub/:uuid", api.NodesValuesSubFlow)
		flowEngNodes.POST("/payload/:uuid", api.SetNodePayload)
		flowEngNodes.GET("/values/:uuid", api.NodesValue)
		flowEngNodes.GET("/pallet", api.NodePallet)
		flowEngNodes.GET("/help", api.NodesHelp)
		flowEngNodes.GET("/help/:node", api.NodeHelpByName)
		flowEngNodes.GET("", api.GetBaseNodesList)
	}

	connections := apiRoutes.Group("/connections")
	{
		connections.GET("", api.GetConnections)
		connections.GET("/:uuid", api.GetConnection)
		connections.PATCH("/:uuid", api.UpdateConnection)
		connections.DELETE("/:uuid", api.DeleteConnection)
		connections.POST("", api.AddConnection)
	}

	return engine
}
