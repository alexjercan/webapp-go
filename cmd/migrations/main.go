package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"webapp-go/migrations"
	"webapp-go/webapp"
	"webapp-go/webapp/config"

	"github.com/uptrace/bun/migrate"

	"github.com/urfave/cli/v2"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        panic(err)
    }

	app := &cli.App{
		Name: "migrations",
        Usage: "run migrations stuff",
		Commands: []*cli.Command{
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
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
