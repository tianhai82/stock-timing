package httprequester

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	"github.com/pkg/errors"
)

var httpClient = &http.Client{
	Timeout: 15 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	httpClient.Jar = jar
}

func MakeGetRequest(urlStr string, output interface{}) (err error) {
	resp, err := makeRequest(urlStr)
	if err != nil {
		err = errors.Wrap(err, "http get fails")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.StatusCode == 403 {
			resp, err = makeRequest(urlStr)
			if err != nil {
				err = errors.Wrap(err, "http get fails")
				return
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return errors.New(resp.Status)
			}
		} else {
			return errors.New(resp.Status)
		}
	}

	var reader io.ReadCloser
	respEncoding := resp.Header.Get("Content-Encoding")
	switch respEncoding {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = json.NewDecoder(reader).Decode(output)
		if err != nil {
			err = errors.Wrap(err, "json decoding fails")
			return
		}
		defer reader.Close()
	default:
		err = json.NewDecoder(resp.Body).Decode(output)
		if err != nil {
			err = errors.Wrap(err, "json decoding fails")
			return
		}
	}
	return
}

func makeRequest(urlStr string) (*http.Response, error) {
	request, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("cache-control", "no-cache")
	request.Header.Add("sec-ch-ua", `"Google Chrome";v="87", " Not;A Brand";v="99", "Chromium";v="87"`)
	request.Header.Add("sec-fetch-dest", "document")
	request.Header.Add("sec-fetch-mode", "navigate")
	request.Header.Add("sec-fetch-site", "none")
	request.Header.Add("upgrade-insecure-requests", "1")
	request.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36")

	resp, err := httpClient.Do(request)
	// b, _ := httputil.DumpResponse(resp, true)
	// fmt.Println(string(b))
	return resp, err
}
