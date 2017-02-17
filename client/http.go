package client

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	tokenServer       string
	tokenUsername     string
	tokenPassword     string
	payClientCertFile string
	payClientKeyFile  string
	payClientRootCert string

	HTTPC  *HTTPClient
	HTTPSC *HTTPSClient
)

func init() {
	HTTPC = &HTTPClient{}
	HTTPSC = NewHTTPSClient()
}

// HTTPSClient HTTPS客户端结构
type HTTPSClient struct {
	http.Client
}

// GetDefaultClient 返回默认的客户端
func GetDefaultClient() *HTTPSClient {
	return HTTPSC
}

// NewHTTPSClient 新建https客户端
// func NewHTTPSClient(certFile string, keyFile string) *HTTPSClient {
// 	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
// 	if err != nil {
// 		log.Error("load x509 cert error:", err)
// 		return nil
// 	}
// 	config := &tls.Config{
// 		Certificates: []tls.Certificate{cert},
// 	}
// 	config.BuildNameToCertificate()
// 	tr := &http.Transport{TLSClientConfig: config}
// 	client := &http.Client{Transport: tr}
// 	return &HTTPSClient{
// 		Client: *client,
// 	}
// }

// NewHTTPSClient 获取默认https客户端
func NewHTTPSClient() *HTTPSClient {
	config := &tls.Config{InsecureSkipVerify: true}
	tr := &http.Transport{TLSClientConfig: config}
	client := http.Client{
		Transport: tr,
		Timeout:   15 * time.Second,
	}
	return &HTTPSClient{
		Client: client,
	}
}

// PostData 提交post数据
func (c *HTTPSClient) PostData(url string, contentType string, data string) ([]byte, error) {
	resp, err := c.Post(url, contentType, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// HTTPClient http客户端
type HTTPClient struct {
	http.Client
}

// PostData post数据
func (c *HTTPClient) PostData(url, format string, data string) ([]byte, error) {
	resp, err := c.Post(url, format, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
