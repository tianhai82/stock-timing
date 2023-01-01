package tiger

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/tianhai82/stock-timing/tiger/model"
	"github.com/tianhai82/stock-timing/tiger/openapi"
)

const (
	__VERSION__  = "2.2.2"
	USER_LICENSE = "user_license"
)

func GetMacAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	var macAddress string
	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && !bytes.Equal(i.HardwareAddr, nil) {
			macAddress = i.HardwareAddr.String()
			break
		}
	}
	return macAddress
}

var COMMON_PARAM_KEYS = []string{model.P_TIMESTAMP, model.P_TIGER_ID, model.P_METHOD, model.P_CHARSET, model.P_VERSION, model.P_SIGN_TYPE, model.P_DEVICE_ID}

type TigerOpenClient struct {
	config       *TigerOpenClientConfig
	logger       *log.Logger
	headers      map[string]string
	deviceID     string
	license      string
	serverInfo   map[string]string
	serverStatus string
}

func NewTigerOpenClient(clientConfig *TigerOpenClientConfig, logger *log.Logger) *TigerOpenClient {
	client := &TigerOpenClient{config: clientConfig, logger: logger}
	client.headers = map[string]string{
		"Content-type":  "application/json;charset=" + client.config.Charset,
		"Cache-Control": "no-cache",
		"Connection":    "Keep-Alive",
		"User-Agent":    "openapi-python-sdk-" + __VERSION__,
	}
	client.deviceID = GetMacAddress()
	client.initLicense()
	client.refreshServerInfo()
	return client
}

func (c *TigerOpenClient) initLicense() {
	if c.config.License == "" && c.config.EnableDynamicDomain {
		c.config.License = c.queryLicense()
	}
}

func (c *TigerOpenClient) queryLicense() string {
	request := &openapi.OpenApiRequest{Method: USER_LICENSE}

	var responseContent map[string]interface{}
	var err error
	responseContent, err = c.Execute(request, "")
	if err != nil {
		c.logger.Panic(err)
	}

	if responseContent != nil {
		response := &openapi.TigerResponse{}
		response.ParseResponseContent(responseContent)
		if response.IsSuccess() {
			data := make(map[string]string)
			if err := json.Unmarshal([]byte(response.Data), &data); err != nil {
				return ""
			}
			return data["license"]
		}
	}
	c.logger.Panicf("failed to query license, response: %s", responseContent)
	return ""
}

func (c *TigerOpenClient) refreshServerInfo() {
	c.config.refreshServerInfo()
}

func (c *TigerOpenClient) getCommonParams(params map[string]string) map[string]string {
	commonParams := make(map[string]string)
	commonParams[model.P_TIMESTAMP] = params[model.P_TIMESTAMP]
	commonParams[model.P_TIGER_ID] = c.config.TigerID
	commonParams[model.P_METHOD] = params[model.P_METHOD]
	commonParams[model.P_CHARSET] = c.config.Charset
	commonParams[model.P_VERSION] = params[model.P_VERSION]
	if commonParams[model.P_VERSION] == "" {
		commonParams[model.P_VERSION] = model.OPEN_API_SERVICE_VERSION
	}
	commonParams[model.P_SIGN_TYPE] = c.config.SignType
	commonParams[model.P_DEVICE_ID] = c.deviceID
	if url, ok := params[model.P_NOTIFY_URL]; ok {
		commonParams[model.P_NOTIFY_URL] = url
	}
	return commonParams
}

func (c *TigerOpenClient) removeCommonParams(params map[string]string) {
	if params == nil {
		return
	}
	for _, key := range COMMON_PARAM_KEYS {
		if _, ok := params[key]; ok {
			delete(params, key)
		}
	}
}

func (c *TigerOpenClient) prepareRequest(request *openapi.OpenApiRequest) ([]byte, error) {
	request.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	params := request.GetParams()

	commonParams := c.getCommonParams(params)
	allParams := make(map[string]string)
	for key, value := range params {
		allParams[key] = value
	}
	for key, value := range commonParams {
		allParams[key] = value
	}
	signContent := getSignContent(allParams)
	sign, err := signWithRSA(c.config.PrivateKey, []byte(signContent), crypto.SHA1)
	if err != nil {
		c.logger.Println("sign error:", err)
		return nil, nil
	}
	allParams[model.P_SIGN] = sign
	b, err := json.Marshal(allParams)
	if err != nil {
		return nil, fmt.Errorf("fail to json.marshal:%w", err)
	}
	return b, nil
}

func (c *TigerOpenClient) parseResponse(responseStr string, timestamp string) (map[string]interface{}, error) {

	responseContent := make(map[string]interface{})
	if err := json.Unmarshal([]byte(responseStr), &responseContent); err != nil {
		return nil, err
	}

	if c.config.TigerPublicKey == "" || responseContent["sign"] == nil || timestamp == "" {
		return responseContent, nil
	}

	sign := responseContent["sign"].(string)
	err := verifyWithRsa(c.config.TigerPublicKey, []byte(timestamp), []byte(sign))
	if err != nil {
		return nil, fmt.Errorf("[%s]response sign verify failed. %s", err, responseStr)
	}
	return responseContent, err
}

func addStartEnd(key string, startMarker string, endMarker string) string {
	if strings.Index(key, startMarker) < 0 {
		key = startMarker + key
	}
	if strings.Index(key, endMarker) < 0 {
		key = key + endMarker
	}
	return key
}

func fillPublicKeyMarker(publicKey string) string {
	return addStartEnd(publicKey, "-----BEGIN PUBLIC KEY-----\n", "\n-----END PUBLIC KEY-----")
}

func verifyWithRsa(publicKey string, message []byte, sign []byte) error {
	publicKey = fillPublicKeyMarker(publicKey)

	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return fmt.Errorf("failed to parse PEM block")
	}

	sign, err := base64.StdEncoding.DecodeString(string(sign))
	if err != nil {
		return err
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	hashed := sha1.Sum(message)

	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, hashed[:], sign)
}

func (c *TigerOpenClient) Execute(request *openapi.OpenApiRequest, url string) (map[string]interface{}, error) {
	if url == "" {
		url = c.config.ServerURL
	}
	body, err := c.prepareRequest(request)
	if err != nil {
		return nil, err
	}
	c.logger.Println("request body:", string(body))
	httpClient := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		c.logger.Println("create request error:", err)
		return nil, err
	}
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		c.logger.Println("send request error:", err)
		return nil, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.logger.Println("read response error:", err)
		return nil, err
	}
	params := request.GetParams()
	timestamp := params[model.P_TIMESTAMP]
	parsedResp, err := c.parseResponse(string(responseBody), timestamp)
	if err != nil {
		return nil, err
	}
	if parsedResp["code"].(float64) != 0 {
		return nil, fmt.Errorf("query %s failed: code:%v, message:%v", request.Method, parsedResp["code"], parsedResp["message"])
	}
	return parsedResp, nil
}

func getSignContent(params map[string]string) string {
	signContent := ""
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := params[k]
		signContent += "&" + k + "=" + v
	}
	signContent = signContent[1:]
	return signContent
}

func getQueryString(params map[string]string) string {
	var keys []string
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var queryString string
	for _, key := range keys {
		queryString += key + "=" + params[key] + "&"
	}
	return queryString[:len(queryString)-1]
}

func fillPrivateKeyMarker(privateKey string) string {
	return addStartEnd(privateKey, "-----BEGIN RSA PRIVATE KEY-----\n", "\n-----END RSA PRIVATE KEY-----")
}

func signWithRSA(privateKey string, signContent []byte, hashFunc crypto.Hash) (string, error) {

	privateKey = fillPrivateKeyMarker(privateKey)
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	h := hashFunc.New()
	h.Write(signContent)
	hashed := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, key, hashFunc, hashed)
	if err != nil {
		return "", err
	}

	sign := base64.StdEncoding.EncodeToString(signature)
	return sign, nil
}
