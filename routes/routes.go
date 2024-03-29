package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/saheemshafi/gin-basic-api/middlewares"
)

func Register(app *gin.Engine) {

	v1 := app.Group("/api/v1")

	// User routes
	users := v1.Group("/users")
	users.POST("/create-account", CreateAccount)
	users.POST("/login", Login)
	users.PUT("/", middlewares.Authorize, UpdateUser)

	// Book routes
	books := v1.Group("/books")
	books.GET("/", GetBooks)
	books.GET("/:bookId", GetBook)
	books.POST("/", middlewares.Authorize, CreateBook)
	books.PUT("/:bookId", middlewares.Authorize, UpdateBook)
	books.DELETE("/:bookId", middlewares.Authorize, DeleteBook)
	books.POST("/:bookId/pages", middlewares.Authorize, AddPage)
	books.PUT("/:bookId/pages/:pageId", middlewares.Authorize, UpdatePage)
	books.DELETE("/:bookId/pages/:pageId", middlewares.Authorize, DeletePage)

	books.PUT("/:bookId/cover", middlewares.Authorize, ChangeBookCover)
	books.PUT("/:bookId/pages/:pageId/cover", middlewares.Authorize, ChangePageCover)
}
