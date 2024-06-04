package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"text/template"
	"webapp-go/migrations"
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
	"github.com/uptrace/bun/migrate"

	"github.com/urfave/cli/v2"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	app := &cli.App{
		Name:  "TOP G WebApp",
		Usage: "webapp-go is a web application for TOP G Document based RAG",
		Commands: []*cli.Command{
			newDbCommand(cfg),
			newAppCommand(cfg),
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func newDbCommand(cfg config.Config) *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "run migrations stuff",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					ctx := context.Background()
					db := webapp.DBConnection(cfg)

					migrator := migrate.NewMigrator(db, migrations.Migrations)
					return migrator.Init(ctx)
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					ctx := context.Background()
					db := webapp.DBConnection(cfg)

					migrator := migrate.NewMigrator(db, migrations.Migrations)

					group, err := migrator.Migrate(ctx)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no new migrations to run\n")
						return nil
					}

					fmt.Printf("migrated to %s\n", group)
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					ctx := context.Background()
					db := webapp.DBConnection(cfg)

					migrator := migrate.NewMigrator(db, migrations.Migrations)

					group, err := migrator.Rollback(ctx)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no groups to roll back\n")
						return nil
					}

					fmt.Printf("rolled back %s\n", group)
					return nil
				},
			},
			{
				Name:  "lock",
				Usage: "lock migrations",
				Action: func(c *cli.Context) error {
					ctx := context.Background()
					db := webapp.DBConnection(cfg)

					migrator := migrate.NewMigrator(db, migrations.Migrations)
					return migrator.Lock(ctx)
				},
			},
			{
				Name:  "unlock",
				Usage: "unlock migrations",
				Action: func(c *cli.Context) error {
					ctx := context.Background()
					db := webapp.DBConnection(cfg)

					migrator := migrate.NewMigrator(db, migrations.Migrations)
					return migrator.Unlock(ctx)
				},
			},
			{
				Name:  "create_go",
				Usage: "create Go migration",
				Action: func(c *cli.Context) error {
					ctx := context.Background()
					db := webapp.DBConnection(cfg)

					migrator := migrate.NewMigrator(db, migrations.Migrations)

					name := strings.Join(c.Args().Slice(), "_")
					mf, err := migrator.CreateGoMigration(ctx, name)
					if err != nil {
						return err
					}
					fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)

					return nil
				},
			},
			{
				Name:  "create_sql",
				Usage: "create up and down SQL migrations",
				Action: func(c *cli.Context) error {
					ctx := context.Background()
					db := webapp.DBConnection(cfg)

					migrator := migrate.NewMigrator(db, migrations.Migrations)

					name := strings.Join(c.Args().Slice(), "_")
					files, err := migrator.CreateSQLMigrations(ctx, name)
					if err != nil {
						return err
					}

					for _, mf := range files {
						fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
					}

					return nil
				},
			},
			{
				Name:  "status",
				Usage: "print migrations status",
				Action: func(c *cli.Context) error {
					ctx := context.Background()
					db := webapp.DBConnection(cfg)

					migrator := migrate.NewMigrator(db, migrations.Migrations)

					ms, err := migrator.MigrationsWithStatus(ctx)
					if err != nil {
						return err
					}
					fmt.Printf("migrations: %s\n", ms)
					fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
					fmt.Printf("last migration group: %s\n", ms.LastGroup())

					return nil
				},
			},
			{
				Name:  "mark_applied",
				Usage: "mark migrations as applied without actually running them",
				Action: func(c *cli.Context) error {
					ctx := context.Background()
					db := webapp.DBConnection(cfg)

					migrator := migrate.NewMigrator(db, migrations.Migrations)

					group, err := migrator.Migrate(ctx, migrate.WithNopMigration())
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no new migrations to mark as applied\n")
						return nil
					}

					fmt.Printf("marked as applied %s\n", group)
					return nil
				},
			},
		},
	}
}

func newAppCommand(cfg config.Config) *cli.Command {
	return &cli.Command{
		Name:  "app",
		Usage: "application commands",
		Subcommands: []*cli.Command{
			{
				Name:  "run",
				Usage: "run the application",
				Action: func(_ *cli.Context) error {
					return runApp(cfg)
				},
			},
		},
	}
}

func runApp(cfg config.Config) error {
	ctx := context.Background()

	db := webapp.DBConnection(cfg)

	defer db.Close()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	llm, err := ollama.New(ollama.WithServerURL(cfg.Ollama.Url), ollama.WithModel(cfg.Ollama.Model))
	if err != nil {
		return err
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

	postsController := controllers.NewPostsController(postsRepository, usersRepository)
	viewController := controllers.NewViewController(postsRepository, usersRepository, documentsRepository, embeddingsService)
	authController := controllers.NewAuthController(cfg, authService, usersService, bearerService)
	documentsController := controllers.NewDocumentsController(documentsRepository, postsRepository, documentChan)
	embeddingsController := controllers.NewEmbeddingsController(documentsRepository, embeddingsService)

	go embeddingsService.Worker(ctx)

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

	authorized.GET("/api/search/:slug", embeddingsController.GetSearchResult)

	authorized.GET("/api/user", authController.GetUser)
	authorized.GET("/api/bearer", authController.BearerToken)

	router.GET("/", viewController.GetIndexPage)
	router.GET("/login", viewController.GetLoginPage)
	authorized.GET("/home", viewController.GetHomePage)
	authorized.GET("/user", viewController.GetUserPage)
	authorized.GET("/posts/:slug", viewController.GetPostPage)
	authorized.GET("/create", viewController.GetCreatePostPage)
	authorized.GET("/search/:slug", viewController.SearchPost)

	router.GET("/auth/anonymous", authController.Anonymous)
	router.GET("/auth/login", authController.Login)
	router.GET("/auth/callback", authController.Callback)
	authorized.GET("/auth/logout", authController.Logout)

	router.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))

	return nil
}
