package middleware

import (
	"fmt"
	"ginredis/controller"
	"ginredis/gredis"

	"github.com/gin-gonic/gin"
)

/*
	redis 缓存逻辑

1. 功能：判断 redis 中是否存在 key，如果存在则返回内存中的缓存数据，调用 c.Abort，不走数据库
2. 否则，调用 c.Next() 放行去查询数据库
*/
func CacheMiddleware(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 判断 key 类型
		// 1. 全量查询（范围查询）
		if key == gredis.CACHE_POSTS || key == gredis.CACHE_USERS {
			if exist := gredis.ExistKey(key); !exist {
				c.Next()
			} else {
				// 按照 key 类型取出缓存
				switch key {
				case gredis.CACHE_POSTS:
					data, _ := gredis.GetCacheAllPosts(key)
					controller.ResponseSuccess(c, data)
				case gredis.CACHE_USERS:
					data, _ := gredis.GetCacheAllUsers(key)
					controller.ResponseSuccess(c, data)
				}
				c.Abort()
			}
		}
		// 2. 单点查询
		if key == gredis.KeyPostIdSet || key == gredis.KeyUserIdSet {
			postId := c.Param("postid")
			currentKey := fmt.Sprintf("%s%s", key, postId)
			if exist := gredis.ExistKey(currentKey); !exist {
				c.Next()
			} else {
				switch key {
				case gredis.KeyPostIdSet:
					data, _ := gredis.GetCachePostById(currentKey)
					controller.ResponseSuccess(c, data)
				case gredis.KeyUserIdSet:
					data, _ := gredis.GetCacheUserById(currentKey)
					controller.ResponseSuccess(c, data)
				}
				c.Abort()
			}
		}
	}
}
