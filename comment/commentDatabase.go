package comment

import (
	"fmt"
	"homework4-backend-go/database"
)

func saveComment(postid int, token, content, nickname, commentDate, color string) (int, error) {
	var commentId int

	_, err := database.Db.Exec("INSERT INTO post_user_vn_table (postid, nickname, color, userid) values ($1,$2,$3,(select userid from users_table where token=$4)) ON CONFLICT DO NOTHING", postid, nickname, color, token)
	if err != nil {
		return 0, err
	}

	err = database.Db.QueryRow("insert into comment_table(postid,userid,content,likes,dislikes,commentdate) values($1,(SELECT userid FROM users_table WHERE token=$2),$3,0,0,$4) returning commentid", postid, token, content, commentDate).Scan(&commentId)
	if err != nil {
		fmt.Println("Error516: when you are inserting comment(commentDatabase.go)", err.Error())
		return 0, err
	}

	return commentId, err
}
