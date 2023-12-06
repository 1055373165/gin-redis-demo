package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func ResponseError(c *gin.Context, data interface{}) {
	c.JSON(http.StatusNotFound, gin.H{
		"data": data,
	})
}
