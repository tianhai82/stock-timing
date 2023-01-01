package tests

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/tianhai82/stock-timing/tiger"
	"github.com/tianhai82/stock-timing/tiger/model"
	"github.com/tianhai82/stock-timing/tiger/openapi"
)

const (
	pk      = "MIICXQIBAAKBgQCd/U+3rcSin5XpOPiPlH4oRiquk1J6uXjUkwnXGUmJqJ8VCtNrTRBUFgoKQEUdAksLMRD0EKv35t7GvqtWmT84kUxQknHUk64ike/+5M/PLOyavd2N7vJNUHkSyNSPtiVxyqpnVCgo2TqF86ImWX5w2rnKR6jVzpK+xQLmmuENkwIDAQABAoGBAJ3IPNP48/dhn4rS/dIO/8ti//9nXCj6kETkMCCkvX+AapfOPwTbauI/PHmuZBerkZy0vPSyrbwf0v7zrxQGak6XPGRYJdmE02uKt0uJWbPp2XjMePKfw5fBN5Q7XLLGnTQ4u2niOM0gtGBWj0/2dDCuN6YbYKsZPHUJ8aqhAxghAkEA0K5Cr3RvJyAA4z6iBTCyC95gWQMQ4H9CPeCQL8pwBIv6oOJJtNSnBSFXQp+/LconnQA5ElxFjwSuwwweai3pAwJBAMHQeAuyt/ZI7FXYBnVIkSEE1KAea4Pp0/5/gNiERRy2DNGJLhwqRcmoNiOLyEMs7Sm+mo5gXVi/30Vhdgm+fDECQFfFg9ziX0IYjucF2AXQ1oJxdRrbVEToocb+5gaD4hu3eKIkq5W4f8uDm301TacHySObDWYwkz01XgBB36UPTFsCQQCwiB5foVg4JnHFOu+e8grmhUzZzvtk+p0SSLZmAAwnO5ZvYEC0fLh2FhXByLcOoKQgCrEiD5nWlWVa/4uREoRxAkBx4DLhz03sYRFJ9Ze12LN+OTspC13MQFab7mMg4Z0QAaMGLrcWoa3HtztHRltTmCMG7bogvbu+IzkLPrdJSOed"
	tigerID = "20151992"
	account = "50287417"
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
	b, _ := json.Marshal(resp)
	log.Printf("resp: %s", string(b))

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
