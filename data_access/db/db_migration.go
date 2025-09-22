package db

import (
	"errors"
	"log"
	action_type "sonit_server/constant/action_type"
	"sonit_server/constant/noti"
	db_server "sonit_server/data_access/db_server"

	"github.com/golang-migrate/migrate"
)

const (
	migration_script_folder_dir string = "file://sql_script/"
)

// Database migration
func MigrateDB(action string, version int, server db_server.ISQLServer, logger *log.Logger) error {
	// Identify action
	var refixDir string
	switch action {
	case action_type.MIGRATION_TYPE:
		refixDir = "migration"
	case action_type.ROLLBACK_TYPE:
		refixDir = "rollback"
	default: // No method found
		return errors.New("Invalid migration command.")
	}

	// Initialize migration
	migration, err := migrate.New(
		migration_script_folder_dir+refixDir,
		server.GetCnnStr(),
	)

	if err != nil {
		logger.Println(noti.DB_MIGRATION_ERR_MSG + err.Error())
		return errors.New(noti.DB_MIGRATION_INFORM_MSG)
	}

	if err := migration.Migrate(uint(version)); err != nil && err != migrate.ErrNoChange {
		logger.Println(noti.DB_CONNECTION_ERR_MSG + err.Error())
		return errors.New(noti.DB_MIGRATION_INFORM_MSG)
	}

	log.Println(action + "successful.")
	return nil
}
