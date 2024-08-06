package db

import (
	"database/sql"
	"fmt"
	"log"
	"github.com/Shoetan/pkg/utils"
	_ "github.com/lib/pq"
)


func Database() (*sql.DB, error)  {
	
	user := utils.GetEnvVariable("POSTGRES_USER")
	password := utils.GetEnvVariable("POSTGRES_PASSWORD")
	dbname := utils.GetEnvVariable("POSTGRES_DB")

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", "localhost", user, password, dbname, "disable") 

		db, err := sql.Open("postgres", connectionString)

		if err != nil {
			log.Fatalln(err.Error())
		}

		if err := db.Ping() ; err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Connection to database was succesful 👍 ")
		}
		return db, err

}