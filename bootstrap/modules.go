package bootstrap

import (
	"gitlab.com/tsmdev/software-development/backend/go-project/api/controllers"
	"gitlab.com/tsmdev/software-development/backend/go-project/api/middlewares"
	"gitlab.com/tsmdev/software-development/backend/go-project/api/routes"
	"gitlab.com/tsmdev/software-development/backend/go-project/lib"
	"gitlab.com/tsmdev/software-development/backend/go-project/repository"
	"gitlab.com/tsmdev/software-development/backend/go-project/services"
	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	controllers.Module,
	routes.Module,
	lib.Module,
	services.Module,
	middlewares.Module,
	repository.Module,
)
