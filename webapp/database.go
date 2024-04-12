package webapp

import (
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func DBConnection(cfg Config) *bun.DB {
    pgconn := pgdriver.NewConnector(
        pgdriver.WithAddr(fmt.Sprintf("%s:%d", cfg.Database.Host, cfg.Database.Port)),
        pgdriver.WithUser(cfg.Database.User),
        pgdriver.WithPassword(cfg.Database.Password),
        pgdriver.WithDatabase(cfg.Database.Database),
        pgdriver.WithInsecure(cfg.Database.Insecure),
    )
    sqldb := sql.OpenDB(pgconn)
    db := bun.NewDB(sqldb, pgdialect.New())

    return db
}
