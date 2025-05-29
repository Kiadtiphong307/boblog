
package routes

import (
	"backend/controller"
	"backend/middleware"

	"github.com/gofiber/fiber/v2"
)

// Comment Routes
func CommentRoutes(app *fiber.App) {
    articles := app.Group("/articles")
    articles.Post("/:slug/comments", middleware.Protected(), controller.CreateComment)
    articles.Get("/:slug/comments", controller.GetComments)
	articles.Put("/:slug/comments/:commentId", middleware.Protected(), controller.UpdateComment)
	articles.Delete("/:slug/comments/:commentId", middleware.Protected(), controller.DeleteComment)
}
