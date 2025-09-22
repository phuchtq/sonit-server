package dbserver

import (
	"os"
	"sonit_server/constant/env"
)

type mySQLServer struct{}

func InitializeMySQLServer() ISQLServer {
	return &mySQLServer{}
}

// GetSQLServer implements ISQLServer.
func (m mySQLServer) GetSQLServer() string {
	return env.MYSQL_SERVER
}

// GetCnnStr implements ISQLServer.
func (m mySQLServer) GetCnnStr() string {
	return os.Getenv(env.MYSQL_DB_CNNSTR)
}
