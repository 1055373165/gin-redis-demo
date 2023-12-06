package main

import (
	"ginredis/controller"
	"ginredis/gredis"
	"ginredis/middleware"
	"ginredis/mysql"

	"github.com/gin-gonic/gin"
)

func main() {
	err := gredis.InitRedis()
	if err != nil {
		panic(err)
	}
	err = mysql.InitMysql()
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	// Get Uesrs Routes Registration
	r.GET("/getalluser", middleware.CacheMiddleware(gredis.CACHE_USERS), controller.HandleGetAllUsers)
	r.GET("/getuserbyid", controller.HandleQueryUserById)
	r.DELETE("/deluserbyid/:id", controller.HandleDeleteUserById)
	r.POST("/updateuserbyid", controller.HandleUpdateUserById)

	// Get Posts Routes Registration
	r.GET("/getallposts", middleware.CacheMiddleware(gredis.CACHE_POSTS), controller.HandleGetAllPost)
	r.GET("/getpostbyid/:id", controller.HandleGetPostById)
	r.GET("/getpostbypostid/:postid", middleware.CacheMiddleware(gredis.KeyPostIdSet), controller.HandleGetPostByPostId)
	r.POST("/updatepostbyid", controller.HandleUpdatePostById)
	r.DELETE("/deletepostbyid", controller.HandleDeletePostById)

	r.Run(":9091")
}
