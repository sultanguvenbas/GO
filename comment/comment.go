package comment

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"homework4-backend-go/database"
	"homework4-backend-go/helper"
	"time"
)

var dateCreated string

func Setup(c *gin.RouterGroup) {
	c.POST("/post", postComment)
	c.GET("/getAllComments/:postid", getAllComments)
}

func postComment(context *gin.Context) {
	body := commentStructure{}

	token := context.GetHeader("token")

	rowData, err := context.GetRawData()
	if err != nil {
		fmt.Println("Input format is wrong(comment.go)", err)
		panic(err)
		context.AbortWithStatus(400)
		return
	}

	err = json.Unmarshal(rowData, &body)
	if err != nil {
		fmt.Println("Something wrong when unmarshal part(comment.go)", err)
		context.AbortWithStatus(400)
	}

	dateCreated = time.Now().Format("2006.01.02 15:04:05")

	row, err := database.Db.Query("SELECT nickname from post_user_vn_table where postid=$1", body.PostId)
	if err != nil {
		fmt.Println("Error: when getting nickname from post_user_vn_table(comment.go)", err)
		context.JSON(400, helper.ErrorStruct{Error: "415"})
		return
	}

	var allNicknames []string
	for row.Next() {
		var nick string
		if err := row.Scan(&nick); err != nil {
			fmt.Println("Something wrong when scanning nickname(comment.go)", err)
			context.JSON(400, helper.ErrorStruct{
				Error: "Something wrong when scanning nickname"})
			return
		}
		allNicknames = append(allNicknames, nick)
	}

	var nickname string
	for chekNicknames(&nickname, allNicknames) {
	}

	cmntId, err := saveComment(body.PostId, token, body.Content, nickname, dateCreated, helper.RandomColor())
	if err != nil {
		fmt.Println("ERROR12312: Something wrong when you are saving comment:", err)
		context.AbortWithStatus(400)
		return
	}

	rows := database.Db.QueryRow("select commentid,comment_table.postid,content,likes,dislikes, post_user_vn_table.nickname,commentdate,post_user_vn_table.color from comment_table left join  post_user_vn_table on comment_table.postid=post_user_vn_table.postid and comment_table.userid=post_user_vn_table.userid where commentid=$1", cmntId)
	var comment commentStructure
	err = rows.Scan(&comment.CommentId, &comment.PostId, &comment.Content, &comment.Likes, &comment.Dislikes, &comment.Nickname, &comment.CommentDate, &comment.Color)
	if err != nil {
		fmt.Println("We could not reach the comment!(comment.go)", err)
		context.JSON(400, helper.ErrorStruct{
			Error: "We could not reach the comment!",
		})
		return
	}
	context.JSON(200, comment)

}

func chekNicknames(nickname *string, allNicknames []string) bool {
	nick := helper.RandomNickname()
	*nickname = nick
	for _, b := range allNicknames {
		if b == *nickname {
			return true
		}
	}
	return false
}

func getAllComments(context *gin.Context) {
	postid := context.Param("postid")

	row, err := database.Db.Query("select commentid,comment_table.postid,content,likes,dislikes,nickname,commentdate,color from comment_table left join post_user_vn_table puvt on comment_table.postid = puvt.postid and comment_table.userid=puvt.userid where comment_table.postid=$1", postid)
	if err != nil {
		fmt.Println("Something wrong when reaching comment table for get all comments(comment.go)", err)
		panic(err)
		context.AbortWithStatus(400)
	}

	defer row.Close()

	var comment []commentStructure

	for row.Next() {
		var com commentStructure
		if err := row.Scan(&com.CommentId, &com.PostId, &com.Content, &com.Likes, &com.Dislikes, &com.Nickname, &com.CommentDate, &com.Color); err != nil {
			fmt.Println("Something wrong when scanning comments(comment.go)", err)
			context.AbortWithStatus(400)
		}
		comment = append(comment, com)
	}
	context.JSON(200, comment)

}
