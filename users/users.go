package users

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"homework4-backend-go/helper"
	"time"
)

func Setup(c *gin.RouterGroup) {
	c.POST("/login", login)
	c.POST("/signup", signup)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
// TODO: encrypt the token so nobody can actually read ıt.
func tokenGenerator() string {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	dt := time.Now().Format("2006-01-02T15:04:05.999999999Z07:00")
	return dt + hex.EncodeToString(b)
}

func signup(context *gin.Context) {
	body := userStruct{}
	rowData, err := context.GetRawData()
	if err != nil {
		context.JSON(400, helper.ErrorStruct{
			Error: "Input format is wrong",
		})
		return

	}
	err = json.Unmarshal(rowData, &body)
	if err != nil {
		context.JSON(400, helper.ErrorStruct{
			Error: "Something wrong when unmarshalling",
		})
		return
	}

	newToken := tokenGenerator()

	if len(body.Password) <= 6 || len(body.Password) > 14 {
		context.JSON(403, helper.ErrorStruct{
			Error: "Password length should be less than 14 and more than 6 please",
		})
		return
	}
	hashPass, _ := HashPassword(body.Password)

	err = SignupDatabase(body.Username, hashPass, newToken)

	if err != nil {
		context.JSON(400,helper.ErrorStruct{
			Error: "User is already exist!!!",
		})
		return
	}
	context.JSON(200, "User is created")
}

func login(context *gin.Context) {
	body := userStruct{}
	rowData, err := context.GetRawData()
	if err != nil {
		context.JSON(400, helper.ErrorStruct{
			Error: "Input format is wrong",
		})
		return

	}
	err = json.Unmarshal(rowData, &body)
	if err != nil {
		context.JSON(400, helper.ErrorStruct{
			Error: "Bad Input",
		})
		return
	}

	pass, token := LoginDatabase(body.Username)

	match := CheckPasswordHash(body.Password, pass)
	if match {
		context.JSON(200, token)
	} else {
		context.JSON(400, helper.ErrorStruct{
			Error: "Something wrong when try to login with password match",
		})
	}
}
