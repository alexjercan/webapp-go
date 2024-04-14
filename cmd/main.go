package main

import (
	"fmt"
	"webapp-go/webapp"
	"webapp-go/webapp/config"
	"webapp-go/webapp/controllers"
	"webapp-go/webapp/repositories"
	"webapp-go/webapp/services"
	"webapp-go/webapp/middlewares"

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

    defer db.Close()

    postsRepository := repositories.NewPostsRepository(db)
    usersRepository := repositories.NewUserRepository(db)
    documentsRepository := repositories.NewDocumentsRepository(db)

    authService := services.NewAuthService(cfg)
    usersService := services.NewUsersService(usersRepository)
    bearerService := services.NewBearerService(cfg)

    postsController := controllers.NewPostsController(postsRepository)
    viewController := controllers.NewViewController(postsRepository, usersRepository)
    authController := controllers.NewAuthController(cfg, authService, usersService, bearerService)
    documentsController := controllers.NewDocumentsController(documentsRepository, postsRepository)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	store := cookie.NewStore([]byte(cfg.AuthStore.Secret))
	router.Use(sessions.Sessions(cfg.AuthStore.Name, store))

	router.LoadHTMLGlob("templates/*")

    authorized := router.Group("/", middlewares.AuthRequired(bearerService))

	authorized.GET("/api/posts/:slug", postsController.GetPost)
	authorized.GET("/api/posts", postsController.GetPosts)
	authorized.POST("/api/posts", postsController.CreatePost)
	authorized.PUT("/api/posts/:slug", postsController.UpdatePost)
	authorized.DELETE("/api/posts/:slug", postsController.DeletePost)

    authorized.GET("/api/posts/:slug/documents/:id", documentsController.GetDocument)
    authorized.GET("/api/posts/:slug/documents", documentsController.GetDocuments)
    authorized.POST("/api/posts/:slug/documents", documentsController.CreateDocument)
    authorized.PUT("/api/posts/:slug/documents/:id", documentsController.UpdateDocument)
    authorized.DELETE("/api/posts/:slug/documents/:id", documentsController.DeleteDocument)

    authorized.GET("/api/user", authController.GetUser)

    authorized.GET("/api/bearer", authController.BearerToken)

	authorized.GET("/", viewController.GetIndexPage)
	authorized.GET("/user", viewController.GetUserPage)
	authorized.GET("/posts/:slug", viewController.GetPostPage)
	authorized.GET("/create", viewController.GetCreatePostPage)

    router.GET("/auth/login", authController.Login)
    router.GET("/auth/callback", authController.Callback)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
    router.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
	// router.Run(":3000") for a hard coded port
}
