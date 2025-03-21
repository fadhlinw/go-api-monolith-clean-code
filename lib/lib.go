package lib

import "go.uber.org/fx"

// Module exports dependency
var Module = fx.Options(
	fx.Provide(NewRequestHandler),
	fx.Provide(NewEnv),
	fx.Provide(GetLogger),
	fx.Provide(NewDatabase),
	fx.Provide(NewSMTP),
	fx.Provide(NewSignature),
	fx.Provide(NewS3Uploader),
	fx.Provide(NewRedisClient),
)
