package cmd

import (
	"fmt"
	"github.com/NubeIO/rubix-edge-wires/config"
	"github.com/NubeIO/rubix-edge-wires/logger"
	"github.com/NubeIO/rubix-edge-wires/server/constants"
	"github.com/NubeIO/rubix-edge-wires/server/router"
	"github.com/spf13/cobra"
	"os"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starting rubix-edge-wires",
	Long:  "it starts a server for edge-wires flow-engine",
	Run:   runServer,
}

func runServer(cmd *cobra.Command, args []string) {
	if err := config.Setup(RootCmd); err != nil {
		fmt.Errorf("error: %s", err) // here log is not setup yet...
	}
	logger.Init()
	if err := os.MkdirAll(config.Config.GetAbsDataDir(), os.FileMode(constants.Permission)); err != nil {
		panic(err)
	}

	logger.Logger.Infoln("starting rubix-edge-wires...")

	r := router.Setup(flgRoot.runFlow)

	host := "0.0.0.0"
	port := config.Config.GetPort()
	logger.Logger.Infof("server is starting at %s:%s", host, port)
	logger.Logger.Fatalf("%v", r.Run(fmt.Sprintf("%s:%s", host, port)))
}

func init() {
	RootCmd.AddCommand(serverCmd)
}
