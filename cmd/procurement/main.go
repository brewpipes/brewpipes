package main

import (
	"context"
	"os"

	"github.com/brewpipes/brewpipes/cmd"
	"github.com/brewpipes/brewpipes/service/procurement"
)

func main() {
	cmd.Main(run)
}

func run(ctx context.Context) error {
	// Entry point for the independent procurement service application.

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = os.Getenv("POSTGRES_DSN")
	}

	svc := procurement.New(procurement.Config{
		PostgresDSN: dsn,
		SecretKey:   os.Getenv("BREWPIPES_SECRET_KEY"),
	})

	return cmd.RunServices(ctx, nil, svc)
}
