package providers

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"

	"github.com/desygner/desygner-toolchain/app/http"
	"github.com/desygner/desygner-toolchain/routes"
)

type RouteServiceProvider struct {
}

func (receiver *RouteServiceProvider) Register(app foundation.Application) {
}

func (receiver *RouteServiceProvider) Boot(app foundation.Application) {
	// Add HTTP middleware
	facades.Route().GlobalMiddleware(http.Kernel{}.Middleware()...)

	receiver.configureRateLimiting()

	// Add routes
	routes.Api()
}

func (receiver *RouteServiceProvider) configureRateLimiting() {

}
