package post

import (
	"fmt"
	"homework4-backend-go/database"
)

func GetSinglePostDatabase(id string) (error, postContent) {
	row := database.Db.QueryRow("SELECT post_table.postid,content,datecreated,likes,dislikes,nickname,color from post_table left join post_user_vn_table on post_table.postid = post_user_vn_table.postid  WHERE post_table.postid=$1", id)
	var post postContent
	err := row.Scan(&post.PostId, &post.Content, &post.DateCreated, &post.Likes, &post.Dislikes, &post.Nickname, &post.Color)

	return err, post
}

func PostDatabase(token, content, dataCreated, rnNick, rnColor string) (error, int) {
	var postId int
	err := database.Db.QueryRow("with pInfo as (insert into post_table (userid, content, datecreated, likes, dislikes) values ((select userid from users_table where token=$1),$2,$3,0,0)returning postid,userid) insert into post_user_vn_table(postid, userid, nickname, color) values((SELECT postid from pInfo),(SELECT userid from pInfo),$5,$4)  returning postid", token, content, dataCreated, rnColor, rnNick).Scan(&postId)
	if err != nil {
		fmt.Println("Error: when you are inserting post", err.Error())
		return err, 0
	}
	err = database.Db.QueryRow("INSERT INTO post_user_vn_table (postid, nickname, color, userid) values ($1,$2,$3,(select userid from users_table where token=$4)) returning postid", postId, rnNick, rnColor, token).Scan(&postId)
	return err, postId
}