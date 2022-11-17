package sdk

import (
	"sync"

	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk"
	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/service"
)

var ossPlugin *secretsManagerOssPlugin
var once sync.Once

func New(endpoint, secretName string) (*ProxyOssClient, error) {
	return getOssClient(endpoint, secretName, nil)
}

func NewProxyOssClientWithHandler(endpoint, secretName string, akExpireHandler service.AKExpireHandler) (*ProxyOssClient, error) {
	return getOssClient(endpoint, secretName, akExpireHandler)
}

func getOssClient(endpoint, secretName string, akExpireHandler service.AKExpireHandler) (*ProxyOssClient, error) {
	err := initSecretsManagerPlugin()
	if err != nil {
		return nil, err
	}
	return ossPlugin.getOssClient(endpoint, secretName, akExpireHandler)
}

func closeOssClient(client *ProxyOssClient, secretName string) error {
	err := initSecretsManagerPlugin()
	if err != nil {
		return err
	}
	return ossPlugin.closeOssClient(client, secretName)
}

func Destroy() {
	if ossPlugin != nil {
		ossPlugin.destroy()
	}
}

func initSecretsManagerPlugin() (err error) {
	err = sdk.Init()
	if err != nil {
		return
	}
	once.Do(func() {
		ossPlugin = NewSecretsManagerOssPlugin()
	})
	return
}
