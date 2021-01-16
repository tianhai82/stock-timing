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
	"golang.org/x/net/publicsuffix"
)

var httpClient = &http.Client{Timeout: 15 * time.Second}

func init() {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
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
	request.Header.Set("sec-fetch-dest", "document")
	request.Header.Set("sec-fetch-mode", "navigate")
	request.Header.Set("sec-fetch-site", "none")
	request.Header.Set("upgrade-insecure-requests", "1")
	request.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36")
	// request.Header.Set("Accept-Encoding", "gzip")

	// b, _ := httputil.DumpRequest(request, false)
	// fmt.Println(string(b))
	return httpClient.Do(request)
}
