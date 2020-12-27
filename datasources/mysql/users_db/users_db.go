package users_db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// is being used in sql.Open()
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const (
	mySQLUsersUsername = "MYSQL_USERS_USERNAME"
	mySQLUsersPassword = "MYSQL_USERS_PASSWORD"
	mySQLUsersHost     = "MYSQL_USERS_HOST"
	mySQLUsersSchema   = "MYSQL_USERS_SCHEMA"
)

var (
	// Client the usersDB Client
	Client *sql.DB

	username = goDotEnvVariable(mySQLUsersUsername)
	password = goDotEnvVariable(mySQLUsersPassword)
	host     = goDotEnvVariable(mySQLUsersHost)
	schema   = goDotEnvVariable(mySQLUsersSchema)
)

func init() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", username, password, host, schema)
	var err error

	Client, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	if err = Client.Ping(); err != nil {
		panic(err)
	}

	log.Println("database successfully configured")
}

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
