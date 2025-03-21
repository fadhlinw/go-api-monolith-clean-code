package validation

import "github.com/golodash/galidator"

var (
	g = galidator.New().CustomMessages(galidator.Messages{
		"required": "$field is required",
		"min":      "$field is too short",
		"max":      "$field is too long",
	})
)

func ParseFieldErrors(err error, fields interface{}) interface{} {
	customizer := g.Validator(fields)
	return customizer.DecryptErrors(err)
}
