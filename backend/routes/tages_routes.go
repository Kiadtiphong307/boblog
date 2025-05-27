package routes

import (
	"backend/controller"

	"github.com/gofiber/fiber/v2"
)

// Get Tags Routes
func GetTagsAll(app *fiber.App) {
	api := app.Group("/tags")
	api.Get("/", controller.GetTags)
}