package firebase

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/storage"
	"github.com/tianhai82/stock-timing/model"
)

var config = &firebase.Config{
	StorageBucket: "stock-timing.appspot.com",
}
var app *firebase.App
var AuthClient *auth.Client
var StorageClient *storage.Client
var FirestoreClient *firestore.Client
var Instruments []model.InstrumentDisplayData

func init() {
	var err error
	ctx := context.Background()
	app, err = firebase.NewApp(ctx, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	AuthClient, err = app.Auth(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	StorageClient, err = app.Storage(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	FirestoreClient, err = app.Firestore(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	bucket, err := StorageClient.DefaultBucket()
	if err != nil {
		return
	}
	reader, err := bucket.Object("etoro_stocks.json").NewReader(context.Background())
	dec := json.NewDecoder(reader)
	err = dec.Decode(&Instruments)
	if err != nil {
		return
	}
}
