package main

import (
	"github.com/gin-gonic/gin"
	"homework4-backend-go/comment"
	"homework4-backend-go/database"
	"homework4-backend-go/post"
	"homework4-backend-go/users"
)

func main() {

	router := gin.Default()
	database.ConnectDatabase()

	router.Use(func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "http://localhost:4200")
		context.Header("Access-Control-Allow-Headers", "*")

		if context.Request.Method=="OPTIONS" {
			context.Status(200)
			context.Abort()
		}
	})

	userGroup := router.Group("/user")
	postGroup := router.Group("/post")
	commentGroup := router.Group("/comment")

	users.Setup(userGroup)
	post.Setup(postGroup)
	comment.Setup(commentGroup)

	router.Run(":8080")
}
