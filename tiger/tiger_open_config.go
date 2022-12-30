package tiger

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	DEFAULT_DOMAIN         = "openapi.tigerfintech.com"
	DEFAULT_SANDBOX_DOMAIN = "openapi-sandbox.tigerfintech.com"
	DOMAIN_GARDEN_ADDRESS  = "https://cg.play-analytics.com/"

	HTTPS_PROTOCAL   = "https://"
	SSL_PROTOCAL     = "ssl"
	GATEWAY_SUFFIX   = "/gateway"
	DOMAIN_SEPARATOR = "-"

	// 老虎证券开放平台网关地址
	SERVER_URL = HTTPS_PROTOCAL + DEFAULT_DOMAIN + GATEWAY_SUFFIX

	// 老虎证券开放平台公钥
	TIGER_PUBLIC_KEY = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDNF3G8SoEcCZh2rshUbayDgLLrj6rKgzNMxDL2HS" +
		"nKcB0+GPOsndqSv+a4IBu9+I3fyBp5hkyMMG2+AXugd9pMpy6VxJxlNjhX1MYbNTZJUT4nudki4uh+LM" +
		"OkIBHOceGNXjgB+cXqmlUnjlqha/HgboeHSnSgpM3dKSJQlIOsDwIDAQAB"
	// 请求签名类型
	SIGN_TYPE = "RSA"
	// 请求字符集
	CHARSET = "UTF-8"
	// 语言
	LANGUAGE = "en_US"
	// 请求超时时间, 单位秒, 默认15s
	TIMEOUT = 15

	// sandbox 环境配置
	SANDBOX_SERVER_URL = HTTPS_PROTOCAL + DEFAULT_SANDBOX_DOMAIN + GATEWAY_SUFFIX

	SANDBOX_TIGER_PUBLIC_KEY = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCbm21i11hgAENGd3/f280PSe4g9YGkS3TEXBY" +
		"MidihTvHHf+tJ0PYD0o3PruI0hl3qhEjHTAxb75T5YD3SGK4IBhHn/Rk6mhqlGgI+bBrBVYaXixm" +
		"HfRo75RpUUuWACyeqQkZckgR0McxuW9xRMIa2cXZOoL1E4SL4lXKGhKoWbwIDAQAB"
)

type ServiceType string

const (
	COMMON ServiceType = "COMMON"
	TRADE              = "TRADE"
	QUOTE              = "QUOTE"
)

type TigerOpenClientConfig struct {
	TigerID             string
	Account             string
	IsPaper             bool
	PrivateKey          string
	License             string
	SignType            string
	SecretKey           string
	Charset             string
	Language            string
	Timezone            string
	Timeout             int
	SandboxDebug        bool
	TigerPublicKey      string
	ServerURL           string
	QuoteServerURL      string
	DomainConf          map[string]string
	EnableDynamicDomain bool
}

func NewTigerOpenClientConfig(sandboxDebug bool, enableDynamicDomain bool, pk, tigerID, account string) *TigerOpenClientConfig {
	return &TigerOpenClientConfig{
		SandboxDebug:        sandboxDebug,
		EnableDynamicDomain: enableDynamicDomain,
		SignType:            SIGN_TYPE,
		Charset:             CHARSET,
		Language:            LANGUAGE,
		Timeout:             TIMEOUT,
		PrivateKey:          pk,
		TigerID:             tigerID,
		Account:             account,
		TigerPublicKey:      TIGER_PUBLIC_KEY,
		ServerURL:           SERVER_URL,
		QuoteServerURL:      SERVER_URL,
	}
}

func (c *TigerOpenClientConfig) Init() {
	if c.SandboxDebug {
		c.TigerPublicKey = SANDBOX_TIGER_PUBLIC_KEY
		c.ServerURL = SANDBOX_SERVER_URL
		c.QuoteServerURL = SANDBOX_SERVER_URL
	}

	if c.EnableDynamicDomain {
		c.DomainConf = c.queryDomains()
		c.refreshServerInfo()
	}
}
func (c *TigerOpenClientConfig) refreshServerInfo() {
	if c.EnableDynamicDomain && len(c.DomainConf) > 0 {
		if c.License != "" {
			c.ServerURL = c.getDomainByType(TRADE, c.License, c.IsPaper) + GATEWAY_SUFFIX
			c.QuoteServerURL = c.getDomainByType(QUOTE, c.License, c.IsPaper) + GATEWAY_SUFFIX
		}
	}
}

func (c *TigerOpenClientConfig) queryDomains() map[string]string {
	result := make(map[string]string)

	resp, err := http.Get(DOMAIN_GARDEN_ADDRESS)
	if err != nil {
		return result
	}
	defer resp.Body.Close()

	items := struct {
		Items []struct {
			Openapi        map[string]string `json:"openapi"`
			OpenapiSandbox map[string]string `json:"openapi-sandbox"`
		} `json:"items"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return result
	}

	for _, item := range items.Items {
		if c.SandboxDebug {
			result = item.OpenapiSandbox
		} else {
			result = item.Openapi
		}
		if len(result) > 0 {
			break
		}
	}

	return result
}

func (c *TigerOpenClientConfig) getDomainByType(serviceType ServiceType, license string, isPaper bool) string {
	licenseValue := license
	commonDomain := c.DomainConf[string(COMMON)]
	if serviceType != COMMON {
		if serviceType == QUOTE {
			return c.DomainConf[fmt.Sprintf("%s-QUOTE", licenseValue)]
		}
		if isPaper {
			return c.DomainConf[fmt.Sprintf("%s-PAPER", licenseValue)]
		}
		return c.DomainConf[licenseValue]
	}
	return commonDomain
}
