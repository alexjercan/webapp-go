package main

import (
	"fmt"
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
    cfg, err := webapp.LoadConfig()
    if err != nil {
        panic(err)
    }

    db := webapp.DBConnection(cfg)

    postsController := controllers.NewPostsController(db)
    viewController := controllers.NewViewController(db)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.GET("/api/posts/:slug", postsController.GetPost)
	router.GET("/api/posts", postsController.GetPosts)
	router.POST("/api/posts", postsController.CreatePost)
	router.PUT("/api/posts/:slug", postsController.UpdatePost)
	router.DELETE("/api/posts/:slug", postsController.DeletePost)

	router.GET("/", viewController.GetIndexPage)
	router.GET("/posts/:slug", viewController.GetPostPage)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
    router.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
	// router.Run(":3000") for a hard coded port
}
