package initializers

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func errorCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetConnection() *gorm.DB {
	myEnv, err := godotenv.Read(".env")
	errorCheck(err)
	host := myEnv["HOST"]
	dbPort := myEnv["DBPORT"]
	user := myEnv["USERNAME"]
	password := myEnv["PASSWORD"]
	dbName := myEnv["NAME"]

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, password, dbPort)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbURI,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	errorCheck(err)
	return db
}

func CloseConnection(db *gorm.DB) {
	db1, err := db.DB()
	errorCheck(err)
	err = db1.Close()
	errorCheck(err)
}
