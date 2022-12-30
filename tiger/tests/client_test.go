package tests

import (
	"log"
	"testing"
	"time"

	"github.com/tianhai82/stock-timing/tiger"
	"github.com/tianhai82/stock-timing/tiger/model"
	"github.com/tianhai82/stock-timing/tiger/openapi"
)

const (
	pk      = ""
	tigerID = ""
	account = ""
)

func TestGetQuotePermission(t *testing.T) {
	client := tiger.NewTigerOpenClient(tiger.NewTigerOpenClientConfig(false, false, pk, tigerID, account), log.Default())
	req := &openapi.OpenApiRequest{
		Method: model.GET_QUOTE_PERMISSION,
	}
	resp, err := client.Execute(req, "")
	if err != nil {
		t.Fatalf("execute failed: %s", err.Error())
	}
	log.Printf("resp: %#v", resp)

	req = &openapi.OpenApiRequest{
		Method: model.KLINE,
	}
	params := map[string]interface{}{}
	params["symbols"] = []string{"AAPL"}
	params["period"] = "day"
	params["begin_time"] = -1
	params["end_time"] = -1
	params["right"] = "br"
	params["limit"] = 60
	params["lang"] = tiger.LANGUAGE
	req.BizModel = params
	resp, err = client.Execute(req, "")
	if err != nil {
		t.Fatalf("execute failed: %s", err.Error())
	}
	log.Printf("resp: %#v", resp)

	req = &openapi.OpenApiRequest{
		Method: model.OPTION_BRIEF,
	}

	tz, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("load timezone error:%v", err)
	}
	param := []map[string]interface{}{
		{
			"symbol": "AAPL",
			"right":  "PUT",
			"expiry": time.Date(2023, time.January, 20, 0, 0, 0, 0, tz).UnixMilli(),
			"strike": "00134000",
		},
	}
	// params = map[string]interface{}{}
	// params["contract"] = param
	// params["version"] = model.OPEN_API_SERVICE_VERSION_V3
	// params["lang"] = tiger.LANGUAGE
	req.BizModel = param
	resp, err = client.Execute(req, "")
	if err != nil {
		t.Fatalf("execute failed: %s", err.Error())
	}
	log.Printf("resp: %#v", resp)
}
