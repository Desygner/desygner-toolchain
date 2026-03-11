package routes

import (
	"github.com/goravel/framework/facades"

	"github.com/desygner/desygner-toolchain/app/http/controllers"
)

func Api() {
	userController := controllers.NewUserController()
	cloudfront := controllers.NewCloudfrontInvalidate()
	fileController := controllers.NewFile()
	healthController := controllers.NewHealthcheck()
	routes := facades.Route()
	routes.Get("/health", healthController.Index)
	routes.Get("/users/{id}", userController.Show)
	routes.Post("/cloudfront/invalidate", cloudfront.Invalidate)
	routes.Post("/cloudfront/get-bucket", cloudfront.GetBucket)
	// file replacement
	routes.Post("/file/replace", fileController.Replace)
}
