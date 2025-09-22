package db

import (
	"database/sql"
	"errors"
	"log"
	"sonit_server/constant/noti"

	_ "github.com/lib/pq" // this is necessary!

	db_server "sonit_server/data_access/db_server"
)

// Database connection
func ConnectDB(logger *log.Logger, server db_server.ISQLServer) (*sql.DB, error) {
	// Open database connection
	cnn, err := sql.Open(server.GetSQLServer(), server.GetCnnStr())

	if err != nil {
		logger.Println(noti.DB_CONNECTION_ERR_MSG + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return cnn, nil
}
