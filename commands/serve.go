package commands

import (
	"github.com/spf13/cobra"
	"gitlab.com/tsmdev/software-development/backend/go-project/api/middlewares"
	"gitlab.com/tsmdev/software-development/backend/go-project/api/routes"
	"gitlab.com/tsmdev/software-development/backend/go-project/lib"
)

// ServeCommand test command
type ServeCommand struct{}

func (s *ServeCommand) Short() string {
	return "serve application"
}

func (s *ServeCommand) Setup(cmd *cobra.Command) {}

func (s *ServeCommand) Run() lib.CommandRunner {
	return func(
		middleware middlewares.Middlewares,
		env lib.Env,
		router lib.RequestHandler,
		route routes.Routes,
		logger lib.Logger,
		database lib.Database,
		// deviceService domains.DeviceService,
		// audioPushHistory domains.AudioPushHistoryService,
		// mqtt lib.MQTTAPI,
	) {
		// Set up middleware
		middleware.Setup()

		// Set up routes
		route.Setup()

		// Initialize MQTT client
		// _, err := lib.NewMQTTClient(env, logger, deviceService, audioPushHistory)
		// if err != nil {
		// 	logger.Fatal("Failed to initialize MQTT client: ", err)
		// }

		// Start server
		logger.Info("Running server")
		if env.ServerPort == "" {
			_ = router.Gin.Run()
		} else {
			_ = router.Gin.Run(":" + env.ServerPort)
		}
	}
}

func NewServeCommand() *ServeCommand {
	return &ServeCommand{}
}
