package webapp

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func DBConnection() *bun.DB {
    pgconn := pgdriver.NewConnector(
        pgdriver.WithAddr("localhost:5432"),
        pgdriver.WithUser("postgres"),
        pgdriver.WithPassword("mysecretpassword"),
        pgdriver.WithDatabase("postgres"),
        pgdriver.WithInsecure(true),
    )
    sqldb := sql.OpenDB(pgconn)
    db := bun.NewDB(sqldb, pgdialect.New())

    return db
}
