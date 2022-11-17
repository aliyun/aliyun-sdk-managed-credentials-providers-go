package sdk

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"sync"
)

var sdkCorePlugin *secretsManagerSdkCorePlugin
var once sync.Once

func GetClient(client interface{}, regionId string, secretName string) (interface{}, error) {
	return GetClientWithOptions(client, regionId, nil, secretName)
}
func GetClientWithOptions(client interface{}, regionId string, config *sdk.Config, secretName string) (interface{}, error) {
	client, err := sdkCorePlugin.getClient(client, regionId, config, secretName)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func CloseSDKCoreClient(client interface{}, secretName string) error {
	initSecretsManagerPlugin()
	return sdkCorePlugin.closeSDKCoreClient(client, secretName)
}

func Destroy() {
	if sdkCorePlugin != nil {
		sdkCorePlugin.destroy()
	}
}

func initSecretsManagerPlugin() {
	once.Do(func() {
		sdkCorePlugin = newSecretsManagerSdkCorePlugin()
	})
	return
}
