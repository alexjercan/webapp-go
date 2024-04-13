package main

import (
	"fmt"
	"webapp-go/webapp"
	"webapp-go/webapp/config"
	"webapp-go/webapp/controllers"
	"webapp-go/webapp/repositories"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        panic(err)
    }

    db := webapp.DBConnection(cfg)

    postsRepository := repositories.NewPostsRepository(db)
    usersRepository := repositories.NewUserRepository(db)

    postsController := controllers.NewPostsController(postsRepository)
    viewController := controllers.NewViewController(postsRepository)
    authController := controllers.NewAuthController(cfg, usersRepository)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	store := cookie.NewStore([]byte(cfg.AuthStore.Secret))
	router.Use(sessions.Sessions(cfg.AuthStore.Name, store))

	router.LoadHTMLGlob("templates/*")

	router.GET("/api/posts/:slug", postsController.GetPost)
	router.GET("/api/posts", postsController.GetPosts)
	router.POST("/api/posts", postsController.CreatePost)
	router.PUT("/api/posts/:slug", postsController.UpdatePost)
	router.DELETE("/api/posts/:slug", postsController.DeletePost)

	router.GET("/", viewController.GetIndexPage)
	router.GET("/user", viewController.GetUserPage)
	router.GET("/posts/:slug", viewController.GetPostPage)

    router.GET("/auth/login", authController.Login)
    router.GET("/auth/callback", authController.Callback)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
    router.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
	// router.Run(":3000") for a hard coded port
}
