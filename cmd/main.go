package main

import (
	"context"
	"errors"
	"fmt"
	"text/template"
	"webapp-go/webapp"
	"webapp-go/webapp/config"
	"webapp-go/webapp/controllers"
	"webapp-go/webapp/middlewares"
	"webapp-go/webapp/models"
	"webapp-go/webapp/repositories"
	"webapp-go/webapp/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db := webapp.DBConnection(cfg)

	defer db.Close()

	llm, err := ollama.New(ollama.WithModel("llama2"))
	if err != nil {
		panic(err)
	}

	documentChan := make(chan models.DocumentChanItem, 128)

	postsRepository := repositories.NewPostsRepository(db)
	usersRepository := repositories.NewUserRepository(db)
	documentsRepository := repositories.NewDocumentsRepository(db)
    embeddingRepository := repositories.NewEmbeddingsRepository(db)

	authService := services.NewAuthService(cfg)
	usersService := services.NewUsersService(usersRepository)
	bearerService := services.NewBearerService(cfg)
	embeddingsService := services.NewEmbeddingsService(documentsRepository, embeddingRepository, llm, documentChan)

	postsController := controllers.NewPostsController(postsRepository)
	viewController := controllers.NewViewController(postsRepository, usersRepository)
	authController := controllers.NewAuthController(cfg, authService, usersService, bearerService)
	documentsController := controllers.NewDocumentsController(documentsRepository, postsRepository, documentChan)
    embeddingsController := controllers.NewEmbeddingsController(documentsRepository, embeddingsService)

	go embeddingsService.Worker(ctx)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	store := cookie.NewStore([]byte(cfg.AuthStore.Secret))
	router.Use(sessions.Sessions(cfg.AuthStore.Name, store))

	router.SetFuncMap(template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	})
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

    authorized.GET("/api/embeddings/:slug", embeddingsController.GetSimilarDocument)

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
