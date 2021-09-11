package database

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"homework4-backend-go/helper"
	"os"
	"strconv"
)

const (
	host     = "localhost"
	port     = 5432
)

var Db *sql.DB //I am gonna use when insert or delete etc something in database
func ConnectDatabase()  {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("", helper.ErrorStruct{
			Error: "loading .env file",
		})
		panic(err)
	}

	host := os.Getenv("host")
	port, _ := strconv.Atoi(os.Getenv("port")) // don't forget to conver int since port is int type.
	user := os.Getenv("user")
	dbname := os.Getenv("dbname")
	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)

	//psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	//	"password=%s dbname=%s sslmode=disable",
	//	host, port, os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"))

	db, err := sql.Open("postgres", psqlSetup)
	if err != nil {
		fmt.Println("",helper.ErrorStruct{Error: "Something wrong when connecting database"})
		panic(err)
	}
	Db=db

	fmt.Println("Successfully connected! Congrats")
}