package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"webapp-go/webapp/models"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")

		_, err := db.NewCreateTable().
			Model((*models.User)(nil)).
			Exec(ctx)
		if err != nil {
			panic(err)
		}

		_, err = db.NewCreateTable().
			Model((*models.Post)(nil)).
			ForeignKey(`("author_id") REFERENCES "users" ("id") ON DELETE CASCADE`).
			Exec(ctx)
		if err != nil {
			panic(err)
		}

		_, err = db.NewCreateTable().
			Model((*models.Document)(nil)).
			ForeignKey(`("post_slug") REFERENCES "posts" ("slug") ON DELETE CASCADE`).
			Exec(ctx)
		if err != nil {
			panic(err)
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")

		_, err := db.NewDropTable().
			Model((*models.User)(nil)).
			IfExists().
			Exec(ctx)
		if err != nil {
			panic(err)
		}

		_, err = db.NewDropTable().
			Model((*models.Post)(nil)).
			IfExists().
			Exec(ctx)
		if err != nil {
			panic(err)
		}

		_, err = db.NewDropTable().
			Model((*models.Document)(nil)).
			IfExists().
			Exec(ctx)
		if err != nil {
			panic(err)
		}

		return nil
	})
}
