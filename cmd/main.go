package main

import (
	"webapp-go/webapp"
	"webapp-go/webapp/controllers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `json:"id"`
	UserName string    `json:"username" binding:"required"`
	Name     string    `json:"name"`
}

func main() {
    db := webapp.DBConnection()

    postsController := controllers.NewPostsController(db)
    viewController := controllers.NewViewController(db)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.GET("/api/posts/:slug", postsController.GetPost)
	// router.GET("/api/posts", getAllPosts)
	// router.POST("/api/posts", createPost)
	// router.PUT("/api/posts/:slug", updatePost)
	// router.DELETE("/api/posts/:slug", deletePost)

	router.GET("/", viewController.GetIndexPage)
	router.GET("/posts/:slug", viewController.GetPostPage)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}
