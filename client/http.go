package client

import (
	"crypto/tls"
	"errors"
	"fmt"
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
	HTTPSC = NewHTTPSClient([]byte{}, []byte{})
}

// HTTPSClient HTTPS客户端结构
type HTTPSClient struct {
	http.Client
}

// GetDefaultClient 返回默认的客户端
func GetDefaultClient() *HTTPSClient {
	return HTTPSC
}

// NewHTTPSClient 获取默认https客户端
func NewHTTPSClient(certPEMBlock, keyPEMBlock []byte) *HTTPSClient {
	config := new(tls.Config)
	if len(certPEMBlock) != 0 && len(keyPEMBlock) != 0 {
		cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
		if err != nil {
			panic("load x509 cert error:" + err.Error())
			return nil
		}
		config = &tls.Config{
			Certificates: []tls.Certificate{
				cert,
			},
		}
	} else {
		config = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	tr := &http.Transport{
		TLSClientConfig: config,
	}
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
		panic(err)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// GetData get数据
func (c *HTTPSClient) GetData(url string) ([]byte, error) {
	resp, err := c.Get(url)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		return []byte{}, errors.New("http.stateCode != 200 : " + fmt.Sprintf("%+v", resp))
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
		panic(err)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
