package migrations

import (
	"context"
	"fmt"
	"webapp-go/webapp/models"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")

		_, err := db.NewCreateTable().
			Model((*models.Document)(nil)).
			Exec(ctx)
		if err != nil {
			panic(err)
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")

        _, err := db.NewDropTable().
			Model((*models.Document)(nil)).
			IfExists().
			Exec(ctx)
		if err != nil {
			panic(err)
		}

		return nil
	})
}
