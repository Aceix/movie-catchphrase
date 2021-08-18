package routes

import (
	"github.com/Aceix/movie-catchprase/controllers"
	"github.com/gofiber/fiber/v2"
)

func CatchphrasesRoute(route fiber.Router) {
	route.Get("/", controllers.GetCatchphrasesBy)
	route.Get("/:id", controllers.GetCatchphrase)
	route.Post("/", controllers.AddCatchphrase)
	route.Put("/:id", controllers.UpdateCatchphrase)
	route.Delete("/:id", controllers.DeleteCatchphrase)
}
