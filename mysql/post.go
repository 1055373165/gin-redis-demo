package mysql

import (
	"ginredis/model"
	"log"
)

// 查询所有的 Post 数据
func QueryAllPosts() (posts []*model.Post, err error) {
	// 1.指定您想要运行数据库操作的表
	// 2.Find 查找与给定条件 conds 匹配的所有记录
	err = db.Table("post").Find(&posts).Error
	if err != nil {
		log.Println("QueryAllPosts error:", err)
		return nil, err
	}
	return posts, nil
}

// 根据 postid 查询数据
func QueryPostByPostId(postId string) (post *model.Post, err error) {
	if err = db.Table("post").Where("post_id = ?", postId).First(&post).Error; err != nil {
		return nil, err
	}
	return post, err
}

// 根据主键 id 查询
func QueryPostById(id string) (post *model.Post, err error) {
	if err = db.Table("post").Where("id = ?", id).Find(&post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func UpdatePostById(postModel *model.Post, post model.Post) (err error) {
	if err = db.Table("post").Model(&postModel).Updates(post).Error; err != nil {
		return err
	}
	return nil
}

func DeletePostById(id string) (err error) {
	var post model.Post
	if err = db.Table("post").Delete(&post, id).Error; err != nil {
		log.Println("mysql delete post id error:", err)
		return
	}
	return nil
}
