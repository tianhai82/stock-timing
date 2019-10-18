package httprequester

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

func MakeGetRequest(url string, output interface{}) (err error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		err = errors.Wrap(err, "http get fails")
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(output)
	if err != nil {
		err = errors.Wrap(err, "json decoding fails")
		return
	}
	return
}
