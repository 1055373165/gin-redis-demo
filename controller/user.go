package controller

import (
	"ginredis/gredis"
	"ginredis/logic"
	"log"

	"github.com/gin-gonic/gin"
)

func HandleGetAllUsers(c *gin.Context) {
	data, err := logic.GetAllUsers()
	if err != nil {
		log.Println("logic GetAllUsers failed, err: ", err)
		return
	}
	if err = gredis.SetCacheAllUsers(data); err != nil {
		log.Println("Set Cache Error: ", err)
		return
	}
	ResponseError(c, data)
}

// 根据 id 查询用户
func HandleQueryUserById(c *gin.Context) {
	/*
		1. 生成 key name，查询 redis 是否存在这个 key
		2. 存在会直接返回结果，否则查询数据库
		3. 查询失败返回，查询成功返回并将数据放入 redis 缓存中（设置 15s 有效时间）
	*/
	id := c.Param("id")
	redisKey := gredis.GetCacheKey(id)
	if exist := gredis.ExistKey(redisKey); exist {
		data, _ := gredis.GetCacheUserById(redisKey)
		ResponseSuccess(c, data)
	} else {
		data, err := logic.GetUserById(id)
		if err != nil {
			log.Println("handle query user by id failed, err: ", err)
			return
		}
		// 从数据库中查询成功，设置缓存
		err = gredis.SetCacheUserById(id, data)
		if err != nil {
			log.Println("handle query user by id set cache failed, err: ", err)
			return
		}
		ResponseSuccess(c, data)
	}
}

// 根据用户 id 删除
func HandleDeleteUserById(c *gin.Context) {
	/*
		1. 查询 redis 中是否存在指定 key 的缓存数据；如果存在逻辑删除（负的过期时间），否则直接跳过
		2. 在 mysql 中根据用户 id 删除该用户
	*/
	id := c.Param("id")
	key := gredis.GetCacheKey(id)
	if exist := gredis.ExistKey(key); exist {
		err := gredis.DelCacheUserById(key)
		if err != nil {
			log.Println("handle delte user by id failed, err: ", err)
			return
		}
	}
	// 无论缓存是否命中，mysql中都需要进行删除
	err := logic.DeletePostById(id)
	if err != nil {
		log.Println("mysql delete user failed, err: ", err)
		return
	}
	ResponseSuccess(c, "delete success")
}

// 根据用户 id 进行更新
func HandleUpdateUserById(c *gin.Context) {
	/*
		1. 直接更新数据库
		2. 更新完后删除 redis 中对应的 key
		3. 重新查询并重新写入 redis 缓冲中
	*/
	id := c.PostForm("user_id")
	email := c.PostForm("email")
	err := logic.UpdateUserById(id, email)
	if err != nil {
		log.Println("logic update user error: ", err)
		return
	}
	// 删除 redis 中对应的 key
	key := gredis.GetCacheKey(id)
	if exist := gredis.ExistKey(key); exist {
		// 删除 redis 中该 key 的缓存数据
		gredis.DelCacheUserById(key)
		data, _ := logic.GetUserById(id)
		// 重新查询将最新值加载进 redis 中
		err := gredis.SetCacheUserById(key, data)
		if err != nil {
			log.Println("delete redis data success, but update failed, err: ", err)
			return
		}
		ResponseSuccess(c, "update mysql and redis success")
		return
	}
	ResponseSuccess(c, "update mysql success, don't need update redis")
}
