package rpcs

import (
	"context"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
)

var config = &firebase.Config{
	StorageBucket: "<BUCKET_NAME>.appspot.com",
}

func AddRpcs(router *gin.RouterGroup) {
	router.POST("/instruments", func(c *gin.Context) {
		app, err := firebase.NewApp(context.Background(), config)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}

		client, err := app.Storage(context.Background())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
		bucket, err := client.DefaultBucket()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
		reader, err := bucket.Object("etoro_stocks.json").NewReader(context.Background())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.DataFromReader(200, -1, "application/json", reader, make(map[string]string))
	})
}
