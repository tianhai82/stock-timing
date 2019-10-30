package rpcs

import (
	"context"
	"log"

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
			log.Fatalln(err)
		}

		client, err := app.Storage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		bucket, err := client.DefaultBucket()
		if err != nil {
			log.Fatalln(err)
		}
		reader, err := bucket.Object("etoro_stocks.json").NewReader(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		c.DataFromReader(200, -1, "application/json", reader, make(map[string]string))
	})
}
