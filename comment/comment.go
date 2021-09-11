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

func Setup(c *gin.RouterGroup)  {
	c.POST("/post",postComment)
	c.GET("/getAllComments/:postid",getAllComments)
}

func postComment(context *gin.Context) {
	body := commentStructure{}

	token := context.GetHeader("token")
	context.JSON(200,token)

	rowData, err := context.GetRawData()
	if err != nil {
		context.JSON(400, helper.ErrorStruct{
			Error: "Input format is wrong",
		})
		return
	}

	err = json.Unmarshal(rowData,&body)
	if err != nil {
		context.JSON(400,helper.ErrorStruct{
			Error: "Something wrong when unmarshal part",
		})
	}

	dateCreated = time.Now().Format("2006.01.02 15:04:05")
	var commentid int

	//checking post is created or not
	nickname := checkPostVnTable(body.PostId,token)

	//checking if owner wants to comment own post
	if nickname != ""{
		err,commentid = postOwnerComment(body.PostId,token,body.Content,dateCreated)
		fmt.Println(commentid)
		if err != nil {
			fmt.Println("Error: when inserting post owner nickname",err)
			context.JSON(400, helper.ErrorStruct{
				Error: "Something wrong when inserting post owner nickname"})
			return
		}
	}else {
		row, err :=database.Db.Query("SELECT nickname from post_user_vn_table where postid=$1",body.PostId)
		if err != nil {
			fmt.Println("Error: when getting nickname from post_user_vn_table",err)
			context.JSON(400, helper.ErrorStruct{
				Error: "Something wrong when getting nickname from post_user_vn_table"})
			return
		}

		var allNicknames[]string
		for row.Next(){
			var nick string
			if err := row.Scan(&nick); err != nil{
				context.JSON(400, helper.ErrorStruct{
					Error: "Something wrong when scanning nickname"})
				return
			}
			allNicknames=append(allNicknames, nick)
		}

		newNickname := helper.RandomNickname()
		for i:=0; i< len(allNicknames); i++{
			if newNickname != allNicknames[i]{
				err = saveComment(body.PostId,token,body.Content,newNickname,dateCreated,helper.RandomColor())
				if err != nil {
					context.JSON(400, helper.ErrorStruct{
						Error: "Something wrong when inserting comment"})
					return
				}
			}
			//if newNickname is equal one of them should I create again randomnickname here ?????
		}
	}

	rows := database.Db.QueryRow("select commentid,comment_table.postid,content,likes,dislikes,nickname,commentdate,color from comment_table left join post_user_vn_table puvt on comment_table.postid = puvt.postid where commentid=$1",commentid)
	var comment commentStructure
	err = rows.Scan(&comment.CommentId,&comment.PostId,&comment.Content,&comment.Likes,&comment.Dislikes,&comment.Nickname,&comment.CommentDate,&comment.Color);
	if err != nil {
		context.JSON(400, helper.ErrorStruct{
			Error: "We could not reach the comment!",
		})
		return
	}
	context.JSON(200,comment)

}




func getAllComments(context *gin.Context) {
	postid :=context.Param("postid")

	row,err := database.Db.	Query("select commentid,comment_table.postid,content,likes,dislikes,nickname,commentdate,color from comment_table left join post_user_vn_table puvt on comment_table.postid = puvt.postid where comment_table.postid=$1",postid)
	if err != nil {
		context.JSON(400,helper.ErrorStruct{Error: "Something wrong when reaching comment table for get all comments"})
		panic(err)
	}

	defer row.Close()

	var comment []commentStructure

	for row.Next(){
		var com commentStructure
		if err:= row.Scan(&com.CommentId,&com.PostId,&com.Content,&com.Likes,&com.Dislikes,&com.Nickname,&com.CommentDate,&com.Color); err != nil{
			context.JSON(400,helper.ErrorStruct{Error: "Something wrong when scanning comments"})
		}
		comment = append(comment,com)
	}
	context.JSON(200,comment)

}























