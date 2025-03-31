package route

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewAppRoute),
	fx.Provide(NewRoutes),
)

type Routes []Route

type Route interface {
	Setup()
}

func NewRoutes(
	appRoutes AppRoute,
) Routes {
	return Routes{
		&appRoutes,
	}
}

func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
