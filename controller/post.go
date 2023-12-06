package controller

import (
	"fmt"
	"ginredis/gredis"
	"ginredis/logic"
	"ginredis/model"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
获取全部的 post
1. 如果 redis 中 key 不存在，查询 mysql -> 返回 json -> 将返回结果存到 redis（kv 存储）
2. 如果 redis 中 key 存在，取出缓存直接返回（json数据）
*/
func HandleGetAllPost(c *gin.Context) {
	var data []*model.Post
	var err error
	data, err = logic.GetAllPosts()
	if err != nil {
		log.Println("logic GetALLPost error: ", err)
		ResponseError(c, data)
	}
}

// 根据 post id 获取 post
func HandleGetPostByPostId(c *gin.Context) {
	postId := c.Param("postid")
	data, err := logic.GetPostByPostId(postId)
	if err != nil {
		log.Println("logic GetPostByPostId error", err)
		ResponseError(c, data)
	}
	ResponseSuccess(c, data)
	// 设置缓存
	if err = gredis.SetCachePostById(data, postId); err != nil {
		log.Println("Set Cache Post By Post Id Error: ", err)
	}
}

/*
根据 id 获取 post 的信息
*/
func HandleGetPostById(c *gin.Context) {
	/*
		router.GET("/user/:id", func(c *gin.Context) {
		    a GET request to /user/john
		    id := c.Param("id") // id == "/john"
		    a GET request to /user/john/
		    id := c.Param("id") // id == "/john/"
		})
	*/
	id := c.Param("id")
	data, err := logic.GetPostById(id)
	if err != nil {
		log.Println("logic GetPostById error", err)
		ResponseError(c, data)
		return
	}
	ResponseSuccess(c, data)
}

/*
修改文章内容
*/
func HandleUpdatePostById(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("ID"))
	title, content := c.PostForm("Title"), c.PostForm("Content")
	currentPost := model.Post{Id: id, Title: title, Content: content}
	if err := logic.UpatePostById(&currentPost, currentPost); err != nil {
		log.Println("logic update post by id error", err)
		ResponseError(c, "update error")
		return
	}
	ResponseSuccess(c, "success")

	// 修改完成后更新缓存
	data, err := logic.GetPostById(c.PostForm("ID"))
	if err != nil {
		log.Println("logic query post by id error: ", err)
		return
	}
	currentKey := fmt.Sprintf("%s%s", gredis.KeyPostIdSet, data.PostId)
	gredis.UpdatePost(currentKey)
}

// 删除文章内容
func HandleDeletePostById(c *gin.Context) {
	id := c.Param("id")
	// 查询 mysql 中的 post_id
	data, err := logic.GetPostById(id)
	if err != nil {
		log.Println("logic query PostById error:", err)
		return
	}
	if err = logic.DeletePostById(id); err != nil {
		log.Println("logic delete PostById Error: ", err)
		ResponseError(c, "Error")
		return
	}
	ResponseSuccess(c, "Success")
	// mysql 更新完成后，缓存必须更新
	postId := data.PostId
	// 根据 postId 组装出 redis Key
	currentKey := fmt.Sprintf("%s%s", gredis.KeyPostIdSet, postId)
	// 设置一个 Now 之前的过期时间 Now().Add(-10*time.Second) 就相当于逻辑删除了
	if err := gredis.SetKeyExpired(currentKey); err != nil {
		log.Println("gredis SetKeyExpired Error: ", err)
		return
	}
	log.Println("gredis cache delete success")
}
