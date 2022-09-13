package cmd

import (
	"fmt"
	"github.com/NubeIO/rubix-rules/config"
	"github.com/NubeIO/rubix-rules/logger"
	"github.com/NubeIO/rubix-rules/server/constants"
	"github.com/NubeIO/rubix-rules/server/router"
	"github.com/spf13/cobra"
	"os"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starting rubix-rules",
	Long:  "it starts a server for rubix-rules engine",
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

	logger.Logger.Infoln("starting rubix-rules...")

	r := router.Setup()

	host := "0.0.0.0"
	port := config.Config.GetPort()
	logger.Logger.Infof("server is starting at %s:%s", host, port)
	logger.Logger.Fatalf("%v", r.Run(fmt.Sprintf("%s:%s", host, port)))
}

func init() {
	RootCmd.AddCommand(serverCmd)
}
