package comment

import (
	"fmt"
	"homework4-backend-go/database"
)

func checkPostVnTable(postid int, token string)  string {

	var nickname string
	err := database.Db.QueryRow("select nickname from post_user_vn_table where postid=$1 and userid=(select userid from users_table where token=$2)", postid, token).Scan(&nickname)
	if err != nil {
		fmt.Println("Error: when you are checking post vn table ", err.Error())
		return err.Error()
	}
	return nickname
}

func postOwnerComment(postid int, token, content, commentDate string) (error,int) {

	var commentId int
	err := database.Db.QueryRow("insert into comment_table(postid,userid,content,likes,dislikes,commentdate) values($1,(SELECT userid FROM users_table WHERE token=$2),$3,0,0,$4) returning commentid", postid, token, content, commentDate).Scan(&commentId)
	return err,commentId
}
func saveComment(postid int, token, content, nickname, commentDate, color string) error {
	var commentId int

	err := database.Db.QueryRow("insert into comment_table(postid,userid,content,likes,dislikes,commentdate) values($1,(SELECT userid FROM users_table WHERE token=$2),$3,0,0,$4)", postid, token, content, commentDate).Scan(&commentId)
	if err != nil {
		fmt.Println("Error: when you are inserting comment", err.Error())
		return err
	}
	err = database.Db.QueryRow("INSERT INTO post_user_vn_table (postid, nickname, color, userid) values ($1,$2,$3,(select userid from users_table where token=$4)) returning postid", postid, nickname, color, token).Scan(&commentId)

	return err
}
