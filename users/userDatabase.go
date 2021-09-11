package users

import (
	"fmt"
	"homework4-backend-go/database"
	"homework4-backend-go/helper"
)

func LoginDatabase(username string) (string,string) {

	var pass string
	var token string
	err := database.Db.QueryRow("SELECT password,token from users_table where username=$1",username).Scan(&pass,&token)
	if err != nil {
		fmt.Println("",helper.ErrorStruct{
			Error: "There is no such a user!!!",
		})
	}
	return pass,token

}

func SignupDatabase(username, password, token string) error{

	_,err := database.Db.Query("INSERT INTO users_table (username, password, token) values ($1,$2,$3)",username,password,token)
	return err

}
