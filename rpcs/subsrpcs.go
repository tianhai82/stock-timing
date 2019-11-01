package rpcs

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tianhai82/stock-timing/model"
)

func AddSubscriptionRpcs(router *gin.RouterGroup) {
	router.POST("/subscribe", subscribe)
}

func subscribe(c *gin.Context) {
	user, exist := c.Get("loginUser")
	if !exist {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	loginUser := user.(model.UserAccount)
	fmt.Println(loginUser.Email)
	c.JSON(200, "Implementing soon!!")
}
