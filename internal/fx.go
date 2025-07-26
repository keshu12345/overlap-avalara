package internal

import (
	"github.com/keshu12345/overlap-avalara/internal/api"
	"github.com/keshu12345/overlap-avalara/internal/overlap"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Invoke(api.RegisterEndpoint),
	fx.Provide(overlap.New),
)
