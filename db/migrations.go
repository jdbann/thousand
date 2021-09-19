package db

import "embed"

var (
	// Migrations is the embedded filesystem which contains all the SQL migration
	// files.
	//go:embed migrations/*.sql
	Migrations embed.FS

	// AppMigrationsPath is the path from the application root for the migrations
	// folder.
	AppMigrationsPath = "./db/migrations"

	// FSMigrationsPath is the path from the Migrations filesystem root for the
	// migrations folder.
	FSMigrationsPath = "migrations"
)
