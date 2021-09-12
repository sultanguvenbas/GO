package post

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"homework4-backend-go/database"
	"homework4-backend-go/helper"
	"time"
)

var dataCreated string

func Setup(p *gin.RouterGroup) {
	p.GET("/getSinglePost/:id", getSinglePost)
	p.POST("/post", post)
	p.GET("/getAllPosts", getAllPosts)
}

func getSinglePost(context *gin.Context) {
	id := context.Param("id")
	fmt.Println("id", id)

	err, row := GetSinglePostDatabase(id)

	if err != nil {
		fmt.Println("Something when you are trying to reach single post!(post.go)", err)
		context.JSON(400, helper.ErrorStruct{
			Error: "Something when you are trying to reach single post !(post.go)",
		})
		return
	}
	context.JSON(200, row)

}

func post(context *gin.Context) {
	body := postContent{}

	token := context.GetHeader("token")

	rowData, err := context.GetRawData()
	if err != nil {
		fmt.Println("Input format is wrong(post.go)", err.Error())
		context.AbortWithStatus(400)
		return

	}
	err = json.Unmarshal(rowData, &body)
	if err != nil {
		fmt.Println("Bad Input(post.go)", err.Error())
		context.AbortWithStatus(400)
		return
	}
	dataCreated = time.Now().Format("2006.01.02 15:04:05")
	rNick := helper.RandomNickname()
	rColor := helper.RandomColor()

	err, data := PostDatabase(token, body.Content, dataCreated, rNick, rColor)
	fmt.Println(data)
	if err != nil {
		fmt.Println("Something wrong when inserting post in post table(post.go)", err.Error())
		context.AbortWithStatus(400)
	}

	rows := database.Db.QueryRow("SELECT post_table.postid,content,datecreated,likes,dislikes, nickname, color from post_table left join post_user_vn_table on post_table.postid = post_user_vn_table.postid where post_table.postid=$1", data)
	if err != nil {
		fmt.Println("ERROR SELECTING POSTS:(post.go)", err.Error())
		fmt.Println("Something wrong when you are reaching post table to new post")
		context.AbortWithStatus(400)
	}
	var postDetail postContent
	err = rows.Scan(&postDetail.PostId, &postDetail.Content, &postDetail.DateCreated, &postDetail.Likes, &postDetail.Dislikes, &postDetail.Nickname, &postDetail.Color)
	if err != nil {
		fmt.Println("We could not reach the post!(post.go)", err.Error())
		context.AbortWithStatus(400)
		return
	}

	context.JSON(200, postDetail)
}

func getAllPosts(context *gin.Context) {
	data, err := getPosts()
	if err != nil {
		fmt.Println("Something wrong when getting getposts(post.go)", err)
		context.AbortWithStatus(400)
	}
	context.JSON(200, data)
}

func getPosts() ([]postContent, error) {
	rows, err := database.Db.Query("SELECT post_table.postid,content,datecreated,likes,dislikes,nickname,color from post_table left join post_user_vn_table on post_table.postid = post_user_vn_table.postid order by datecreated desc ")
	if err != nil {
		fmt.Println("Something wrong when you are reaching post table where getPosts(post.go)", err)
		panic(err)
	}

	defer rows.Close()

	var posts []postContent

	for rows.Next() {
		var p postContent

		if err := rows.Scan(&p.PostId, &p.Content, &p.DateCreated, &p.Likes, &p.Dislikes, &p.Nickname, &p.Color); err != nil {
			panic(err)
			return posts, err
		}

		posts = append(posts, p)
	}
	if err = rows.Err(); err != nil {
		return posts, err
	}
	return posts, nil
}
