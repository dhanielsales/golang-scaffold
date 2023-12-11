package category_http

import (
	"fmt"

	"github.com/dhanielsales/golang-scaffold/internal/http"
)

func SetupRoutes(httpServer *http.HttpServer, controller *CategoryController) {
	// add middlewares here

	httpServer.App.POST("/category", controller.create)
	httpServer.App.GET("/category", controller.getAll)

	fmt.Println("Category routes loaded")
}
