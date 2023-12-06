package logic

import (
	"ginredis/model"
	"ginredis/mysql"
	"log"
)

func GetAllPosts() (data []*model.Post, err error) {
	data, err = mysql.QueryAllPosts()
	if err != nil {
		log.Println("GetAllPosts error:", err)
		return nil, err
	}
	return data, nil
}

func GetPostById(id string) (data *model.Post, err error) {
	data, err = mysql.QueryPostById(id)
	if err != nil {
		log.Println("mysql query post by id error: ", err)
		return nil, err
	}
	return data, nil
}

func GetPostByPostId(postId string) (data *model.Post, err error) {
	data, err = mysql.QueryPostByPostId(postId)
	if err != nil {
		log.Println("mysql query post by post id error: ", err)
		return nil, err
	}
	return data, nil
}

func UpatePostById(postModel *model.Post, post model.Post) (err error) {
	err = mysql.UpdatePostById(postModel, post)
	return err
}

func DeletePostById(postId string) (err error) {
	err = mysql.DeletePostById(postId)
	if err != nil {
		log.Println("mysql delete post by post id failed err: ", err)
		return err
	}
	return nil
}
